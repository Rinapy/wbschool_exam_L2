package main

import (
	"fmt"
	"time"
)

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

type Command interface {
	RunCommand()
}

type SendCommand struct {
	recipient *Recipient
}

func (sc *SendCommand) RunCommand() {
	sc.recipient.SendInfo()
}

type AcceptCommand struct {
	recipient *Recipient
}

func (ac *AcceptCommand) RunCommand() {
	ac.recipient.AcceptInfo()
}

type Recipient struct {
}

func (r *Recipient) SendInfo() {
	fmt.Println("Отправил информацию")
}

func (r *Recipient) AcceptInfo() {
	fmt.Println("Принял информацию")
}

type Initiator struct {
	commands []Command
}

func (i *Initiator) SetCommand(command Command) {
	i.commands = append(i.commands, command)
}

func (i *Initiator) RunCommands() {
	for _, command := range i.commands {
		command.RunCommand()
	}
}

func main() {
	rc := &Recipient{}
	senC := &SendCommand{rc}
	accC := &AcceptCommand{rc}
	inir := &Initiator{}

	inir.SetCommand(senC)
	inir.SetCommand(accC)
	time.Sleep(6 * time.Second)
	inir.RunCommands()
}
