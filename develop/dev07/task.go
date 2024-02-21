package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:

	sig := func(after time.Duration) <- chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
	}()

return c
}

start := time.Now()
<-or (

	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),

)

fmt.Printf(“fone after %v”, time.Since(start))
*/
// Функция or принимает несколько каналов для объединения и возвращает объединенный канал
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		// Если не передано ни одного канала, возвращается nil
		return nil
	case 1:
		// Если передан только один канал, он возвращается без изменений
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			// Если передано два канала, используется select для объединения
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			// Если передано более двух каналов, рекурсивно вызывается функция or
			selectCases := make([]<-chan interface{}, len(channels)-1)
			for i := range selectCases {
				selectCases[i] = channels[i+1]
			}
			select {
			case <-channels[0]:
			case <-or(selectCases...):
			}
		}
	}()

	return orDone
}

func main() {
	// Функция sig создает канал, который закроется после заданного времени
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	// Замеряется время начала выполнения
	start := time.Now()

	// Используется or для ожидания закрытия одного из каналов
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(60*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// Выводится время, прошедшее с начала выполнения
	fmt.Printf("Done after %v", time.Since(start))
}
