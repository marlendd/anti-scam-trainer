package feedback

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"text/template"
	"time"
)

type Service struct {
	repo Repository
	llm  LLMProvider
	log  *slog.Logger
}

func NewService(repo Repository, llm LLMProvider, log *slog.Logger) *Service {
	return &Service{repo: repo, llm: llm, log: log}
}

const systemPrompt = `Ты — эксперт по кибербезопасности и онлайн-мошенничеству на платформе Авито. Твоя задача — проанализировать действия пользователя в симуляторе общения с мошенниками и предоставить персонализированный, полезный и мотивирующий фидбек.
Сгенерируй ответ СТРОГО в JSON-формате без markdown.`

const userTemplate = `КОНТЕКСТ СЦЕНАРИЯ:
Роль: {{.Role}} ({{if eq .Role "buyer"}}покупатель{{else}}продавец{{end}})
Сценарий: {{.ScenarioTitle}}
Описание: {{.ScenarioDescription}}
Итоговая оценка: {{.TotalScore}} из 100

ДЕЙСТВИЯ ПОЛЬЗОВАТЕЛЯ:
{{range .Answers}}
Шаг {{.StepNumber}}: "{{.NodeQuestion}}"
→ Выбранный ответ: "{{.ChoiceText}}"
→ Оценка выбора: {{.ChoiceScore}}/100 {{if eq .ChoiceScore 100}}✅{{else if eq .ChoiceScore 50}}⚠️{{else}}❌{{end}}
→ Категории риска: {{.RiskCategories}}
→ Последствие: {{.Consequence}}
→ Объяснение: {{.Explanation}}
{{end}}

ДОСТУПНЫЕ КАТЕГОРИИ РИСКОВ:
- phishing, social_engineering, fake_payment, data_leak, fake_delivery, verification_scam, prepayment_scam, account_hijack, counterfeit

ЗАДАЧА: Сгенерируй структурированный фидбек в JSON-формате, следуя этим правилам:
1. strengths (массив строк) — что пользователь сделал правильно.
2. weaknesses (массив строк) — где пользователь ошибся.
3. risk_profile (объект dominant_risk, risk_count, description):
   - dominant_risk: СТРОГО один код из списка "ДОСТУПНЫЕ КАТЕГОРИИ РИСКОВ" выше, без изменений, без перевода, без дополнительных слов.
   - risk_count: число.
   - description: связный текст на РУССКОМ языке, объясняющий суть риска. НЕ используй в description английские термины, коды категорий или обозначения в скобках (например, не пиши "(urgency_pressure)" или "(social_engineering)") — только обычный человекочитаемый русский текст.
4. recommendations (массив строк) — применимые советы по безопасности на Авито, на русском языке, без английских терминов.
5. learning_tips (массив строк) — на какие аспекты обратить внимание в будущем, на русском языке.
6. motivation (строка) — ободряющая фраза на русском языке.

Формат ответа: Только чистый JSON, без markdown.`

// riskLabels переводит внутренние коды категорий риска в человекочитаемые
// русские названия. Коды нужны как стабильный контракт с LLM и для внутренней
// логики (подсчет, выбор описания/рекомендаций), а человекочитаемые названия —
// то, что видит пользователь.
var riskLabels = map[string]string{
	"phishing":           "Фишинг",
	"social_engineering": "Социальная инженерия",
	"fake_payment":       "Фальшивая оплата",
	"data_leak":          "Утечка данных",
	"fake_delivery":      "Фальшивая доставка",
	"verification_scam":  "Мошенничество через 'верификацию'",
	"prepayment_scam":    "Мошенничество с предоплатой",
	"account_hijack":     "Угон аккаунта",
	"counterfeit":        "Подделка товара",
	"urgency_pressure":   "Давление срочностью",
	"general":            "Общие риски",
}

// humanizeRisk возвращает человекочитаемое название категории риска.
// Если код неизвестен, возвращает его как есть (fallback).
func humanizeRisk(code string) string {
	if label, ok := riskLabels[strings.TrimSpace(code)]; ok {
		return label
	}
	return code
}

func (s *Service) Generate(ctx context.Context, userID, attemptID string) (AttemptFeedback, error) {
	data, err := s.repo.GetAttemptDataForFeedback(ctx, userID, attemptID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AttemptFeedback{}, errors.New("not_found")
		}
		return AttemptFeedback{}, err
	}

	tmpl, err := template.New("prompt").Parse(userTemplate)
	if err != nil {
		return AttemptFeedback{}, err
	}

	var promptBuf bytes.Buffer
	if err := tmpl.Execute(&promptBuf, data); err != nil {
		return AttemptFeedback{}, err
	}

	llmCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	jsonStr, err := s.llm.GenerateJSON(llmCtx, systemPrompt, promptBuf.String())
	if err != nil {
		s.log.Warn("LLM failed or timed out, using fallback mechanism",
			"error", err, "attempt_id", attemptID)
		return s.generateFallbackFeedback(data), nil
	}

	jsonStr = cleanLLMJSON(jsonStr)

	var feedback AttemptFeedback
	if err := json.Unmarshal([]byte(jsonStr), &feedback); err != nil {
		s.log.Warn("LLM returned invalid JSON, using fallback mechanism",
			"error", err, "attempt_id", attemptID)
		return s.generateFallbackFeedback(data), nil
	}

	// LLM возвращает код категории риска (например "social_engineering") —
	// переводим его в человекочитаемое русское название перед отдачей клиенту.
	feedback.RiskProfile.DominantRisk = humanizeRisk(feedback.RiskProfile.DominantRisk)
	feedback.Source = FeedbackSourceAI

	return feedback, nil
}

