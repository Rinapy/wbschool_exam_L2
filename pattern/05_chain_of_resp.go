package main

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/
type Handler interface {
	SetNextHandler(next Handler)
	Handle(r Request)
}

type Request struct {
	Type string
	Data interface{}
}

type Handler1 struct {
	NextHandler Handler
}

func (h *Handler1) SetNextHandler(next Handler) {
	h.NextHandler = next
}

func (h *Handler1) Handle(r Request) {
	if r.Type == "Тип1" {
		fmt.Println("Обработка данных по 1 хендлеру")
	} else if h.NextHandler != nil {
		h.NextHandler.Handle(r)
	}
}

type Handler2 struct {
	NextHandler Handler
}

func (h *Handler2) SetNextHandler(next Handler) {
	h.NextHandler = next
}

func (h *Handler2) Handle(r Request) {
	if r.Type == "Тип2" {
		fmt.Println("Обработка данных по 2 хендлеру")
	} else if h.NextHandler != nil {
		h.NextHandler.Handle(r)
	}
}

func main() {

	h1 := &Handler1{}
	h2 := &Handler2{}

	h1.SetNextHandler(h2)

	r := Request{Type: "Тип2", Data: "Данные запроса"}
	h1.Handle(r)
}
