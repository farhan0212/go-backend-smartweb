package handlers

import (
	"context"
	"log"
	"math"
	"net/http"
	"time"

	"detection-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Prediction struct {
	Class string    `json:"class"`
	Score float64   `json:"score"`
	Bbox  []float64 `json:"bbox"`
}

type DetectionHandler struct {
	Collection *mongo.Collection
}

func (h *DetectionHandler) CreateDetection(c *gin.Context) {
	var request struct {
		Predictions []Prediction `json:"predictions"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	var savedDetections []models.Detection
	ctx := context.Background()

	for i, pred := range request.Predictions {
		log.Printf("🟡 Received prediction [%d]: class=%s score=%f bbox=%v", 
			i, pred.Class, pred.Score, pred.Bbox)

		// Validasi
		if pred.Class == "" || len(pred.Bbox) != 4 {
			log.Printf("⚠️ Skipping invalid prediction at index %d", i)
			continue
		}

		detection := models.Detection{
			Label:      pred.Class,
			Confidence: math.Round(pred.Score * 100), // Convert to percentage
			BoundingBox: models.BoundingBox{
				X:      pred.Bbox[0],
				Y:      pred.Bbox[1],
				Width:  pred.Bbox[2],
				Height: pred.Bbox[3],
			},
			CreatedAt: time.Now(),
		}

		result, err := h.Collection.InsertOne(ctx, detection)
		if err != nil {
			log.Printf("❌ Failed to save detection: %v", err)
			continue
		}

		detection.ID = result.InsertedID.(primitive.ObjectID)
		savedDetections = append(savedDetections, detection)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Detections saved successfully",
		"count":      len(savedDetections),
		"detections": savedDetections,
	})
}