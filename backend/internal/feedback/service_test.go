package feedback

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubRepository struct {
	data *PromptData
}

func (r stubRepository) GetAttemptDataForFeedback(context.Context, string, string) (*PromptData, error) {
	return r.data, nil
}

type stubLLM struct {
	response string
	err      error
}

func (l stubLLM) GenerateJSON(context.Context, string, string) (string, error) {
	return l.response, l.err
}

func TestGenerateSetsAISource(t *testing.T) {
	service := NewService(
		stubRepository{data: testPromptData()},
		stubLLM{response: `{
			"strengths":["Проверил ссылку"],
			"weaknesses":[],
			"risk_profile":{"dominant_risk":"phishing","risk_count":1,"description":"Описание"},
			"recommendations":["Не переходить по ссылкам"],
			"learning_tips":["Проверять адрес сайта"],
			"motivation":"Отлично"
		}`},
		testLogger(),
	)

	result, err := service.Generate(context.Background(), "user", "attempt")

	require.NoError(t, err)
	require.Equal(t, FeedbackSourceAI, result.Source)
}

func TestGenerateSetsFallbackSourceWhenLLMFails(t *testing.T) {
	service := NewService(
		stubRepository{data: testPromptData()},
		stubLLM{err: errors.New("unavailable")},
		testLogger(),
	)

	result, err := service.Generate(context.Background(), "user", "attempt")

	require.NoError(t, err)
	require.Equal(t, FeedbackSourceFallback, result.Source)
}

func testPromptData() *PromptData {
	return &PromptData{
		Role:                "buyer",
		ScenarioTitle:       "Проверка ссылки",
		ScenarioDescription: "Тестовый сценарий",
		TotalScore:          100,
		Answers: []AnswerData{
			{
				StepNumber:     1,
				NodeQuestion:   "Перейти по ссылке?",
				ChoiceText:     "Проверить ссылку",
				ChoiceScore:    100,
				RiskCategories: "Нет",
				Consequence:    "Данные защищены",
				Explanation:    "Проверка адреса снижает риск фишинга",
			},
		},
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
