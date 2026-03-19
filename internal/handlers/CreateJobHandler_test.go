package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/azuar4e/api-gateway-tfg/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type mockDynamo struct{}

func (m *mockDynamo) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

type mockSQS struct{}

func (m *mockSQS) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return &sqs.SendMessageOutput{}, nil
}

func TestCreateJob(t *testing.T) {
	DynamoClient = &mockDynamo{}
	SQSClient = &mockSQS{}
	gin.SetMode(gin.TestMode)

	body := `{
        "url": "https://www.amazon.es/producto",
        "target_price": 10
    }`

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req, _ := http.NewRequest("POST", "/jobs", strings.NewReader(body))
	ctx.Request = req
	ctx.Set("user", models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:    "tfg@tfg.com",
		Password: "1234",
	})
	CreateJobHandler(ctx)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusAccepted)
	}

}