func (s *Service) generateFallbackFeedback(data *PromptData) AttemptFeedback {
	var strengths []string
	var weaknesses []string
	riskCounts := make(map[string]int)

	for _, ans := range data.Answers {
		if ans.ChoiceScore == 100 {
			strengths = append(strengths, ans.Explanation)
		} else {
			weaknesses = append(weaknesses, ans.Consequence)

			if ans.RiskCategories != "Нет" && ans.RiskCategories != "" {
				risks := strings.Split(ans.RiskCategories, ", ")
				for _, r := range risks {
					riskCounts[strings.TrimSpace(r)]++
				}
			}
		}
	}

	if len(strengths) == 0 {
		strengths = append(strengths, "В этот раз выявить мошенника не удалось, но это хороший опыт для обучения.")
	}
	if len(weaknesses) == 0 {
		weaknesses = append(weaknesses, "Вы действовали безупречно и не допустили серьезных ошибок.")
	}

	dominantRisk := "general"
	maxRisk := 0
	for risk, count := range riskCounts {
		if count > maxRisk {
			maxRisk = count
			dominantRisk = risk
		}
	}

	var motivation string
	if data.TotalScore >= 80 {
		motivation = "Отличный результат! Вы отлично распознаете мошенничество, но всегда есть куда расти!"
	} else if data.TotalScore >= 50 {
		motivation = "Хорошая попытка! Вы уже многое знаете о безопасности, осталось закрепить несколько важных правил."
	} else {
		motivation = "Не расстраивайтесь! Это только обучение. Каждая ошибка — это ценный опыт, который поможет вам не попасться мошенникам в реальной жизни."
	}

	return AttemptFeedback{
		Source:     FeedbackSourceFallback,
		Strengths:  strengths,
		Weaknesses: weaknesses,
		RiskProfile: RiskProfile{
			DominantRisk: humanizeRisk(dominantRisk),
			RiskCount:    maxRisk,
			Description:  getRiskDescription(dominantRisk),
		},
		Recommendations: getRecommendations(data.Role, dominantRisk),
		LearningTips:    []string{"Внимательно проверяйте профили", "Не переходите по внешним ссылкам", "Используйте безопасную сделку"},
		Motivation:      motivation,
	}
}

func getRiskDescription(risk string) string {
	descriptions := map[string]string{
		"phishing":           "Фишинг нацелен на кражу ваших данных через поддельные ссылки. Будьте осторожны с внешними сайтами.",
		"social_engineering": "Социальная инженерия использует психологическое давление. Не принимайте поспешных решений.",
		"fake_payment":       "Фальшивые чеки и скриншоты оплат — частая уловка. Проверяйте баланс только в приложении банка.",
		"data_leak":          "Передача личных данных или кодов из SMS ведет к потере контроля над аккаунтом или деньгами.",
		"fake_delivery":      "Ложная доставка пытается увести вас с платформы. Пользуйтесь только официальной доставкой Авито.",
		"verification_scam":  "Мошенники маскируются под службу поддержки. Авито никогда не просит данные карты в чате.",
		"prepayment_scam":    "Требование предоплаты без гарантий — высокий риск потерять деньги.",
		"account_hijack":     "Попытка угона аккаунта грозит тем, что от вашего имени начнут обманывать других людей.",
		"counterfeit":        "Риск купить подделку под видом оригинала. Всегда тщательно проверяйте товар при получении.",
		"general":            "Вы столкнулись с общими рисками при сделке. Соблюдайте базовые правила безопасности.",
	}

	if desc, ok := descriptions[risk]; ok {
		return desc
	}
	return descriptions["general"]
}

func getRecommendations(role, risk string) []string {
	recs := []string{
		"Ведите переписку только во встроенном чате Авито.",
		"Не сообщайте коды из SMS, пароли и CVC-коды никому.",
	}

	if role == "buyer" {
		recs = append(recs, "Используйте Авито.Доставку — это защитит ваши деньги до проверки товара.")
		recs = append(recs, "Не переводите предоплату продавцу на карту до получения товара.")
	} else {
		recs = append(recs, "Отправляйте товар только после того, как увидите статус 'Покупатель оплатил товар' в чате Авито.")
		recs = append(recs, "Не переходите по ссылкам 'для получения средств от покупателя' — деньги зачисляются автоматически.")
	}

	return recs
}

func cleanLLMJSON(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}
	if before, ok := strings.CutSuffix(s, "```"); ok {
		s = before
	}
	return strings.TrimSpace(s)
}
