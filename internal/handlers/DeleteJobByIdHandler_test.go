package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/azuar4e/api-gateway-tfg/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type mockDynamoDel struct{}

func (m *mockDynamoDel) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, nil
}

func TestDeleteJobById(t *testing.T) {
	DynamoClientDel = &mockDynamoDel{}
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req, _ := http.NewRequest("DELETE", "/jobs/1", nil)
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctx.Set("user", models.User{
		Model:    gorm.Model{ID: 1},
		Email:    "tfg@tfg.com",
		Password: "1234",
	})

	DeleteJobHandler(ctx)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusNoContent)
	}
}
