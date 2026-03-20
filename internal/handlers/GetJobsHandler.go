package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/azuar4e/api-gateway-tfg/internal/initializers"
	"github.com/azuar4e/api-gateway-tfg/internal/models"
	"github.com/gin-gonic/gin"
)

type DynamoInterfaceQuery interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

var DynamoClientQuery DynamoInterfaceQuery = initializers.DY

func GetJobsHandler(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid := user.(models.User).ID

	out, err := DynamoClientQuery.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("TABLE_NAME")),
		KeyConditionExpression: aws.String("PK = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(uid), 10)},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list jobs"})
		return
	}

	jobs := make([]models.Job, 0, len(out.Items))
	for _, m := range out.Items {
		var item models.JobDynamoItem
		if err := attributevalue.UnmarshalMap(m, &item); err != nil {
			continue
		}
		jobs = append(jobs, item.ToJob())
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}
