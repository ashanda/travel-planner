package models

type TripRequest struct {
	Destination string   `json:"destination" binding:"required"`
	StartDate   string   `json:"start_date"` // YYYY-MM-DD optional
	Days        int      `json:"days" binding:"required,min=1,max=30"`
	Budget      string   `json:"budget"` // low|mid|high
	Interests   []string `json:"interests"`
	Pace        string   `json:"pace"` // chill|balanced|fast
	Notes       string   `json:"notes"`
}

type TripPlan struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	InputHash string      `json:"input_hash"`
	Request   TripRequest `json:"request"`

	Itinerary any `json:"itinerary"`

	// ✅ NEW: for frontend display
	Weather any `json:"weather,omitempty"`
	Places  any `json:"places,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
