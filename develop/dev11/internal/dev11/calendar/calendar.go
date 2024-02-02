package calendar

import (
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
)

type GetEvents struct {
	UserUID string
	Date    time.Time
	DateIn  bool
}

func NewGetEvents() *GetEvents {
	return &GetEvents{}
}

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

func (c *Calendar) UIDGen() string {
	var key string
	b := uuid.New()
	key = b.String()
	return key
}

func (c *Calendar) GetEvent(UID string) (*Event, error) {
	c.mu.RLock()

	if e, ok := c.events[UID]; ok {
		return e, nil
	}
	if e, ok := c.dayEvents[UID]; ok {
		return e, nil
	}
	if e, ok := c.weekEvents[UID]; ok {
		return e, nil
	}
	if e, ok := c.monthEvents[UID]; ok {
		return e, nil
	}

	return nil, &EventNotFound{}
}

func (c *Calendar) CreateEvent(e *Event) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	UID := c.UIDGen()
	c.events[UID] = e
	return UID

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

	if _, ok := c.events[UID]; ok {
		delete(c.events, UID)
	} else if _, ok = c.dayEvents[UID]; ok {
		delete(c.dayEvents, UID)
	} else if _, ok = c.weekEvents[UID]; ok {
		delete(c.weekEvents, UID)
	} else if _, ok = c.monthEvents[UID]; ok {
		delete(c.monthEvents, UID)
	} else {
		return &EventNotFound{}
	}
	return nil
}

func (c *Calendar) GetDayEvents(getEvents *GetEvents) []*Event {
	res := make([]*Event, 0)
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !getEvents.DateIn {
		for _, event := range c.dayEvents {
			if event.UserUID == getEvents.UserUID {
				res = append(res, event)
			}
		}
		return res
	}
	for _, v := range c.events {
		if v.UserUID == getEvents.UserUID {
			if middle(v.StartDT.Year(), getEvents.Date.Year(), v.EndDT.Year()) &&
				middle(int(v.StartDT.Month()), int(getEvents.Date.Month()), int(v.EndDT.Month())) &&
				middle(v.StartDT.Day(), getEvents.Date.Day(), v.EndDT.Day()) {
				res = append(res, v)
			}
		}
	}
	for _, v := range c.dayEvents {
		if v.UserUID == getEvents.UserUID {
			if middle(v.StartDT.Year(), getEvents.Date.Year(), v.EndDT.Year()) &&
				middle(int(v.StartDT.Month()), int(getEvents.Date.Month()), int(v.EndDT.Month())) &&
				middle(v.StartDT.Day(), getEvents.Date.Day(), v.EndDT.Day()) {
				res = append(res, v)
			}
		}
	}
	return res
}

func (c *Calendar) GetWeekEvents(getEvents *GetEvents) []*Event {
	res := make([]*Event, 0)
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !getEvents.DateIn {
		for _, event := range c.weekEvents {
			if event.UserUID == getEvents.UserUID {
				res = append(res, event)
			}
		}
		return res
	}
	var y1, w1, yx, wx, y2, w2 int
	for _, v := range c.events {
		if v.UserUID == getEvents.UserUID {
			y1, w1 = v.StartDT.ISOWeek()
			yx, wx = getEvents.Date.ISOWeek()
			y2, w2 = v.EndDT.ISOWeek()
			if middle(y1, yx, y2) && middle(w1, wx, w2) {
				res = append(res, v)
			}
		}
	}
	for _, v := range c.weekEvents {
		if v.UserUID == getEvents.UserUID {
			y1, w1 = v.StartDT.ISOWeek()
			yx, wx = getEvents.Date.ISOWeek()
			y2, w2 = v.EndDT.ISOWeek()
			if middle(y1, yx, y2) && middle(w1, wx, w2) {
				res = append(res, v)
			}
		}
	}
	return res
}

