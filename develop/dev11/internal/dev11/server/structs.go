package server

import (
	"dev11/internal/dev11/calendar"
	"fmt"
	"time"
)

type EventRaw struct {
	UserUID     string `json:"user_uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDT     string `json:"start_dt"`
	EndDT       string `json:"end_dt"`
}

func NewEventRaw() *EventRaw {
	return &EventRaw{}
}

func (e *EventRaw) convertToEvent() *calendar.Event {
	eve := calendar.NewEvent()
	eve.UserUID = e.UserUID
	eve.Name = e.Name
	eve.Description = e.Description
	eve.StartDT, _ = time.Parse(time.DateOnly, e.StartDT)
	fmt.Println(time.Parse(time.DateOnly, e.StartDT))
	return eve
}

type EventRawUpdate struct {
	UID         string `json:"uid"`
	UserUID     string `json:"user_uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDT     string `json:"start_dt"`
	EndDT       string `json:"end_dt"`
}

func NewEventUpdateRaw() *EventRawUpdate {
	return &EventRawUpdate{}
}

func (e *EventRawUpdate) convertToEvent() *calendar.Event {
	eve := calendar.NewEvent()
	eve.UserUID = e.UserUID
	eve.Name = e.Name
	eve.Description = e.Description
	eve.StartDT, _ = time.Parse(time.RFC3339, e.StartDT)
	eve.EndDT, _ = time.Parse(time.RFC3339, e.EndDT)
	return eve
}

// Result тип для ответа, имеет поле result
type Result struct {
	Result interface{} `json:"result"`
}

// Error тип для ответа, имеет поле error
type Error struct {
	Error interface{} `json:"error"`
}
