package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/azuar4e/microservices-tfg/internal/initializers"
	"github.com/azuar4e/microservices-tfg/internal/models"
	"github.com/gin-gonic/gin"
)

func GetJobByIdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	u := user.(models.User)
	userID := u.ID
	jobID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	out, err := initializers.DY.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("PriceAlerts"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(userID), 10)},
			"SK": &types.AttributeValueMemberN{Value: strconv.FormatInt(jobID, 10)},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get job"})
		return
	}
	if out.Item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	var item models.JobDynamoItem
	if err := attributevalue.UnmarshalMap(out.Item, &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read job"})
		return
	}

	c.JSON(http.StatusOK, item.ToJob())
}
