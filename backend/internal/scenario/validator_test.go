package scenario

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func validScenario() Scenario {
	endings := []Ending{
		{
			ID:     "ending-safe",
			Header: "Сделка остановлена вовремя",
			Result: "Пользователь не передал секретный код и сохранил защиту платформы.",
		},
		{
			ID:     "ending-scam",
			Header: "Данные переданы мошеннику",
			Result: "Пользователь сообщил секретный код вне защищённого канала.",
		},
	}

	choices := []Choice{
		{
			ID:          "choice-stay-on-platform",
			Text:        "Остаться в чате площадки",
			Consequence: "Собеседник продолжает настаивать на оформлении доставки.",
			Explanation: "Общение внутри платформы помогает сохранить её защитные механизмы.",
			Weight:      WeightMedium,
			Score:       ScoreSafe,
			NextNodeID:  "node-platform",
		},
		{
			ID:             "choice-go-to-telegram",
			Text:           "Перейти в Telegram",
			Consequence:    "Собеседник присылает ссылку для подтверждения оплаты.",
			Explanation:    "В стороннем мессенджере защита платформы не действует.",
			Weight:         WeightMedium,
			Score:          ScoreRisky,
			RiskCategories: []RiskCategory{RiskExternalMessenger},
			NextNodeID:     "node-telegram",
		},
		{
			ID:          "choice-platform-refuse-link",
			Text:        "Не открывать ссылку",
			Consequence: "Собеседник меняет тактику и просит код подтверждения.",
			Explanation: "Подозрительные ссылки могут вести на поддельные страницы оплаты.",
			Weight:      WeightHigh,
			Score:       ScoreSafe,
			NextNodeID:  "node-sms",
		},
		{
			ID:             "choice-platform-open-link",
			Text:           "Открыть ссылку",
			Consequence:    "На странице появляется форма, похожая на оформление доставки.",
			Explanation:    "Правдоподобное оформление не делает ссылку безопасной.",
			Weight:         WeightMedium,
			Score:          ScoreRisky,
			RiskCategories: []RiskCategory{RiskExternalLink},
			NextNodeID:     "node-sms",
		},
		{
			ID:          "choice-telegram-return-platform",
			Text:        "Вернуться в чат площадки",
			Consequence: "Собеседник продолжает давить и просит подтвердить действие кодом.",
			Explanation: "Возврат на платформу снижает риск, но не делает запрос кода безопасным.",
			Weight:      WeightMedium,
			Score:       ScoreSafe,
			NextNodeID:  "node-sms",
		},
		{
			ID:             "choice-telegram-open-link",
			Text:           "Открыть ссылку из Telegram",
			Consequence:    "Страница запрашивает код из SMS для подтверждения действия.",
			Explanation:    "Ссылка из стороннего мессенджера может вести на фишинговую страницу.",
			Weight:         WeightHigh,
			Score:          ScoreRisky,
			RiskCategories: []RiskCategory{RiskExternalMessenger, RiskExternalLink},
			NextNodeID:     "node-sms",
		},
		{
			ID:          "choice-refuse-sms",
			Text:        "Не сообщать код из SMS",
			Consequence: "Собеседник прекращает общение.",
			Explanation: "Код из SMS нельзя передавать третьим лицам.",
			Weight:      WeightHigh,
			Score:       ScoreSafe,
			EndingID:    "ending-safe",
		},
		{
			ID:             "choice-share-sms",
			Text:           "Сообщить код из SMS",
			Consequence:    "Мошенник получает возможность подтвердить действие от имени пользователя.",
			Explanation:    "Передача кода из SMS создаёт критический риск потери денег или доступа.",
			Weight:         WeightHigh,
			Score:          ScoreCritical,
			RiskCategories: []RiskCategory{RiskSMSCode, RiskUrgencyPressure},
			EndingID:       "ending-scam",
		},
	}
	nodes := []Node{
		{
			ID:     "node-start",
			Author: "seller",
			Text:   "Давайте перейдём в Telegram, там удобнее обсудить сделку.",
			Choices: []Choice{
				choices[0], choices[1],
			},
		},
		{
			ID:     "node-platform",
			Author: "seller",
			Text:   "Тогда откройте ссылку для оформления доставки.",
			Choices: []Choice{
				choices[2], choices[3],
			},
		},
		{
			ID:     "node-telegram",
			Author: "seller",
			Text:   "Вот ссылка для подтверждения оплаты.",
			Choices: []Choice{
				choices[4], choices[5],
			},
		},
		{
			ID:     "node-sms",
			Author: "system",
			Text:   "Собеседник просит назвать код из SMS.",
			Choices: []Choice{
				choices[6], choices[7],
			},
		},
	}

	s := Scenario{
		ID:          "1",
		LogicalID:   "valid-scenario",
		Version:     1,
		Role:        RoleBuyer,
		Title:       "Безопасная доставка",
		Description: "Тренировка распознавания внешних ссылок и запросов SMS-кода.",
		StartNodeID: nodes[0].ID,
		Nodes:       nodes,
		Endings:     endings,
	}

	return s
}

