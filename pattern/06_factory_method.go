package main

import (
	"errors"
	"fmt"
)

/*
	Реализовать паттерн «фабричный метод».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/Factory_method_pattern
*/
const (
	IPhoneType  = "IPhone"
	SamsungType = "Samsung"
)

type Product interface {
	SpendMoneyOnAdvertising()
}

type IPhone struct {
}

func (i *IPhone) SpendMoneyOnAdvertising() {
	fmt.Println("Отправили много денег на рекламу")
}

type Samsung struct {
}

func (s *Samsung) SpendMoneyOnAdvertising() {
	fmt.Println("Отправили немного денег на рекламу")
}

type GodClass struct {
}

func (gc *GodClass) New(product string) (Product, error) {
	switch product {
	default:
		err := errors.New("the transferred product is missing")
		return nil, err
	case IPhoneType:
		return &IPhone{}, nil
	case SamsungType:
		return &Samsung{}, nil
	}
}

func main() {
	gs := GodClass{}

	IPh, err := gs.New("IPhone")
	if err != nil {
		fmt.Println(err.Error())
	}
	IPh.SpendMoneyOnAdvertising()

	Sg, err := gs.New("Samsung")
	if err != nil {
		fmt.Println(err.Error())
	}
	Sg.SpendMoneyOnAdvertising()
}
