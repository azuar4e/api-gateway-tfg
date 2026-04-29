package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/azuar4e/api-gateway-tfg/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type mockDynamoDelAll struct{}

func (m *mockDynamoDelAll) BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
	return &dynamodb.BatchWriteItemOutput{}, nil
}

func (m *mockDynamoDelAll) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{
		Items: []map[string]types.AttributeValue{
			{
				"PK": &types.AttributeValueMemberN{Value: "1"},
				"SK": &types.AttributeValueMemberN{Value: "100"},
			},
			{
				"PK": &types.AttributeValueMemberN{Value: "1"},
				"SK": &types.AttributeValueMemberN{Value: "200"},
			},
		},
	}, nil
}

func TestDeleteJobs(t *testing.T) {
	DynamoClientDelAll = &mockDynamoDelAll{}
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req, _ := http.NewRequest("DELETE", "/jobs", nil)
	ctx.Request = req

	ctx.Set("user", models.User{
		Model:    gorm.Model{ID: 1},
		Email:    "tfg@tfg.com",
		Password: "1234",
	})

	DeleteJobsHandler(ctx)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusNoContent)
	}
}
