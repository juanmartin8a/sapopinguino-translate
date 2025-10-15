// main.go
//go:build prod || dev

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	aiutils "sapopinguino-translate/internal/ai"
	awsutils "sapopinguino-translate/internal/aws"
	"sapopinguino-translate/internal/config"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/openai/openai-go/v2"
)

func init() {
	err := awsutils.ConfigAWS()
	if err != nil {
		log.Fatal(err)
	}

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = awsutils.ConfigAWSGateway(c)
	if err != nil {
		log.Fatal(err)
	}

	err = aiutils.ConfigOpenAI(c)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := event.RequestContext.ConnectionID

	bodyBytes := []byte(event.Body)

	var bodyS awsutils.Body

	err := json.Unmarshal(bodyBytes, &bodyS)
	if err != nil {
		error := fmt.Errorf("Failed to unmarshal request's body: %v", err)
		log.Println(error)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `"Internal server error :/"`,
		}, error
	}

	tokenStreamChannel := aiutils.StreamResponse(ctx, openai.ChatModelGPT5Nano, bodyS.Message)

	for res := range tokenStreamChannel {
		if res.Error != nil {
			log.Printf("Error while streaming LLM's response: %v", res.Error)
			_, err = awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: &connectionID,
				Data:         []byte("<error:/>"),
			})
			if err != nil {
				log.Printf("Error sending error token to client: %v", err)
				awsutils.HandleDeleteConnection(ctx, &connectionID, "sending \"<error:/>\" in PostConnection")
			}
			break
		}

		var jsonData []byte
		jsonData, err = json.Marshal(res.Response)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			_, err := awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: &connectionID,
				Data:         []byte("<error:/>"),
			})
			if err != nil {
				log.Printf("Error sending error token to client: %v", err)
				awsutils.HandleDeleteConnection(ctx, &connectionID, "sending \"<error:/>\" in PostConnection")
			}
			break
		}

		_, err = awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         jsonData,
		})
		if err != nil {
			if strings.Contains(err.Error(), "410") {
				log.Printf("Client disconnected: %v", err)
				break
			} else {
				log.Printf("Error sending token to client: %v", err)
				awsutils.HandleDeleteConnection(ctx, &connectionID, "sending token in PostConnection")
				break
			}
		}
	}

	if err == nil {
		_, err = awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         []byte("<end:)>"),
		})
		if err != nil {
			log.Printf("Error sending <end:)> thingy to client: %v", err)
			awsutils.HandleDeleteConnection(ctx, &connectionID, "sending \"<end:/>\" in PostConnection")
		}
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `"SIIUUUUU! :D"`,
	}

	return response, nil
}

func main() {
	lambda.Start(handler)
}
