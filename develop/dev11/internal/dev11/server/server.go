package server

import (
	"dev11/internal/dev11/calendar"
	"fmt"
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
		logString := fmt.Sprintf("[Method]%s -- [URL]%s -- %s --[RTime] %dnano\n", r.Method, r.URL, timeBefore, time.Since(timeBefore).Nanoseconds())
		log.Print(logString)
	})
}

func (s *Server) runServer(err chan error) {
	go func() {
		err <- http.ListenAndServe(s.cfg.address, nil) // дефолтный http.DefaultServeMux
	}()
}

// Run запускает сервер
func (s *Server) Run() chan os.Signal {
	s.handle()
	sigint := make(chan os.Signal)
	errors := make(chan error)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM) // ctrl+c или kill <pid>
	go s.calendar.EventBalanced(sigint)
	log.Println("Балансировщик -- Running")
	go s.calendar.EventsCleaner(sigint)
	log.Println("Клинер мап -- Running")
	s.runServer(errors)
	log.Printf("Server start to address -- http://%v\n", s.cfg.address)
	select {
	case <-sigint:
		log.Println("server stopped")
		return nil
	case err := <-errors:
		log.Println(err)
	}
	return sigint
}
