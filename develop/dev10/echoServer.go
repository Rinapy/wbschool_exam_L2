package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Сервер запущен")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при установлении соединения:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Новое соединение:", conn.RemoteAddr())
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			return
		}
		fmt.Printf("Получено от %s: %s", conn.RemoteAddr(), string(buf[:n]))
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("Ошибка записи:", err)
			return
		}
	}
}
