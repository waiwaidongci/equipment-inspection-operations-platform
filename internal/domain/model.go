package domain

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"displayName"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Device struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Model        string    `json:"model"`
	SerialNumber string    `json:"serialNumber"`
	Manufacturer string    `json:"manufacturer"`
	Location     string    `json:"location"`
	InstallDate  string    `json:"installDate"`
	WarrantyEnd  string    `json:"warrantyEnd"`
	Owner        string    `json:"owner"`
	Status       string    `json:"status"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type InspectionPlan struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	ScopeType string          `json:"scopeType"`
	DeviceID  *int64          `json:"deviceId,omitempty"`
	Category  string          `json:"category"`
	Period    string          `json:"period"`
	StartDate string          `json:"startDate"`
	EndDate   string          `json:"endDate"`
	Owner     string          `json:"owner"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Items     []PlanCheckItem `json:"items,omitempty"`
}

type PlanCheckItem struct {
	ID            int64  `json:"id"`
	PlanID        int64  `json:"planId"`
	Name          string `json:"name"`
	ItemType      string `json:"itemType"`
	StandardRange string `json:"standardRange"`
	Required      bool   `json:"required"`
	Description   string `json:"description"`
	SortOrder     int    `json:"sortOrder"`
}

type InspectionTask struct {
	ID               int64             `json:"id"`
	PlanID           int64             `json:"planId"`
	DeviceID         int64             `json:"deviceId"`
	PlannedAt        time.Time         `json:"plannedAt"`
	Status           string            `json:"status"`
	ActualExecutorID *int64            `json:"actualExecutorId,omitempty"`
	ExecutedAt       *time.Time        `json:"executedAt,omitempty"`
	Results          []TaskCheckResult `json:"results"`
	Remark           string            `json:"remark"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
}

type TaskCheckResult struct {
	ItemID        int64  `json:"itemId"`
	Name          string `json:"name"`
	ItemType      string `json:"itemType"`
	StandardRange string `json:"standardRange"`
	Required      bool   `json:"required"`
	Description   string `json:"description"`
	Value         string `json:"value"`
	Passed        bool   `json:"passed"`
	Abnormal      bool   `json:"abnormal"`
	Remark        string `json:"remark"`
}

type InspectionRecord struct {
	ID            int64     `json:"id"`
	TaskID        int64     `json:"taskId"`
	DeviceID      int64     `json:"deviceId"`
	PlanID        int64     `json:"planId"`
	ExecutorID    int64     `json:"executorId"`
	ExecutedAt    time.Time `json:"executedAt"`
	Status        string    `json:"status"`
	ResultSummary string    `json:"resultSummary"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Anomaly struct {
	ID                 int64      `json:"id"`
	DeviceID           int64      `json:"deviceId"`
	TaskID             *int64     `json:"taskId,omitempty"`
	DiscovererID       int64      `json:"discovererId"`
	DiscoveredAt       time.Time  `json:"discoveredAt"`
	Severity           string     `json:"severity"`
	Description        string     `json:"description"`
	Status             string     `json:"status"`
	AssigneeID         *int64     `json:"assigneeId,omitempty"`
	Progress           string     `json:"progress"`
	CauseAnalysis      string     `json:"causeAnalysis"`
	CloseNote          string     `json:"closeNote"`
	VerificationResult string     `json:"verificationResult"`
	ClosedAt           *time.Time `json:"closedAt,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

type Repair struct {
	ID         int64     `json:"id"`
	DeviceID   int64     `json:"deviceId"`
	AnomalyID  *int64    `json:"anomalyId,omitempty"`
	RepairType string    `json:"repairType"`
	Content    string    `json:"content"`
	Vendor     string    `json:"vendor"`
	StartedAt  string    `json:"startedAt"`
	EndedAt    string    `json:"endedAt"`
	Cost       float64   `json:"cost"`
	Result     string    `json:"result"`
	Remark     string    `json:"remark"`
	CreatedBy  int64     `json:"createdBy"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type DeviceHistory struct {
	Records []InspectionRecord `json:"records"`
	Repairs []Repair           `json:"repairs"`
}
