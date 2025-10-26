package reset

// Config holds the application configuration
type Config struct {
	Token string
}

// ActiveSubscription represents a filtered active subscription with essential fields
type ActiveSubscription struct {
	ID             int     `json:"id"`
	ResetTimes     int     `json:"resetTimes"`
	CurrentCredits float64 `json:"currentCredits"`
	CreditLimit    float64 `json:"creditLimit"` // Total credit limit
}

// Stats tracks reset operation statistics
type Stats struct {
	Total   int
	Success int
	Failed  int
}

// SubscriptionResponse represents the API response structure
type SubscriptionResponse struct {
	Code     int            `json:"code"`
	Level    any            `json:"level"`
	Msg      string         `json:"msg"`
	OK       bool           `json:"ok"`
	Data     []Subscription `json:"data"`
	DataType int            `json:"dataType"`
}

// Subscription represents a single subscription
type Subscription struct {
	ResetTimes           int              `json:"resetTimes"`
	ID                   int              `json:"id"`
	EmployeeID           int              `json:"employeeId"`
	EmployeeName         any              `json:"employeeName"`
	EmployeeEmail        string           `json:"employeeEmail"`
	CurrentCredits       float64          `json:"currentCredits"`
	LastCreditUpdate     any              `json:"lastCreditUpdate"`
	SubscriptionPlanID   int              `json:"subscriptionPlanId"`
	SubscriptionPlanName string           `json:"subscriptionPlanName"`
	Cost                 float64          `json:"cost"`
	StartDate            string           `json:"startDate"`
	EndDate              string           `json:"endDate"`
	BillingCycle         string           `json:"billingCycle"`
	BillingCycleDesc     string           `json:"billingCycleDesc"`
	RemainingDays        int              `json:"remainingDays"`
	SubscriptionStatus   string           `json:"subscriptionStatus"`
	SubscriptionPlan     SubscriptionPlan `json:"subscriptionPlan"`
	IsActive             bool             `json:"isActive"`
	AutoRenew            bool             `json:"autoRenew"`
	AutoResetWhenZero    bool             `json:"autoResetWhenZero"`
	LastCreditReset      any              `json:"lastCreditReset"`
	CreatedBy            any              `json:"createdBy"`
	CreatedAt            string           `json:"createdAt"`
	UpdatedAt            string           `json:"updatedAt"`
}

// SubscriptionPlan represents the subscription plan details
type SubscriptionPlan struct {
	ID                     int     `json:"id"`
	SubscriptionName       string  `json:"subscriptionName"`
	BillingCycle           string  `json:"billingCycle"`
	Cost                   float64 `json:"cost"`
	OriginalPrice          any     `json:"originalPrice"`
	Features               string  `json:"features"`
	HotTag                 string  `json:"hotTag"`
	ConcurrencyLimit       int     `json:"concurrencyLimit"`
	CreditLimit            float64 `json:"creditLimit"`
	EnableModelRestriction bool    `json:"enableModelRestriction"`
	RestrictedModels       any     `json:"restrictedModels"`
	IsPublic               any     `json:"isPublic"`
	PlanType               string  `json:"planType"`
	SortOrder              any     `json:"sortOrder"`
	CreatedAt              any     `json:"createdAt"`
	UpdatedAt              any     `json:"updatedAt"`
}
