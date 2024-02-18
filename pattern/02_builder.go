package main

import "fmt"

/*
	Реализовать паттерн «строитель».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/Builder_pattern
*/

const (
	FerrariCollectorType = "Ferrari"
	ToyotaCollectorType  = "Toyota"
)

type Car struct {
	Brand  string
	Model  string
	Engine string
	Color  string
}

func (c *Car) carInfo() string {
	return fmt.Sprintf("Car: %s, Model: %s, Engine: %s, Color: %s\n", c.Brand, c.Model, c.Engine, c.Color)
}

type Collector interface {
	SetBrand()
	SetModel()
	SetEngine()
	SetColor()
	GetCar() Car
}

func GetCollector(collType string) Collector {
	switch collType {
	default:
		return nil
	case FerrariCollectorType:
		return &FerrariCollector{}
	case ToyotaCollectorType:
		return &ToyotaCollector{}

	}
}

type FerrariCollector struct {
	Brand  string
	Model  string
	Engine string
	Color  string
}

func (f *FerrariCollector) SetBrand() {
	f.Brand = "Ferrari"
}

func (f *FerrariCollector) SetModel() {
	f.Model = "F40"
}

func (f *FerrariCollector) SetEngine() {
	f.Engine = "V12"
}

func (f *FerrariCollector) SetColor() {
	f.Color = "Red"
}

func (f *FerrariCollector) GetCar() Car {
	return Car{
		Brand:  f.Brand,
		Engine: f.Engine,
		Color:  f.Color,
		Model:  f.Model,
	}
}

type ToyotaCollector struct {
	Brand  string
	Model  string
	Engine string
	Color  string
}

func (t *ToyotaCollector) SetBrand() {
	t.Brand = "Toyota"
}

func (t *ToyotaCollector) SetModel() {
	t.Model = "Camry"
}

func (t *ToyotaCollector) SetEngine() {
	t.Engine = "V8"
}

func (t *ToyotaCollector) SetColor() {
	t.Color = "black"
}

func (t *ToyotaCollector) GetCar() Car {
	return Car{
		Brand:  t.Brand,
		Color:  t.Color,
		Engine: t.Engine,
		Model:  t.Model,
	}
}

type Director struct {
	Collector Collector
}

func NewDirector(collector Collector) *Director {
	return &Director{Collector: collector}
}

func (d *Director) SetCollector(collector Collector) {
	d.Collector = collector
}

func (d *Director) BuildCar() Car {
	d.Collector.SetBrand()
	d.Collector.SetColor()
	d.Collector.SetEngine()
	d.Collector.SetModel()
	return d.Collector.GetCar()
}

func main() {
	ToyotaColle := GetCollector("Toyota")
	FerrariColle := GetCollector("Ferrari")

	director := NewDirector(ToyotaColle)
	ToyotaCar := director.BuildCar()
	fmt.Println(ToyotaCar.carInfo())

	director.SetCollector(FerrariColle)
	FerrariCar := director.BuildCar()
	fmt.Println(FerrariCar.carInfo())
}
