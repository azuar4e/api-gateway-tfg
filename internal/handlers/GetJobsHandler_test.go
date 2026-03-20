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

type mockDynamoQuery struct{}

func (m *mockDynamoQuery) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return &dynamodb.QueryOutput{}, nil
}

func TestGetJobs(t *testing.T) {
	DynamoClientQuery = &mockDynamoQuery{}
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req, _ := http.NewRequest("Get", "/jobs", nil)
	ctx.Request = req
	ctx.Set("user", models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:    "tfg@tfg.com",
		Password: "1234",
	})
	GetJobsHandler(ctx)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusOK)
	}

}
