package feedback

// Данные для генерации промпта, взятые из бд

type PromptData struct {
	Role                string
	ScenarioTitle       string
	ScenarioDescription string
	TotalScore          int
	Answers             []AnswerData
}

type AnswerData struct {
	StepNumber     int
	NodeQuestion   string
	ChoiceText     string
	ChoiceScore    int
	RiskCategories string // В виде строки через запятую
	Consequence    string
	Explanation    string
}

// Вспомогательная структура для парсинга JSONB из БД answers.response
type AnswerResponseRaw struct {
	NodeQuestion string `json:"node_question"`
	ChoiceText   string `json:"choice_text"`
}

// Структуры ответа от LLM

type RiskProfile struct {
	DominantRisk string `json:"dominant_risk"`
	RiskCount    int    `json:"risk_count"`
	Description  string `json:"description"`
}

type AttemptFeedback struct {
	Strengths       []string    `json:"strengths"`
	Weaknesses      []string    `json:"weaknesses"`
	RiskProfile     RiskProfile `json:"risk_profile"`
	Recommendations []string    `json:"recommendations"`
	LearningTips    []string    `json:"learning_tips"`
	Motivation      string      `json:"motivation"`
}
