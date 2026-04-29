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

type DynamoInterfaceDelAll interface {
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

var DynamoClientDelAll DynamoInterfaceDelAll

func getDynamoDelAll() DynamoInterfaceDelAll {
	if DynamoClientDelAll != nil {
		return DynamoClientDelAll
	}
	return initializers.DY
}

func DeleteJobsHandler(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid := user.(models.User).ID

	result, err := getDynamoDelAll().Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String(os.Getenv("TABLE_NAME")),
		FilterExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(uid), 10)},
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan jobs"})
		return
	}

	if len(result.Items) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	var requests []types.WriteRequest
	for _, item := range result.Items {
		requests = append(requests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"PK": item["PK"],
					"SK": item["SK"],
				},
			},
		})
	}

	for i := 0; i < len(requests); i += 25 {
		end := i + 25
		if end > len(requests) {
			end = len(requests)
		}
		_, err = getDynamoDelAll().BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				os.Getenv("TABLE_NAME"): requests[i:end],
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete jobs"})
			return
		}
	}

	c.AbortWithStatus(http.StatusNoContent)
}
