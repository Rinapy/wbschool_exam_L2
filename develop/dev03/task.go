package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	k int
	n bool
)

type Line struct {
	Fields []string
}

type LineSlice []Line

func (ls LineSlice) Len() int {
	return len(ls)
}

func (ls LineSlice) Less(i, j int) bool {
	if n {
		val1, err := strconv.Atoi(ls[i].Fields[k])
		val2, err := strconv.Atoi(ls[j].Fields[k])
		if err != nil {
			return ls[i].Fields[k] < ls[j].Fields[k]
		}
		return val1 < val2
	}
	return ls[i].Fields[k] < ls[j].Fields[k]
}

func (ls LineSlice) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func main() {
	flag.IntVar(&k, "k", 0, "Индекс колнки по которой будет сортировка")
	flag.BoolVar(&n, "n", false, "Индекс колнки по которой будет сортировка")
	flag.Parse()
	fmt.Println(k)
	fmt.Println(n)
	fileName := "test.txt"
	lines := fillLineSlice(fileName)
	sort.Sort(lines)
	for _, line := range lines {
		fmt.Println(line.Fields)
	}
}
func fillLineSlice(filename string) LineSlice {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	reader := bufio.NewReader(file)
	lineSlice := make(LineSlice, 0)
	for {
		lineText, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
		}

		line := Line{strings.Fields(lineText)}
		lineSlice = append(lineSlice, line)
	}
	return lineSlice
}
