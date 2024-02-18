package main

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

type Visitor interface {
	visitTextDocument(doc *TextDocument)
	visitAudioDocument(doc *AudioDocument)
	visitGraphDocument(doc *GraphDocument)
}

type Document interface {
	AcceptVisitor(visitor Visitor)
}

type TextDocument struct {
}

func (doc *TextDocument) AcceptVisitor(visitor Visitor) {
	visitor.visitTextDocument(doc)
}

type AudioDocument struct {
}

func (doc *AudioDocument) AcceptVisitor(visitor Visitor) {
	visitor.visitAudioDocument(doc)
}

type GraphDocument struct {
}

func (doc *GraphDocument) AcceptVisitor(visitor Visitor) {
	visitor.visitGraphDocument(doc)
}

// Конкрентый посититель
type ChangeDocName struct {
}

func (cn *ChangeDocName) visitTextDocument(doc *TextDocument) {
	fmt.Println("Смена имени тестового документа")
}

func (cn *ChangeDocName) visitAudioDocument(doc *AudioDocument) {
	fmt.Println("Смена имени аудио документа")
}

func (cn *ChangeDocName) visitGraphDocument(doc *GraphDocument) {
	fmt.Println("Смена имени графического документа")
}

func main() {
	documents := []Document{
		&TextDocument{},
		&GraphDocument{},
		&AudioDocument{},
	}

	changeName := &ChangeDocName{}

	for _, doc := range documents {
		doc.AcceptVisitor(changeName)
	}
}
