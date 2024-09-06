package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iqbalrosiadi/wati_hackation_dikes/repo"
	model "github.com/iqbalrosiadi/wati_hackation_dikes/repo"
	gorse "github.com/zhenghaoz/gorse/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var RANDOM_LABELS = []string{
	"Nebula",
	"Stardust",
	"Stellar",
	"Galactic",
	"Interstellar",
	"Celestial",
	"Cosmic",
	"Quantum",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Canvas",
	"Collage",
	"Nights",
	"Ink",
	"Leap",
	"Serenade",
	"Kaleidoscope",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
	"Symphony",
	"Spectacle",
	"Graffiti",
	"Ink",
	"Canvas",
	"Collage",
	"Nights",
}

func main() {
	uri := "mongodb://localhost:27017/dikes_hackathon"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	dbName := "dikes_hackathon"
	broadcastCollName := "Broadcast"
	templateCollName := "Template"
	// compiledContactProfileCollName := "CompiledContactProfile"
	broadcastColl := client.Database(dbName).Collection(broadcastCollName)
	templateColl := client.Database(dbName).Collection(templateCollName)
	// compiledContactProfileColl := client.Database(dbName).Collection(compiledContactProfileCollName)

	broadcastRepo := repo.NewBroadcastRepo(broadcastColl)
	templateRepo := repo.NewTemplateRepo(templateColl)
	// compiledContactProfileRepo := repo.NewCompiledContactProfileRepo(compiledContactProfileColl)

	r := gin.Default()

	// Handle Create message template
	r.POST("/api/v1/message-template", func(c *gin.Context) {
		var messageTemplate model.Template
		if err := c.ShouldBindJSON(&messageTemplate); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		rs, err := templateRepo.Create(context.Background(), messageTemplate)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		id, ok := rs.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(500, gin.H{"error": "Failed to get inserted ID"})
			return
		}
		// Call to AI model to get template's labels
		// Create user
		client := gorse.NewGorseClient("http://127.0.0.1:8087", "")
		if _, err := client.InsertUser(c.Request.Context(), gorse.User{
			UserId: id.String(),
		}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		// Create item with labels
		labels := mockLabels()
		itemId := fmt.Sprintf("template:%s", id.String())
		if _, err := client.InsertItem(c.Request.Context(), gorse.Item{
			ItemId:    itemId,
			Labels:    labels,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}

		c.JSON(200, gin.H{
			"created": true,
		})
	})

	// Handle Get message templates
	r.GET("/api/v1/message-templates", func(c *gin.Context) {
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
	})

	// Handle Get message template by id

	// Handle Get broadcasts
	r.GET("/api/v1/broadcasts", func(c *gin.Context) {
		cursor, err := broadcastRepo.Find(context.Background(), bson.D{})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		// Unpacks the cursor into a slice
		var results []repo.Broadcast
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		c.JSON(200, results)
	})

	// Handle Get broadcasts

	// Handle Create broadcast
	r.POST("/api/v1/broadcast", func(c *gin.Context) {
		var broadcast model.Broadcast
		if err := c.ShouldBindJSON(&broadcast); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := broadcastRepo.Create(context.Background(), broadcast); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"created": true,
		})
	})

	// Handle Get broadcast recommendation contacts
	r.GET("/api/v1/broadcast/recommend-contacts", func(c *gin.Context) {
		// Get Message Template from request
		templateId := c.Query("templateId")

		// Call to AI service // Create a client
		client := gorse.NewGorseClient("http://127.0.0.1:8087", "")

		// // Call to AI model to get template's labels
		// // Create user
		// if _, err := client.InsertUser(c.Request.Context(), gorse.User{
		// 	UserId: templateId,
		// }); err != nil {
		// 	c.JSON(500, gin.H{"error": err.Error()})
		// }
		// // Create item with labels
		// labels := mockLabels()
		// itemId := fmt.Sprintf("template:%s:%d", templateId, time.Now().UnixMilli())
		// if _, err := client.InsertItem(c.Request.Context(), gorse.Item{
		// 	ItemId:    itemId,
		// 	Labels:    labels,
		// 	Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		// }); err != nil {
		// 	c.JSON(500, gin.H{"error": err.Error()})
		// }

		// Insert feedback
		// if _, err := client.InsertFeedback(c.Request.Context(), []gorse.Feedback{
		// 	{FeedbackType: "star", UserId: templateId, ItemId: itemId, Timestamp: "2022-02-24"},
		// }); err != nil {
		// 	c.JSON(500, gin.H{"error": err.Error()})
		// }

		// Get recommendation.
		rs, err := client.GetItemNeighbors(c.Request.Context(), "template:mock_template_999:1725604023024", templateId, 10, 0)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, rs)
	})
	r.Run("0.0.0.0:8099") // listen and serve on 0.0.0.0:8080
}

func mockLabels() []string {
	rand.Shuffle(len(RANDOM_LABELS), func(i, j int) {
		RANDOM_LABELS[i], RANDOM_LABELS[j] = RANDOM_LABELS[j], RANDOM_LABELS[i]
	})
	return RANDOM_LABELS[:5]
}
