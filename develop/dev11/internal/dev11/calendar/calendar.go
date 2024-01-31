package calendar

import (
	"crypto/rand"
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
	UserUID     string    `json:"userID"`
	StartDT     time.Time `json:"startDT"`
	EndDT       time.Time `json:"endDT"`
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
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		key = string(b)
		if _, ok := c.events[key]; !ok {
			break
		}
	}
	return key, nil
}

func (c *Calendar) GetEvent(eUID string) (*Event, error) {
	c.mu.RLock()

	if e, ok := c.events[eUID]; !ok {
		return nil, &EventNotFound{}
	} else {
		return e, nil
	}
}

func (c *Calendar) CreateEvent(e *Event) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	eUID, err := c.UIDGen()
	if err != nil {
		return "", err
	}
	c.events[eUID] = e
	return eUID, nil

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

func (c *Calendar) DeleteEvent(eUID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.events[eUID]
	if !found {
		return &EventNotFound{}
	} else {
		delete(c.events, eUID)
	}
	return nil
}
