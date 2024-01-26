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
	k         int
	n, r, u   bool
	inputFile string
)

type Line struct {
	Fields []string
}

type LineSlice []Line

func (ls *LineSlice) Len() int {
	return len(*ls)
}

func (ls *LineSlice) Less(i, j int) bool {
	switch true {
	case n && !r:
		val1, err := strconv.Atoi((*ls)[i].Fields[k])
		val2, err := strconv.Atoi((*ls)[j].Fields[k])
		if err != nil {
			return (*ls)[i].Fields[k] < (*ls)[j].Fields[k]
		}
		return val1 < val2
	case n && r:
		val1, err1 := strconv.Atoi((*ls)[i].Fields[k])
		val2, err2 := strconv.Atoi((*ls)[j].Fields[k])
		// Если оба значения могут быть преобразованы в числа
		if err1 == nil && err2 == nil {
			return val1 > val2 // Сортировка чисел в убывающем порядке
		}
		// Если только одно значение может быть преобразовано в число
		if err1 == nil || err2 == nil {
			return err1 != nil // Поместить текстовые значения выше числовых значений
		}
		// Если оба значения являются текстовыми
		return (*ls)[i].Fields[k] > (*ls)[j].Fields[k] // Сортировка текста в обратном алфавитном порядке
	default:
		return (*ls)[i].Fields[k] > (*ls)[j].Fields[k]
	}
}

func (ls *LineSlice) Swap(i, j int) {
	(*ls)[i], (*ls)[j] = (*ls)[j], (*ls)[i]
}

func (ls *LineSlice) delDuplicate() {
	seenLines := map[string]bool{}
	for i, line := range *ls {
		if !seenLines[line.Fields[k]] {
			seenLines[line.Fields[k]] = true
		} else {
			(*ls)[i] = (*ls)[ls.Len()-1]
			*ls = (*ls)[:ls.Len()-1]
		}
	}
}

func Sorter() {
	flag.IntVar(&k, "k", 0, "Индекс колонки по которой будет сортировка, колонка 1 == 0 2 == 1 и т.д")
	flag.BoolVar(&n, "n", false, "Сортировать по числам, при вхождении в колонку текстовых значений, они будут выноситься вверх")
	flag.BoolVar(&r, "r", false, "Сортировать по убыванию, при вхождении в колонку текстовых значений, они будут выноситься вверх")
	flag.BoolVar(&u, "u", false, "Удаляет дубликат в колонке")
	flag.Parse()
	inputFile = flag.Arg(0)

	fmt.Println(n)
	lines := fillLineSlice(inputFile)
	if u {
		lines.delDuplicate()
	}
	sort.Sort(lines)
	fillNewFile(lines)
	//	for _, line := range *lines {
	//		fmt.Println(&line.Fields)
	//	}
}

func fillNewFile(slice *LineSlice) {
	file, err := os.Create("output.txt")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			log.Fatal(err)
		}
	}(writer)
	for _, v := range *slice {
		_, err := writer.WriteString(strings.Join(v.Fields, " ") + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func fillLineSlice(filename string) *LineSlice {
	file, err := os.Open(filename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	if err != nil {
		log.Fatal(err)
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
	return &lineSlice
}
