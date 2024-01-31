package server

import (
	"dev11/internal/dev11/calendar"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cfg      *Config
	calendar *calendar.Calendar
}

func NewServer(cfg *Config) *Server {
	return &Server{
		cfg:      cfg,
		calendar: calendar.NewCalendar(),
	}
}

func (s *Server) handle() {
	http.Handle("/create_event", logger(http.HandlerFunc(s.addEvent)))
	http.Handle("/update_event", logger(http.HandlerFunc(s.updateEvent)))
	http.Handle("/delete_event", logger(http.HandlerFunc(s.deleteEvent)))
	http.Handle("/events_for_day", logger(http.HandlerFunc(s.dayEvents)))
	http.Handle("/events_for_week", logger(http.HandlerFunc(s.weekEvents)))
	http.Handle("/events_for_month", logger(http.HandlerFunc(s.monthEvents)))
}

func logger(handler http.Handler) http.Handler { // Работает как middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeBefore := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s\n", r.Method, r.URL, timeBefore, time.Since(timeBefore))
	})
}

func (s *Server) runServer(err chan error) {
	go func() {
		err <- http.ListenAndServe(s.cfg.address, nil) // дефолтный http.DefaultServeMux
	}()
}

// Run запускает сервер
func (s *Server) Run() {
	s.handle()
	sigint := make(chan os.Signal)
	errors := make(chan error)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM) // ctrl+c или kill <pid>
	s.runServer(errors)
	select {
	case <-sigint:
		log.Println("server stopped")
		return
	case err := <-errors:
		log.Println(err)
	}
}
