package awsutils

import (
	"context"
	"fmt"
	"sapopinguino-translate/internal/config"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

var (
	APIGatewayClient *apigatewaymanagementapi.Client
)

func ConfigAWS() error {
	_, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		return fmt.Errorf("Error while loading the AWS config: %s", err)
	}

	return nil
}

func ConfigAWSGateway(c *config.Config) error {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		return fmt.Errorf("Error while loading the AWS config: %s", err)
	}

	APIGatewayClient = apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = c.WebsocketEndpoint()
	})

	return nil
}
