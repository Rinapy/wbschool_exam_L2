package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	errMethod   = &IncorrectMethod{}
	errInput    = &InvalidInput{}
	errServer   = &ServerError{}
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
	fmt.Println(eventRaw)
	if err != nil {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	if !eventRaw.IsValid() {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	event := eventRaw.convertToEvent()
	uid, err := s.calendar.CreateEvent(event)
	if err != nil {
		s.response(false, w, errServer.Error(), http.StatusInternalServerError)
		return
	}
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

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.response(false, w, errMethod.Error(), http.StatusServiceUnavailable)
		return
	}
	//user, foundUser := r.URL.Query()["user_uid"]
	eventRaw, foundEvent := r.URL.Query()["event_uid"]
	eventUID := eventRaw[0]
	if !foundEvent {
		s.response(false, w, errInput.Error(), http.StatusBadRequest)
		return
	}
	event, err := s.calendar.GetEvent(eventUID)
	if err != nil {
		s.response(false, w, err.Error(), http.StatusBadRequest)
		return
	}
	s.response(true, w, event, http.StatusOK)
}
