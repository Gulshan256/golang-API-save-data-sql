// models/waitlist.go
package models

type Waitlist struct {
	Date          string `json:"date"`
	Time          string `json:"time"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	FromLocation  string `json:"from"`
	ComponentName string `json:"componentName"`
}

func NewWaitlist(date string, time string, name string, email string, phone string, from string, componentName string) *Waitlist {
	return &Waitlist{date, time, name, email, phone, from, componentName}
}