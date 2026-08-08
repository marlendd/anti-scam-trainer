package scenario

type RiskCategory string
type ChoiceScore uint8
type Weight uint8

type ScenarioID string
type LogicalScenarioID string
type NodeID string
type ChoiceID string
type EndingID string
type AuthorID string

const (
	RoleBuyer  Role = "buyer"
	RoleSeller Role = "seller"

	ScoreCritical ChoiceScore = 0
	ScoreRisky    ChoiceScore = 50
	ScoreSafe     ChoiceScore = 100

	WeightLow    Weight = 1
	WeightMedium Weight = 2
	WeightHigh   Weight = 3

	RiskExternalMessenger   RiskCategory = "external_messenger"
	RiskExternalLink        RiskCategory = "external_link"
	RiskSMSCode             RiskCategory = "sms_code"
	RiskOffPlatformPayment  RiskCategory = "off_platform_payment"
	RiskFakePaymentDelivery RiskCategory = "fake_payment_delivery"
	RiskUrgencyPressure     RiskCategory = "urgency_pressure"
	RiskDisableProtection   RiskCategory = "disable_protection"
	RiskUnverifiedEvidence  RiskCategory = "unverified_evidence"
)

type Role string

type Scenario struct {
	ID          ScenarioID        // ID версии сценария
	LogicalID   LogicalScenarioID // ID сценария
	Version     int
	Role        Role
	Title       string
	Description string
	StartNodeID NodeID
	Nodes       []Node
	Endings     []Ending
}

type Node struct {
	ID      NodeID   `json:"id"`
	Author  AuthorID `json:"author"`
	Text    string   `json:"text"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	ID             ChoiceID       `json:"id"`
	Text           string         `json:"text"`
	Consequence    string         `json:"consequence"`
	Explanation    string         `json:"explanation"`
	Weight         Weight         `json:"weight"` // важность выбора
	Score          ChoiceScore    `json:"score"`  // оценка безопасности выбора (больше -> безопаснее)
	RiskCategories []RiskCategory `json:"risk_categories"`
	NextNodeID     NodeID         `json:"next_node_id,omitempty"`
	EndingID       EndingID       `json:"ending_id,omitempty"`
	// у одного выбора должно быть заполнено только одно поле из пары NextNodeID / EndingID.
}

type Ending struct {
	ID     EndingID `json:"id"`
	Header string   `json:"header"`
	Result string   `json:"result"`
}

type Content struct {
	StartNodeID NodeID   `json:"start_node_id"`
	Nodes       []Node   `json:"nodes"`
	Endings     []Ending `json:"endings"`
}
