package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type EventData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserUID     string `json:"user_uid"`
	StartDt     string `json:"start_dt"`
	EndDt       string `json:"end_dt"`
}

func TestAddEventHandler(t *testing.T) {
	cfg := DefaultCfg()
	server := NewServer(cfg)
	eventData := getEventDataMap()
	tests := []struct {
		name     string
		method   string
		data     EventData
		wantCode int
	}{
		{
			name:     "valid data POST to /create_event endpoint",
			method:   "POST",
			data:     eventData["Event01Correct"],
			wantCode: http.StatusOK,
		},
		{
			name:     "valid data GET to /create_event endpoint",
			method:   "GET",
			data:     eventData["Event01Correct"],
			wantCode: http.StatusServiceUnavailable,
		},
		{
			name:     "invalid data POST to /create_event endpoint",
			method:   "POST",
			data:     eventData["Event02Incorrect"],
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid JSON POST to /create_event endpoint",
			method:   "POST",
			data:     eventData["Event03IncorrectJSON"],
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		fmt.Println(tt.name)
		jsonData, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(tt.method, "/create_event", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := logger(http.HandlerFunc(server.addEvent))

		handler.ServeHTTP(rr, req)
		if rr.Code != tt.wantCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, tt.wantCode)
		}
	}
}

//func TestUpdateEventHandler_ValidRequest(t *testing.T) {
//	server := &Server{calendar: calendar.NewCalendar()}
//	event := &calendar.Event{
//		Name:        "Event01",
//		Description: "Desc a Event01",
//		UserUID:     "1",
//		StartDT:     time.Now(),
//		EndDT:       time.Now()}
//	uid := server.calendar.CreateEvent(event)
//	eventRawUpd := &EventRawUpdate{
//		UID:         uid,
//		UserUID:     "1",
//		Name:        "Event01",
//		Description: "Desc a updated Event01",
//		StartDT:     "2024-01-01",
//		EndDT:       "2024-01-01",
//	}
//	//Преобразование структуры в JSON
//	payload, _ := json.Marshal(eventRawUpd)
//
//	//// Создание тестового запроса
//	req, err := http.NewRequest("POST", "/update_event", bytes.NewBuffer(payload))
//	if err != nil {
//		t.Fatal(err)
//	}
//	//
//	//// Тестовый ответ
//	rr := httptest.NewRecorder()
//	//
//	//// Вызов обработчика
//	handler := http.HandlerFunc(server.updateEvent)
//	handler.ServeHTTP(rr, req)
//	//
//	//// Проверки
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("Handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//	//
//	//// Проверка результата
//	expectedResult, _ := json.Marshal(event)
//	if rr.Body.String() != string(expectedResult) {
//		t.Errorf("Handler returned unexpected body: got %v want %v",
//			rr.Body.String(), string(expectedResult))
//	}
//}

func getEventDataMap() map[string]EventData {
	dataMap := make(map[string]EventData)
	dataMap["Event01Correct"] = EventData{
		Name:        "Event01",
		Description: "Desc a Event01",
		UserUID:     "1",
		StartDt:     "2024-01-01",
		EndDt:       "2024-01-02",
	}

	dataMap["Event02IncorrectData"] = EventData{
		Name:        "Event02",
		Description: "Desc a Event02",
		UserUID:     "2",
		StartDt:     "2024-02-1",
		EndDt:       "2024-02-02",
	}

	dataMap["Event03IncorrectJSON"] = EventData{
		Name:        "Event03",
		Description: "Desc a Event03",
		StartDt:     "2024-03-01",
		EndDt:       "2024-03-02",
	}

	dataMap["Event04"] = EventData{
		Name:        "Event04",
		Description: "Desc a Event04",
		UserUID:     "4",
		StartDt:     "2024-04-01",
		EndDt:       "2024-04-02",
	}

	dataMap["Event05"] = EventData{
		Name:        "Event05",
		Description: "Desc a Event05",
		UserUID:     "5",
		StartDt:     "2024-05-01",
		EndDt:       "2024-05-02",
	}
	return dataMap
}
