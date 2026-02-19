package initializers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DY *dynamodb.Client

func ConnectToDynamo() {
	// INFO: ver historial de commits para config localstack
	ctx := context.Background()
	dynamoConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("failed to load AWS config: " + err.Error())
	}

	DY = dynamodb.NewFromConfig(dynamoConfig)
}
