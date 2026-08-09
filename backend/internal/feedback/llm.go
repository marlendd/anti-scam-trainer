package feedback

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type LLMProvider interface {
	GenerateJSON(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

type OpenRouterLLMProvider struct {
	client *openai.Client
	model  string
}

func NewOpenRouterLLMProvider(apiKey, baseURL, model string) *OpenRouterLLMProvider {
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)

	return &OpenRouterLLMProvider{
		client: &client,
		model:  model,
	}
}

func (p *OpenRouterLLMProvider) GenerateJSON(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	response, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: p.model,

		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String(systemPrompt),
					},
				},
			},
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(userPrompt),
					},
				},
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("ошибка при обращении к llm: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("модель вернула пустой ответ")
	}

	return response.Choices[0].Message.Content, nil
}
