package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/azuar4e/api-gateway-tfg/internal/initializers"
	"github.com/azuar4e/api-gateway-tfg/internal/models"
	"github.com/gin-gonic/gin"
)

type DynamoInterfaceDel interface {
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

var DynamoClientDel DynamoInterfaceDel

func getDynamoDel() DynamoInterfaceDel {
	if DynamoClientDel != nil {
		return DynamoClientDel
	}
	return initializers.DY
}

func DeleteJobHandler(c *gin.Context) {
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

	_, err = getDynamoDel().DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(userID), 10)},
			"SK": &types.AttributeValueMemberN{Value: strconv.FormatInt(jobID, 10)},
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete job"})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
