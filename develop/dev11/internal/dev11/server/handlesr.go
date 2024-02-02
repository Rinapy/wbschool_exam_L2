package server

import (
	"dev11/internal/dev11/calendar"
	"encoding/json"
	"net/http"
	"time"
)

var (
	errMethod   = &IncorrectMethod{}
	errInput    = &InvalidInput{}
	errServer   = &ServerError{}
	errDate     = &InvalidDate{}
	errEventUID = &InvalidEventUID{}
)

func (s *Server) response(isResult bool, w http.ResponseWriter, payload interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if isResult {
		payload = Result{Result: payload}
	} else {
		payload = Error{Error: payload}
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "response error", http.StatusInternalServerError)
	}

}

func (s *Server) addEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.response(false, w, errMethod.Error(), http.StatusServiceUnavailable)
		return
	}
	eventRaw := NewEventRaw()
	err := json.NewDecoder(r.Body).Decode(&eventRaw)
	if err != nil {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	if !eventRaw.IsValid() {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	event := eventRaw.convertToEvent()
	uid := s.calendar.CreateEvent(event)
	s.response(true, w, struct {
		UID string `json:"uid"`
	}{UID: uid}, http.StatusOK)
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.response(false, w, errMethod.Error(), http.StatusServiceUnavailable)
		return
	}
	eventRawUpd := NewEventUpdateRaw()
	err := json.NewDecoder(r.Body).Decode(&eventRawUpd)
	if err != nil {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	if !eventRawUpd.IsValid() {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	event := eventRawUpd.convertToEvent()
	curEve, err := s.calendar.UpdateEvent(event, eventRawUpd.UID)
	if err != nil {
		s.response(false, w, nil, http.StatusBadRequest)
	} else {
		s.response(true, w, curEve, http.StatusOK)
	}
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.response(false, w, errMethod.Error(), http.StatusServiceUnavailable)
		return
	}
	eveRaw := NewEventUpdateRaw()
	err := json.NewDecoder(r.Body).Decode(&eveRaw)
	if err != nil {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	if !eveRaw.IsValid() {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	err = s.calendar.DeleteEvent(eveRaw.UID)
	if err != nil {
		s.response(false, w, errEventUID.Error(), http.StatusBadRequest)
	} else {
		s.response(true, w, struct {
			Msg string `json:"msg"`
		}{Msg: "Event deleted."}, http.StatusOK)
	}
}

func (s *Server) dayEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.response(false, w, errMethod.Error(), http.StatusServiceUnavailable)
		return
	}
	userRaw, foundUser := r.URL.Query()["user_id"]
	dateRaw, foundDate := r.URL.Query()["date"]
	getEvents := calendar.NewGetEvents()
	if !foundUser {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	getEvents.UserUID = userRaw[0]
	if !foundDate {
		eves := s.calendar.GetDayEvents(getEvents)
		s.response(true, w, eves, http.StatusOK)
		return
	}
	date, err := time.Parse("2006-01-02", dateRaw[0])
	if err != nil {
		s.response(false, w, errDate.Error(), http.StatusBadRequest)
		return
	}
	getEvents.Date = date
	getEvents.DateIn = true
	eves := s.calendar.GetDayEvents(getEvents)
	s.response(true, w, eves, http.StatusOK)
}

func (s *Server) weekEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.response(false, w, errMethod.Error(), http.StatusServiceUnavailable)
		return
	}
	userRaw, foundUser := r.URL.Query()["user_id"]
	dateRaw, foundDate := r.URL.Query()["date"]
	getEvents := calendar.NewGetEvents()
	if !foundUser {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	getEvents.UserUID = userRaw[0]
	if !foundDate {
		eves := s.calendar.GetWeekEvents(getEvents)
		s.response(true, w, eves, http.StatusOK)
		return
	}
	date, err := time.Parse("2006-01-02", dateRaw[0])
	if err != nil {
		s.response(false, w, errDate.Error(), http.StatusBadRequest)
		return
	}
	getEvents.Date = date
	getEvents.DateIn = true
	eves := s.calendar.GetWeekEvents(getEvents)
	s.response(true, w, eves, http.StatusOK)
}
func (s *Server) monthEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.response(false, w, errMethod.Error(), http.StatusServiceUnavailable)
		return
	}
	userRaw, foundUser := r.URL.Query()["user_id"]
	dateRaw, foundDate := r.URL.Query()["date"]
	getEvents := calendar.NewGetEvents()
	if !foundUser {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	getEvents.UserUID = userRaw[0]
	if !foundDate {
		eves := s.calendar.GetMonthEvents(getEvents)
		s.response(true, w, eves, http.StatusOK)
		return
	}
	date, err := time.Parse("2006-01-02", dateRaw[0])
	if err != nil {
		s.response(false, w, errDate.Error(), http.StatusBadRequest)
		return
	}
	getEvents.Date = date
	getEvents.DateIn = true
	eves := s.calendar.GetMonthEvents(getEvents)
	s.response(true, w, eves, http.StatusOK)
}
