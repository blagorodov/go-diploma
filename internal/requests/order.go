package requests

import "time"

type Order struct {
	ID         int       `json:"id,omitempty"`
	UserID     int       `json:"user_id,omitempty"`
	Number     string    `json:"number"`
	Accrual    int       `json:"accrual,omitempty"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}
