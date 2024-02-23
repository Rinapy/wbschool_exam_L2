package netcat

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// Netcat - тип описывающий утилиту
type Netcat struct {
	reader io.Reader
	writer io.Writer
}

// NewNC создает экземпляр Netcat
func NewNC(r io.Reader) *Netcat {
	return &Netcat{
		reader: r,
	}
}

// Run запускает Netcat
func (nc *Netcat) Run(req []string) error {
	if len(req) != 3 { // первый аргумент -t или -u для определения типа протокола, дальше адрес и порт
		return &ArgsError{}
	}
	if req[0] == "-u" {
		return nc.UDP(req[1] + ":" + req[2])
	} else if req[0] == "-t" {
		return nc.TCP(req[1] + ":" + req[2])
	} else {
		return &NoConnectionType{}
	}
}

// TCP создает подключение по TCP
func (nc *Netcat) TCP(url string) error {
	s, err := net.ResolveTCPAddr("tcp4", url)
	if err != nil { // это чтоб localhost в 127.0.0.1 красиво превращался
		return err
	}
	c, errD := net.DialTCP("tcp4", nil, s)
	if err != nil { // создание соединения
		return errD
	}
	defer c.Close() // отложенное закрытие соединения
	nc.writer = c
	if err = nc.session(); err != nil { // передача данных
		return err
	}

	return nil
}

// UDP создает подключение по UDP
func (nc *Netcat) UDP(url string) error {
	s, err := net.ResolveUDPAddr("udp4", url) // тут все аналогично TCP()
	if err != nil {
		return err
	}
	c, errD := net.DialUDP("udp4", nil, s)
	if err != nil {
		return errD
	}
	defer c.Close()
	nc.writer = c
	if err = nc.session(); err != nil {
		return err
	}
	return nil
}

func (nc *Netcat) session() error {
	r := bufio.NewScanner(nc.reader)
	for r.Scan() { // считывание из ридера
		line := r.Text()
		if line == "exit" { // обработка выхода
			break
		}
		if _, err := fmt.Fprint(nc.writer, line+"\n"); err != nil { // отправка данных
			return err
		}
	}
	return nil
}
