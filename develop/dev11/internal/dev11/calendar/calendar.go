package calendar

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type EventNotFound struct{}

func (err *EventNotFound) Error() string {
	return "event with such UID was not found"
}

type Event struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserUID     string    `json:"user_uid"`
	StartDT     time.Time `json:"start_dt"`
	EndDT       time.Time `json:"end_dt"`
}

func NewEvent() *Event {
	return &Event{}
}

type Calendar struct {
	events map[string]*Event
	mu     *sync.RWMutex
}

func NewCalendar() *Calendar {
	return &Calendar{
		events: make(map[string]*Event),
		mu:     &sync.RWMutex{},
	}
}

func (c *Calendar) UIDGen() (string, error) {
	var key string
	for {
		b := uuid.New()
		key = b.String()
		if _, ok := c.events[key]; !ok {
			break
		}
	}
	return key, nil
}

func (c *Calendar) GetEvent(UID string) (*Event, error) {
	c.mu.RLock()

	if e, ok := c.events[UID]; !ok {
		return nil, &EventNotFound{}
	} else {
		return e, nil
	}
}

func (c *Calendar) CreateEvent(e *Event) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	UID, err := c.UIDGen()
	if err != nil {
		return "", err
	}
	c.events[UID] = e
	return UID, nil

}

func (c *Calendar) UpdateEvent(e *Event, eUID string) (*Event, error) {
	if changEvent, err := c.GetEvent(eUID); err != nil {
		return nil, err
	} else {
		c.mu.Lock() // Может ли тут быть гонка ??
		defer c.mu.Unlock()
		if e.Name != "" {
			changEvent.Name = e.Name
		}
		if e.Description != "" {
			changEvent.Description = e.Description
		}
		if !e.StartDT.IsZero() {
			changEvent.StartDT = e.StartDT
		}
		if !e.EndDT.IsZero() {
			changEvent.EndDT = e.EndDT
		}
		return changEvent, nil
	}
}

func (c *Calendar) DeleteEvent(UID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.events[UID]
	if !found {
		return &EventNotFound{}
	} else {
		delete(c.events, UID)
	}
	return nil
}
