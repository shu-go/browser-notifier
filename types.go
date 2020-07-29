package notifier

import (
	"time"
)

type Notification struct {
	App       string    `json:"app"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Timestamp string    `json:"timestamp"`
	RawTS     time.Time `json:"-"`
}
