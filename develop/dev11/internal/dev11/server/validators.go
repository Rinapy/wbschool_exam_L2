package server

import "time"

func (e *EventRaw) IsValid() bool {
	if e.UserUID == "" {
		return false
	}
	if _, err := time.Parse(time.DateOnly, e.StartDT); err != nil {
		return false
	}
	if _, err := time.Parse(time.DateOnly, e.EndDT); err != nil {
		return false
	}
	return true
}

func (e *EventRawUpdate) IsValid() bool {
	if e.UID == "" {
		return false
	}
	return true
}
