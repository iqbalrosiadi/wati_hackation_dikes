package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iqbalrosiadi/wati_hackation_dikes/repo"
	model "github.com/iqbalrosiadi/wati_hackation_dikes/repo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	gorse "github.com/zhenghaoz/gorse/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	pkgConfig "github.com/ClareAI/wati-go-common/pkg/config"
	pkgMongo "github.com/ClareAI/wati-go-common/pkg/database/mongo"
	pkgLogger "github.com/ClareAI/wati-go-common/pkg/logger"
)

type RecommendContact struct {
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	ItemId string `json:"itemId"`
}

const (
	BROADCAST_COLLECTION_NAME = "Broadcast"
	TEMPLATE_COLLECTION_NAME  = "Template"
	CONTACT_COLLECTION_NAME   = "Contact"
)

var (
	templateRepo *repo.TemplateRepo
	contactRepo  *repo.ContactRepo
	logger       zerolog.Logger
	config       *viper.Viper
	mongoManager *pkgMongo.MongoManager
	db           *mongo.Database
	gorseClient  *gorse.GorseClient
)

func init() {
	config = pkgConfig.GetConfig()

	pkgLogger.InitDefaultLogger()
	logger = log.Logger

	pkgMongo.NewMongoManager()
	mongoManager = pkgMongo.GetMongoManager()
	mongoConn := pkgMongo.Config{
		Protocol: config.GetString("database.mongo.protocol"),
		Host:     config.GetString("database.mongo.host"),
		Port:     config.GetString("database.mongo.port"),
		User:     config.GetString("database.mongo.user"),
		Password: config.GetString("database.mongo.password"),
	}
	if err := mongoManager.NewMongoConnection(mongoConn); err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize database")
	}
	dbConn := mongoManager.GetConnection(config.GetString("database.mongo.host")).GetClient()
	db = dbConn.Database(config.GetString("database.mongo.name"))

	gorseClient = gorse.NewGorseClient(fmt.Sprintf("http://%s", net.JoinHostPort(config.GetString("recommender.host"), config.GetString("recommender.port"))), "")

}

func main() {
	defer func() {
		mongoManager.Shutdown()
	}()

	templateRepo = repo.NewTemplateRepo(db.Collection(TEMPLATE_COLLECTION_NAME))
	contactRepo = repo.NewContactRepo(db.Collection(CONTACT_COLLECTION_NAME))

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "DIKES Backend API Server",
		})
	})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                            // Allows all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // HTTP methods
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Template routes
	r.POST("/api/v1/templates", CreateTemplate)
	r.GET("/api/v1/templates", ListTemplate)
	r.GET("/api/v1/templates/:id", GetTemplateById)

	// Recommend contacts
	r.GET("/api/v1/recommend-contacts", RecommendContacts)

	addr := net.JoinHostPort(config.GetString("server.http.host"), config.GetString("server.http.port"))
	r.Run(addr)
}

func CreateTemplate(c *gin.Context) {
	var messageTemplate model.Template
	if err := c.ShouldBindJSON(&messageTemplate); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Insert the template to the database
	var templateId string
	if rs, err := templateRepo.Create(context.Background(), messageTemplate); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		objID, ok := rs.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(500, gin.H{"error": "Failed to get inserted ID"})
			return
		}
		templateId = objID.Hex()
	}
	messageTemplate.Id = templateId

	err := createTemplateOnRecommender(c, messageTemplate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created":    true,
		"templateId": templateId,
	})
}

func ListTemplate(c *gin.Context) {
	cursor, err := templateRepo.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Unpacks the cursor into a slice
	var results []repo.Template
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	c.JSON(200, results)
}

func GetTemplateById(c *gin.Context) {

}

func RecommendContacts(c *gin.Context) {
	// Get Message Template from request
	templateId := c.Query("templateId")

	// Get template from recommender system
	var defaultContactId string
	template, err := gorseClient.GetUser(c, templateId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if len(template.Subscribe) != 0 {
		defaultContactId = template.Subscribe[0]
	}

	// Get recommended user for this template
	rs, err := gorseClient.GetItemNeighbors(c.Request.Context(), defaultContactId, template.UserId, 10, 0)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	recommendContacts := []RecommendContact{}
	for _, elt := range rs {
		result, err := contactRepo.FindById(c.Request.Context(), elt.Id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		if result != nil {
			contact := model.Contact{}
			if err := result.Decode(&contact); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}

			recommendContacts = append(recommendContacts, RecommendContact{
				Phone:  contact.Phone,
				Name:   contact.Name,
				ItemId: contact.Id,
			})
		}
	}
	c.JSON(200, recommendContacts)
}

func createTemplateOnRecommender(ctx context.Context, messageTemplate model.Template) error {
	var (
		httpClient = http.DefaultClient
	)

	// Generate labels for the template
	labelerAddr := net.JoinHostPort(config.GetString("labeler.host"), config.GetString("labeler.port"))
	jsonByte, err := json.Marshal(messageTemplate)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "http://"+labelerAddr+"/api/v1/labeler/template", bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	labelResponse := struct {
		Labels []string `json:"labels"`
	}{}
	if err := json.Unmarshal(respBodyBytes, &labelResponse); err != nil {
		return err
	}

	// Insert mock contacts to the recommender
	tempUserId := uuid.New().String()
	if _, err := gorseClient.InsertItem(ctx, gorse.Item{
		ItemId:    tempUserId,
		Labels:    labelResponse.Labels,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		return err
	}

	// Insert the template to the recommender
	// Since we recommend contacts based on the template, we need to insert the template to the recommender as User
	if _, err := gorseClient.InsertUser(ctx, gorse.User{
		UserId:    messageTemplate.Id,
		Labels:    labelResponse.Labels,
		Subscribe: []string{tempUserId},
	}); err != nil {
		return err
	}

	// Insert mock feedback (the SDK is broken, so we need to use the HTTP API)
	recommenderUrl := fmt.Sprintf("http://%s", net.JoinHostPort(config.GetString("recommender.host"), config.GetString("recommender.port")))
	feedbacks := []gorse.Feedback{
		{FeedbackType: "star", UserId: messageTemplate.Id, ItemId: tempUserId, Timestamp: time.Now().Format("2006-01-02")},
	}
	jsonByte, err = json.Marshal(feedbacks)
	if err != nil {
		return err
	}
	feedbackRequest, err := http.NewRequest(http.MethodPost, recommenderUrl+"/api/feedback", bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	feedbackRequest.Header.Set("Content-Type", "application/json")
	_, err = httpClient.Do(feedbackRequest)
	if err != nil {
		return err
	}
	return nil
}
