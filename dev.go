// dev.go
//go:build !prod && !dev

package main

import (
	"context"
	"log"
	aiutils "sapopinguino-translate/internal/ai"
	"sapopinguino-translate/internal/config"

	"github.com/openai/openai-go/v2"
)

func main() {
	config.LoadDotEnv()

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	aiutils.ConfigOpenAI(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := `{"source_language":"English","target_language":"Spanish","input":"Hi! How are you doing!? I heard that you won the comp. That's amazing!!"}`

	tokenStreamChannel := aiutils.StreamResponse(ctx, openai.ChatModelGPT5Nano, input)

	for res := range tokenStreamChannel {
		if res.Error != nil {
			log.Printf("error: %v", res.Error)
		}

		log.Printf("res: %v", res.Response)
	}
}
