package main

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/iqbalrosiadi/wati_hackation_dikes/ai/labeler"
	"github.com/spf13/viper"

	pkgConfig "github.com/ClareAI/wati-go-common/pkg/config"
)

var (
	config *viper.Viper
)

func init() {
	config = pkgConfig.GetConfig()
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "DIKES AI Model Server",
		})
	})

	// Template routes
	r.POST("/api/v1/labeler/template", GenerateTemplateLabels)
	r.GET("/api/v1/labeler/contact", GenerateTemplateLabels)

	addr := net.JoinHostPort(config.GetString("server.http.host"), config.GetString("server.http.port"))
	r.Run(addr)
}

func GenerateTemplateLabels(c *gin.Context) {
	type Template struct {
		Id      string `json:"id"`
		Content string `json:"content"`
	}

	var payload Template
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	labelerHdl := labeler.NewTemplateLabeler()
	labels, err := labelerHdl.CreateLabelForTemplate(c, payload.Content)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"id":     payload.Id,
		"labels": labels,
	})
}

func GenerateContactLabels(c *gin.Context) {
	labelerHdl := labeler.NewTemplateLabeler()
	labels, err := labelerHdl.CreateLabelForContact(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(labels)

	c.JSON(200, gin.H{
		"labels": labels,
	})
}
