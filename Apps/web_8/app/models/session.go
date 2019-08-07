package models

import "time"

// Session is model
type Session struct {
	SessionID  string
	Account    string
	ExpireDate time.Time
}
