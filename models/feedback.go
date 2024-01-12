package models


type Feedback struct {
	Date     string `json:"date"`
	Time     string `json:"time"`
	Feedback string `json:"feedback"`
}


func NewFeedback(date string, time string, feedback string) *Feedback {
		return &Feedback{date, time, feedback}
	}