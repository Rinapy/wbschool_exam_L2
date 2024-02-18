package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

type State interface {
	Handle()
}

type StandingState struct{}

func (s *StandingState) Handle() {
	fmt.Println("Персонаж стоит на месте")
}

type WalkingState struct{}

func (s *WalkingState) Handle() {
	fmt.Println("Персонаж идет")
}

type RunningState struct{}

func (s *RunningState) Handle() {
	fmt.Println("Персонаж бежит")
}

type Context struct {
	currentState State
}

func (c *Context) ChangeState(state State) {
	c.currentState = state
}

func (c *Context) Handle() {
	c.currentState.Handle()
}

func main() {
	context := &Context{}
	walkingState := &WalkingState{}
	standingState := &StandingState{}
	runningState := &RunningState{}

	// Открытие устройства клавиатуры
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Отслеживание нажатий кнопок
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		switch key {
		default:
			fmt.Println("Данная клавиша не управляет персонажем")
		case keyboard.KeyArrowUp:
			context.ChangeState(runningState)
			context.Handle()
		case keyboard.KeyArrowDown:
			context.ChangeState(standingState)
			context.Handle()
		case keyboard.KeyArrowRight:
			context.ChangeState(walkingState)
			context.Handle()
		}

	}
}
