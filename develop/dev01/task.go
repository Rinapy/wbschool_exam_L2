package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"os"
	"time"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	ntpTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		log.Printf("error get time :%v\n", err)
		os.Exit(1)
	}

	currentTime := time.Now()
	fmt.Println("Текущее время на ПК: ", currentTime.Format(time.RFC3339))
	fmt.Println("Точное время: ", ntpTime.Format(time.RFC3339))

}