func TestValidate_ValidScenario(t *testing.T) {
	s := validScenario()

	err := Validate(s)

	if err != nil {
		t.Fatalf("expected valid scenario, got error: %v", err)
	}
}

func TestValidate_InvalidScenarios(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Scenario)
		wantErr error
	}{
		{
			name: "empty start node ID",
			mutate: func(s *Scenario) {
				s.StartNodeID = ""
			},
			wantErr: ErrEmptyStartNodeID,
		},
		{
			name: "unknown start node",
			mutate: func(s *Scenario) {
				s.StartNodeID = NodeID("unknown")
			},
			wantErr: ErrUnknownStartNode,
		},
		{
			name: "empty node ID",
			mutate: func(s *Scenario) {
				s.Nodes[0].ID = ""
			},
			wantErr: ErrEmptyNodeID,
		},
		{
			name: "duplicated node ID",
			mutate: func(s *Scenario) {
				s.Nodes = append(s.Nodes, s.Nodes[0])
			},
			wantErr: ErrDuplicateNodeID,
		},
		{
			name: "empty ending ID",
			mutate: func(s *Scenario) {
				s.Endings[0].ID = ""
			},
			wantErr: ErrEmptyEndingID,
		},
		{
			name: "duplicated ending ID",
			mutate: func(s *Scenario) {
				s.Endings = append(s.Endings, s.Endings[0])
			},
			wantErr: ErrDuplicateEndingID,
		},
		{
			name: "empty choice ID",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].ID = ""
			},
			wantErr: ErrEmptyChoiceID,
		},
		{
			name: "duplicated choice ID",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[1].ID = s.Nodes[0].Choices[0].ID
			},
			wantErr: ErrDuplicateChoiceID,
		},
		{
			name: "choice has no target",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].NextNodeID = ""
			},
			wantErr: ErrInvalidChoiceTarget,
		},
		{
			name: "choice has two targets",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].EndingID = s.Endings[0].ID
			},
			wantErr: ErrInvalidChoiceTarget,
		},
		{
			name: "choice refers to unknown node",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].NextNodeID = "unknown-node"
			},
			wantErr: ErrUnknownNode,
		},
		{
			name: "choice refers to unknown ending",
			mutate: func(s *Scenario) {
				s.Nodes[3].Choices[0].EndingID = "unknown-ending"
			},
			wantErr: ErrUnknownEnding,
		},
		{
			name: "weight is too low",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].Weight = Weight(0)
			},
			wantErr: ErrInvalidWeight,
		},
		{
			name: "weight is too high",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].Weight = Weight(4)
			},
			wantErr: ErrInvalidWeight,
		},
		{
			name: "invalid score",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].Score = ChoiceScore(25)
			},
			wantErr: ErrInvalidScore,
		},
		{
			name: "unsafe choice has no risk categories",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[1].RiskCategories = nil
			},
			wantErr: ErrMissingRiskCategory,
		},
		{
			name: "unreachable node",
			mutate: func(s *Scenario) {
				s.Nodes = append(s.Nodes, Node{
					ID: "node-disconnected",
					Choices: []Choice{
						{
							ID:         "choice-disconnected",
							Weight:     WeightLow,
							Score:      ScoreSafe,
							NextNodeID: "node-sms",
						},
					},
				})
			},
			wantErr: ErrUnreachableNode,
		},
		{
			name: "unreachable ending",
			mutate: func(s *Scenario) {
				s.Endings = append(s.Endings, Ending{
					ID: "unreachable ending",
				})
			},
			wantErr: ErrUnreachableEnding,
		},
		{
			name: "too few reachable endings",
			mutate: func(s *Scenario) {
				s.Endings = s.Endings[:1]
				s.Nodes[3].Choices[1].EndingID = s.Endings[0].ID
			},
			wantErr: ErrTooFewReachableEndings,
		},
		{
			name: "cycle in graph",
			mutate: func(s *Scenario) {
				choice := &s.Nodes[3].Choices[0]

				choice.EndingID = ""
				choice.NextNodeID = s.Nodes[0].ID
			},
			wantErr: ErrCycleDetected,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := validScenario()
			tc.mutate(&s)

			err := Validate(s)

			require.Error(t, err)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
