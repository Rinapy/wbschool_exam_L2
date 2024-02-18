package main

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

type Storage struct {
	data map[string]struct{}
}

func (s *Storage) availabilityCheck(item string) bool {
	_, ok := s.data[item]
	fmt.Println("Проверяем наличие на складе!")
	if ok {
		return true
	}
	return false
}

type Package struct {
}

func (p *Package) packageItem(item string) bool {
	fmt.Printf("Упакавали %s.\n", item)
	return true
}

type Delivery struct {
}

func (d *Delivery) Sending() bool {
	fmt.Println("Отправили заказ")
	return true
}

type Oredering struct {
	storage  *Storage
	pack     *Package
	delivery *Delivery
}

func (o *Oredering) NewOrder(item string) {
	availability := o.storage.availabilityCheck(item)
	if availability == false {
		fmt.Println("К сожадению данного товара нету на складе(")
		return
	}
	packingResult := o.pack.packageItem(item)
	if packingResult == false {
		fmt.Println("Произошла ошибка на моменте упаковки, попробуйте сделать заказ позже.")
		return
	}
	deliverRes := o.delivery.Sending()
	if deliverRes == false {
		fmt.Println("Произошла ошибка на моменте отправки вашего товара.")
		return
	}

}

func main() {

	storageData := make(map[string]struct{})

	storageData["Шкаф"] = struct{}{}
	storageData["Стул"] = struct{}{}

	s := Storage{storageData}

	order := Oredering{
		storage:  &s,
		delivery: &Delivery{},
		pack:     &Package{},
	}

	order.NewOrder("Шкаф")
	order.NewOrder("Стул")
	order.NewOrder("Ложка")
}
