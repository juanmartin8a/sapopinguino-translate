package aiutils

import (
	"context"
	"fmt"

	"sapopinguino-translate/internal/config"

	openai "github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/responses"
)

var OpenAIClient *openai.Client

type StreamRes struct {
	Response string
	Error    error
}

func ConfigOpenAI(c *config.Config) error {

	openaiClient := openai.NewClient(
		option.WithAPIKey(c.OpenAIKey()),
	)
	OpenAIClient = &openaiClient

	return nil
}

func StreamResponse(context context.Context, model string, input string) <-chan StreamRes {

	tokenStreamChannel := make(chan StreamRes)

	go func() {
		stream := OpenAIClient.Responses.NewStreaming(
			context,
			responses.ResponseNewParams{
				Model: model,
				Input: responses.ResponseNewParamsInputUnion{
					OfString: openai.String(input),
				},
				Prompt: responses.ResponsePromptParam{
					ID:      "pmpt_68d6dd3df0cc8195b092c08b02bfe24e05d616bbdf9c857c",
					Version: openai.String("5"),
				},
				Reasoning: openai.ReasoningParam{
					Effort: openai.ReasoningEffortMinimal,
				},
			},
		)

		for stream.Next() {
			data := stream.Current()

			token := data.Delta

			if len(data.Delta) > 0 {
				tokenStreamChannel <- StreamRes{
					Response: token,
					Error:    nil,
				}
			}
		}

		if err := stream.Err(); err != nil {
			tokenStreamChannel <- StreamRes{
				Response: "",
				Error:    fmt.Errorf("Error while or during LLM's response stream : %v", err),
			}
		}

		close(tokenStreamChannel)
	}()

	return tokenStreamChannel
}
