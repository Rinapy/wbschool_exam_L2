package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	A               int
	B               int
	C               int
	c               bool
	i               bool
	v               bool
	F               bool
	n               bool
	searchStr       string
	nameOrPattern   string
	ErrFindFiles    = errors.New("ошибка поиска файла или файлов")
	ErrNoLinesFound = errors.New("совпадения не найдены")
	ErrFileIsEmpty  = errors.New("файл пустой")
)

type Line struct {
	Text string
}
type lineSlice []Line

type File struct {
	name    string
	lines   lineSlice
	findStr []int
}

type FileSlice []File

func parseFlags() {
	flag.IntVar(&A, "A", 0, "\"after\" печатать +N строк после совпадения")
	flag.IntVar(&B, "B", 0, "\"before\" печатать +N строк до совпадения")
	flag.IntVar(&C, "C", 0, "\"context\" (A+B) печатать ±N строк вокруг совпадения")
	flag.BoolVar(&c, "c", false, "\"count\" (количество строк)")
	flag.BoolVar(&i, "i", false, "\"ignore-case\" (игнорировать регистр)")
	flag.BoolVar(&v, "v", false, "\"invert\" (вместо совпадения, исключать)")
	flag.BoolVar(&F, "F", false, "\"fixed\", точное совпадение со строкой, не паттерн")
	flag.BoolVar(&n, "n", false, "\"line num\", печатать номер строки")
	flag.Parse()
	nameOrPattern = flag.Arg(1)
	searchStr = flag.Arg(0)
}

func parseNameOrPattern(arg string) ([]string, error) {
	fs := make([]string, 0)
	isPattern := false
	for _, val := range arg {
		if string(val) == "*" {
			isPattern = true
			break
		}
	}
	if isPattern {
		files, err := filepath.Glob(arg)
		if len(files) == 0 {
			return nil, ErrFindFiles
		}
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			fs = append(fs, file)
		}
	} else {
		if _, err := os.Stat(arg); err != nil {
			if os.IsNotExist(err) {
				return nil, ErrFindFiles
			}
		} else {
			fs = append(fs, arg)
		}
	}
	return fs, nil
}

func fillFileSlice(filenames []string) (FileSlice, error) {
	fileSlice := FileSlice{}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		lines, err := readLines(file)
		if err != nil && !errors.Is(err, ErrFileIsEmpty) {
			return nil, err
		} else if errors.Is(err, ErrFileIsEmpty) {
			log.Println(filename, err)
		}

		fileInfo := File{
			name:    filename,
			lines:   lines,
			findStr: []int{},
		}

		fileSlice = append(fileSlice, fileInfo)
	}

	return fileSlice, nil
}

func readLines(file *os.File) ([]Line, error) {
	reader := bufio.NewReader(file)
	var lines lineSlice

	for {
		lineText, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		cleanStr := strings.TrimRight(lineText, "\r\n")
		line := Line{Text: cleanStr}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return nil, ErrFileIsEmpty
	}
	return lines, nil
}

func printer(fileSlice FileSlice) {
	match := "Match -- %v\n"
	before := "Before -- %v\n"
	after := "After  -- %v\n"
	matchN := "Match -- %v [%v]\n"
	beforeN := "Before -- %v [%v]\n"
	afterN := "After  -- %v [%v]\n"
	spliter := "---------------------------%v---------------------------\n"

	for _, file := range fileSlice {
		if len(file.findStr) > 0 {
			fmt.Printf(spliter, file.name)
		}

		if c {
			fmt.Printf(match, len(file.findStr))
			continue
		}

		for _, idx := range file.findStr {
			if C > 0 {
				A = C
				B = C
			}

			if B > 0 {
				for iB := 1; iB <= B; iB++ {
					beforeIdx := idx - iB
					if beforeIdx >= 0 {
						if n {
							fmt.Printf(beforeN, file.lines[beforeIdx].Text, beforeIdx+1)
							continue
						}
						fmt.Printf(before, file.lines[beforeIdx].Text)
					}
				}
			}

			if n {
				fmt.Printf(matchN, file.lines[idx].Text, idx+1)
			} else {
				fmt.Printf(match, file.lines[idx].Text)
			}

			if A > 0 {
				for iA := 1; iA <= A; iA++ {
					afterIdx := idx + iA
					if afterIdx < len(file.lines) {
						if n {
							fmt.Printf(afterN, file.lines[afterIdx].Text, afterIdx+1)
							continue
						}
						fmt.Printf(after, file.lines[afterIdx].Text)
					}
				}
			}
		}
	}
}

func finder(fileSlice FileSlice, text string) (FileSlice, error) {
	found := false
	for idx := range fileSlice {
		file := &fileSlice[idx]
		for idx, line := range file.lines {
			if !F {
				if !v {
					if (i && strings.Contains(strings.ToLower(line.Text), strings.ToLower(text))) ||
						(!i && strings.Contains(line.Text, text)) {
						file.findStr = append(file.findStr, idx)
						found = true
					}
				} else {
					if A != 0 || B != 0 || C != 0 {
						A = 0
						B = 0
						C = 0
						fmt.Println("Использование флагов -A -B -C при флаге -v недопустимо, данные флаги отключены.")
					}

					if (i && !strings.Contains(strings.ToLower(line.Text), strings.ToLower(text))) ||
						(!i && !strings.Contains(line.Text, text)) {
						file.findStr = append(file.findStr, idx)
						found = true
					}
				}
			} else {
				if !v {
					if (i && strings.ToLower(line.Text) == strings.ToLower(text)) ||
						(!i && line.Text == text) {
						file.findStr = append(file.findStr, idx)
						found = true

					}
				} else {
					if A != 0 || B != 0 || C != 0 {
						A = 0
						B = 0
						C = 0
						fmt.Println("Использование флагов -A -B -C при флаге -v недопустимо, данные флаги отключены.")
					}

					if (i && strings.ToLower(line.Text) != strings.ToLower(text)) ||
						(!i && line.Text != text) {
						file.findStr = append(file.findStr, idx)
						found = true
					}
				}
			}

		}
	}
	if found {
		return fileSlice, nil
	}
	return nil, ErrNoLinesFound
}
func Find() {
	parseFlags()
	fsn, err := parseNameOrPattern(nameOrPattern)
	if err != nil {
		log.Fatal(err)
	}
	fs, err := fillFileSlice(fsn)
	if err != nil {
		log.Fatal(err)
	}
	fs, err = finder(fs, searchStr)
	if err != nil {
		fmt.Println("Совпадения не найдены")
	} else {
		printer(fs)
	}
}

func main() {
	Find()
}