func (c *Calendar) GetMonthEvents(getEvents *GetEvents) []*Event {
	res := make([]*Event, 0)
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !getEvents.DateIn {
		for _, event := range c.monthEvents {
			if event.UserUID == getEvents.UserUID {
				res = append(res, event)
			}
		}
		return res
	}
	for _, v := range c.events {
		if v.UserUID == getEvents.UserUID {
			if middle(v.StartDT.Year(), getEvents.Date.Year(), v.EndDT.Year()) &&
				middle(int(v.StartDT.Month()), int(getEvents.Date.Month()), int(v.EndDT.Month())) {
				res = append(res, v)
			}
		}
	}
	for _, v := range c.monthEvents {
		if v.UserUID == getEvents.UserUID {
			if middle(v.StartDT.Year(), getEvents.Date.Year(), v.EndDT.Year()) &&
				middle(int(v.StartDT.Month()), int(getEvents.Date.Month()), int(v.EndDT.Month())) {
				res = append(res, v)
			}
		}
	}
	return res
}

func (c *Calendar) EventBalanced(sigint <-chan os.Signal) {
	for {
		select {
		case <-sigint:
			return
		default:
			c.mu.Lock()
			for uid, event := range c.events {
				added := false
				if c.inDayRange(event) {
					//log.Printf("EventsBalanced add event to dayEvents\n")
					c.dayEvents[uid] = event
					added = true
				}
				if c.inWeekRange(event) {
					//log.Printf("EventsBalanced add event to weekEvents\n")
					c.weekEvents[uid] = event
					added = true
				}
				if c.inMonthRange(event) {
					//log.Printf("EventsBalanced add event to monthEvents\n")
					c.monthEvents[uid] = event
					added = true
				}
				if added {
					//log.Printf("EventsBalanced del event to events")
					delete(c.events, uid)
				}

			}
			c.mu.Unlock()
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *Calendar) EventsCleaner(sigint <-chan os.Signal) {
	for {
		select {
		case <-sigint:
			return
		default:
			c.mu.Lock()
			for uid, event := range c.dayEvents {
				if !c.inDayRange(event) {
					delete(c.dayEvents, uid)
					//log.Printf("EventsCleaner del event %v to dayEvents", event.Name)
					_, InWeek := c.weekEvents[uid]
					_, InMonth := c.monthEvents[uid]
					if !InWeek && !InMonth {
						c.events[uid] = event
						//log.Printf("EventsCleaner add event %v to events", event.Name)
					}
				}
			}
			for uid, event := range c.weekEvents {
				if !c.inWeekRange(event) {
					delete(c.weekEvents, uid)
					//log.Printf("EventsCleaner del event %v to weekEvents", event.Name)
					_, InMonth := c.monthEvents[uid]
					if !InMonth {
						c.events[uid] = event
						//log.Printf("EventsCleaner add event %v to events", event.Name)
					}
				}
			}
			for uid, event := range c.monthEvents {
				if !c.inMonthRange(event) {
					delete(c.monthEvents, uid)
					//log.Printf("EventsCleaner del event %v to monthEvents", event.Name)
					c.events[uid] = event
					//log.Printf("EventsCleaner add event %v to events", event.Name)
				}
			}
		}
		c.mu.Unlock()
		//fmt.Println("events", c.events)
		//fmt.Println("dayEvents", c.dayEvents)
		//fmt.Println("weekEvents", c.weekEvents)
		//fmt.Println("monthEvents", c.monthEvents)
		time.Sleep(60 * time.Second)
	}
}

func (c *Calendar) inDayRange(event *Event) bool {
	today := time.Now().Truncate(24 * time.Hour)
	return middle(event.StartDT.Year(), today.Year(), event.EndDT.Year()) &&
		middle(int(event.StartDT.Month()), int(today.Month()), int(event.EndDT.Month())) &&
		middle(event.StartDT.Day(), today.Day(), event.EndDT.Day())
}
func (c *Calendar) inWeekRange(event *Event) bool {
	today := time.Now().Truncate(24 * time.Hour)
	y1, w1 := event.StartDT.ISOWeek()
	yx, wx := today.ISOWeek()
	y2, w2 := event.EndDT.ISOWeek()
	return middle(y1, yx, y2) && middle(w1, wx, w2)
}
func (c *Calendar) inMonthRange(event *Event) bool {
	today := time.Now().Truncate(24 * time.Hour)
	thisMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	return middle(event.StartDT.Year(), thisMonth.Year(), event.EndDT.Year()) &&
		middle(int(event.StartDT.Month()), int(thisMonth.Month()), int(event.EndDT.Month()))
}

func middle(first, x, second int) bool {
	if first <= x && x <= second {
		return true
	}
	return false
}
