package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/azuar4e/api-gateway-tfg/internal/initializers"
	"github.com/azuar4e/api-gateway-tfg/internal/models"
	"github.com/gin-gonic/gin"
)

type DynamoInterface interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type SQSInterface interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

var DynamoClient DynamoInterface

func getDynamo() DynamoInterface {
	if DynamoClient != nil {
		return DynamoClient
	}
	return initializers.DY
}

var SQSClient SQSInterface

func getSQS() SQSInterface {
	if SQSClient != nil {
		return SQSClient
	}
	return initializers.SQS
}

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
	userObj, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
		return
	}

	item := models.JobDynamoItem{
		PK:          int64(userObj.ID),
		SK:          jobID,
		URL:         req.URL,
		TargetPrice: req.TargetPrice,
		Status:      "pending",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	job := item.ToJob()
	body, _ := json.Marshal(job)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build item"})
		return
	}

	_, err = getDynamo().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      av,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save job"})
		return
	}

	//logica para encolar en sqs
	_, err = getSQS().SendMessage(
		context.TODO(),
		&sqs.SendMessageInput{
			QueueUrl:    aws.String(os.Getenv("SQS_QUEUE_URL")),
			MessageBody: aws.String(string(body)),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue job", "log": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"id":     jobID,
		"status": "Job queued successfully",
	})
}
