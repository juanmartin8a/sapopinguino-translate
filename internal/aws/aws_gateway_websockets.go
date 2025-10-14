package awsutils

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

type Body struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func HandleDeleteConnection(ctx context.Context, connectionID *string, afterErrorMessage string) {
	_, err := APIGatewayClient.DeleteConnection(ctx, &apigatewaymanagementapi.DeleteConnectionInput{
		ConnectionId: connectionID,
	})
	if err != nil {
		log.Printf("Error while deleting connection after %s failed: %s", afterErrorMessage, err)
	}
}
