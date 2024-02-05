package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

// Комментируем функцию main()
func main() {
	// Инициализация конфигурации
	c := Configure()
	// Создание нового объекта Telnet с передачей конфигурации
	t := NewTelnet(c)
	// Запуск Telnet соединения
	t.Run()
}

type Config struct {
	Timeout time.Duration // Период ожидания для подключения
	Host    string        // Хост (адрес сервера)
	Port    string        // Порт
}

func Configure() *Config {
	c := &Config{}                                               // Создание нового объекта Config
	flag.DurationVar(&c.Timeout, "timeout", 0, "telnet timeout") // Регистрация флага -timeout для периода ожидания
	flag.Parse()                                                 // Парсинг аргументов командной строки
	left := flag.Args()                                          // Получение оставшихся аргументов командной строки
	if len(left) != 2 {                                          // Если количество аргументов не равно 2
		log.Fatalln("incorrect host and port") // Завершение программы с сообщением об ошибке
	}
	c.Host = left[0] // Задание хоста из аргумента командной строки
	c.Port = left[1] // Задание порта из аргумента командной строки
	return c         // Возврат указателя на структуру Config
}

type TelnetUtil struct {
	config  *Config    // Конфигурация
	conn    net.Conn   // Соединение
	errorCh chan error // Канал для передачи ошибок
}

func NewTelnet(config *Config) *TelnetUtil {
	return &TelnetUtil{
		config:  config,           // Задание конфигурации
		errorCh: make(chan error), // Создание канала для ошибок
	}
}

func (t *TelnetUtil) Run() {
	t.connect()                                          // Установка соединения
	sigint := make(chan os.Signal)                       // Создание канала для сигналов
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM) // Подписка на сигналы прерывания
	go t.receive()                                       // Запуск метода приема данных в отдельной горутине
	go t.send()                                          // Запуск метода отправки данных в отдельной горутине
	select {
	// Ожидание канала
	case err := <-t.errorCh: // Если пришла ошибка
		log.Println(err) // Вывести сообщение об ошибке
		t.disconnect()   // Закрыть соединение
		return
	case <-sigint: // Если пришел сигнал прерывания
		t.disconnect() // Закрыть соединение
		return
	}
}

func (t *TelnetUtil) connect() {
	conn, err := net.DialTimeout("tcp4", t.config.Host+":"+t.config.Port, t.config.Timeout) // Попытка установить соединение с сервером и заданными параметрами
	if err != nil {                                                                         // В случае ошибки
		time.Sleep(t.config.Timeout)    // Пауза на протяжении заданного периода
		log.Fatalln("connection error") // Завершение программы с сообщением об ошибке
	}
	t.conn = conn // Задание соединения
}

func (t *TelnetUtil) disconnect() {
	t.conn.Close() // Закрытие соединения
}

func (t *TelnetUtil) receive() {
	r := bufio.NewReader(t.conn) // Создание нового читателя для соединения
	for {                        // Бесконечный цикл
		line, err := r.ReadString('\n') // Чтение строки из соединения до символа новой строки
		if err != nil {                 // Если произошла ошибка (закрытие сокета)
			t.errorCh <- err // Отправка ошибки в канал
			return
		}
		fmt.Print(line) // Вывод строки
	}
}

func (t *TelnetUtil) send() {
	r := bufio.NewReader(os.Stdin) // Создание нового читателя для стандартного ввода
	for {                          // Бесконечный цикл
		line, err := r.ReadString('\n') // Чтение строки из стандартного ввода до символа новой строки
		if err != nil {                 // Если произошла ошибка (нажата комбинация ctrl+d для завершения ввода)
			t.errorCh <- err // Отправка ошибки в канал
			return
		}
		_, err = t.conn.Write([]byte(line)) // Отправка строки через соединение
		if err != nil {                     // Если произошла ошибка при отправке
			t.errorCh <- err // Отправка ошибки в канал
			return
		}
	}
}
