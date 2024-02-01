package calendar

import (
	"fmt"
	"github.com/google/uuid"
	"os"
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
	events      map[string]*Event
	dayEvents   map[string]*Event
	weekEvents  map[string]*Event
	monthEvents map[string]*Event
	mu          *sync.RWMutex
}

func NewCalendar() *Calendar {
	return &Calendar{
		events:      make(map[string]*Event),
		dayEvents:   make(map[string]*Event),
		weekEvents:  make(map[string]*Event),
		monthEvents: make(map[string]*Event),
		mu:          &sync.RWMutex{},
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

func (c *Calendar) getDayEvents(userUID string) []*Event {
	res := make([]*Event, 0)
	c.mu.RLock()
	for _, event := range c.dayEvents {
		if event.UserUID == userUID {
			res = append(res, event)
		}
	}
	return res
}

func (c *Calendar) getWeekEvents(userUID string) []*Event {
	res := make([]*Event, 0)
	c.mu.RLock()
	for _, event := range c.weekEvents {
		if event.UserUID == userUID {
			res = append(res, event)
		}
	}
	return res
}

func (c *Calendar) getMonthEvents(userUID string) []*Event {
	res := make([]*Event, 0)
	c.mu.RLock()
	for _, event := range c.monthEvents {
		if event.UserUID == userUID {
			res = append(res, event)
		}
	}
	return res
}

func (c *Calendar) EventBalanced(target string, sigint <-chan os.Signal) {
	//targetMap := map[string]map[string]*Event{
	//	"day":   c.dayEvents,
	//	"week":  c.weekEvents,
	//	"month": c.monthEvents,
	//} прикол конечно, стоил мне 3 часов жизни и 10 минут мидла
	targetMap := make(map[string]map[string]*Event)
	targetMap["day"] = c.dayEvents
	targetMap["week"] = c.weekEvents
	targetMap["month"] = c.monthEvents
	for {
		select {
		case <-sigint:
			return
		default:
			today := time.Now().Truncate(24 * time.Hour)
			thisWeek := today.AddDate(0, 0, -int(today.Weekday()))
			thisMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
			c.mu.Lock()
			fmt.Println(len(c.events), len(c.weekEvents))
			fmt.Println(c.weekEvents)
			for uid, event := range c.events {
				systemMap := map[string]bool{
					"day":   event.StartDT.Before(today.Add(24*time.Hour)) && event.EndDT.After(today),
					"week":  event.StartDT.After(thisWeek) && event.EndDT.Before(thisWeek.Add(7*24*time.Hour)),
					"month": event.StartDT.After(thisMonth) && event.EndDT.Before(thisMonth.AddDate(0, 1, 0)),
				}
				if systemMap[target] {
					fmt.Println(systemMap[target], target)
					targetMap[target][uid] = event

					delete(c.events, uid)
				}
			}
			for uid, event := range targetMap[target] {
				systemMap := map[string]bool{
					"day":   event.StartDT.Before(today.Add(24*time.Hour)) && event.EndDT.After(today),
					"week":  event.StartDT.After(thisWeek) && event.EndDT.Before(thisWeek.Add(7*24*time.Hour)),
					"month": event.StartDT.After(thisMonth) && event.EndDT.Before(thisMonth.AddDate(0, 1, 0)),
				}
				if !systemMap[target] {
					fmt.Println(systemMap[target], target)
					c.events[uid] = event
					Map := targetMap[target]
					delete(Map, uid)
				}
			}
			c.mu.Unlock()
			time.Sleep(10 * time.Second)
		}
	}
}
