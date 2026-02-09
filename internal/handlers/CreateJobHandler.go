package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/azuar4e/api-gateway-tfg/internal/initializers"
	"github.com/azuar4e/api-gateway-tfg/internal/models"
	"github.com/gin-gonic/gin"
)

const tableName = "PriceAlerts"

func CreateJobHandler(c *gin.Context) {
	var req struct {
		URL         string  `json:"url" binding:"required,url"`
		TargetPrice float64 `json:"target_price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jobID := time.Now().UnixMicro()
	item := models.JobDynamoItem{
		PK:          int64(user.(models.User).ID),
		SK:          jobID,
		URL:         req.URL,
		TargetPrice: req.TargetPrice,
		Status:      "pending",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(item)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build item"})
		return
	}

	_, err = initializers.DY.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save job"})
		return
	}

	//logica para encolar en sqs
	_, err = initializers.SQS.SendMessage(
		context.TODO(),
		&sqs.SendMessageInput{
			QueueUrl:    aws.String(os.Getenv("SQS_QUEUE_URL")),
			MessageBody: aws.String(string(body)),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue job", "log": err})
		return
	}

	//logica para suscribir al usuario al topic de sns
	output, err := initializers.SNS.Subscribe(context.TODO(), &sns.SubscribeInput{
		Protocol: aws.String("email"),
		TopicArn: aws.String(os.Getenv("SNS_TOPIC_ARN")),
		Endpoint: aws.String(user.(models.User).Email),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf(
				"failed to subscribe email %v to topic %v",
				user.(models.User).Email,
				os.Getenv("SNS_TOPIC_ARN"),
			),
			"log": err,
		})
		return
	}

	filterPolicy := map[string][]uint{
		"user_id": {user.(models.User).ID},
	}
	filterPolicyJSON, _ := json.Marshal(filterPolicy)

	_, err = initializers.SNS.SetSubscriptionAttributes(context.TODO(), &sns.SetSubscriptionAttributesInput{
		SubscriptionArn: output.SubscriptionArn,
		AttributeName:   aws.String("FilterPolicy"),
		AttributeValue:  aws.String(string(filterPolicyJSON)),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf(
				"failed to set subscription policy of email %v",
				user.(models.User).Email,
			),
			"log": err,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"id":     jobID,
		"status": "Job queued successfully",
		"email":  "Email subscribed successfully",
	})
}
