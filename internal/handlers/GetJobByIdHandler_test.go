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

type mockDynamoGet struct{}

func (m *mockDynamoGet) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"PK":           &types.AttributeValueMemberN{Value: "1"},
			"SK":           &types.AttributeValueMemberN{Value: "1"},
			"url":          &types.AttributeValueMemberS{Value: "https://www.amazon.es/producto"},
			"target_price": &types.AttributeValueMemberN{Value: "500"},
		},
	}, nil
}

func TestGetJobById(t *testing.T) {
	DynamoClientGet = &mockDynamoGet{}
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req, _ := http.NewRequest("Get", "/jobs/1", nil)
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctx.Set("user", models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:    "tfg@tfg.com",
		Password: "1234",
	})
	GetJobByIdHandler(ctx)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusOK)
	}

}
