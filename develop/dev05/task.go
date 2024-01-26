package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

func main() {
	// Определение флагов командной строки
	after := flag.Int("A", 0, "печатать +N строк после совпадения")
	before := flag.Int("B", 0, "печатать +N строк до совпадения")
	context := flag.Int("C", 0, "печатать ±N строк вокруг совпадения")
	count := flag.Bool("c", false, "количество строк")
	ignoreCase := flag.Bool("i", false, "игнорировать регистр")
	invert := flag.Bool("v", false, "вместо совпадения, исключать")
	fixed := flag.Bool("F", false, "точное совпадение со строкой, не паттерн")
	lineNum := flag.Bool("n", false, "напечатать номер строки")
	// Парсинг флагов командной строки
	flag.Parse()

	// Получение паттерна и файлов из аргументов командной строки
	pattern := flag.Arg(0)
	files := flag.Args()[1:]

	// Проверка наличия паттерна и файлов
	if pattern == "" || len(files) == 0 {
		fmt.Println("Использование: grep [опции] паттерн файлы...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Открытие файлов и фильтрация строк
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Ошибка при открытии файла %s: %v\n", file, err)
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 0
		printLine := func(line string) {
			if *lineNum {
				fmt.Printf("%d:", lineNum)
			}
			fmt.Println(line)
		}

		// Флаги "before" и "context"
		beforeCounter := 0
		contextCounter := 0
		beforeLines := []string{}
		contextLines := []string{}

		for scanner.Scan() {
			line := scanner.Text()
			lineNum++

			match := false
			if *ignoreCase {
				match = strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
			} else if *fixed {
				match = line == pattern
			} else {
				match = strings.Contains(line, pattern)
			}

			if *invert {
				match = !match
			}

			// Флаг "before"
			if *before > 0 && beforeCounter > 0 {
				beforeLines = append(beforeLines, line)
				if len(beforeLines) > *before {
					beforeLines = beforeLines[1:]
				}
			}

			// Флаг "context"
			if *context > 0 && contextCounter > 0 {
				contextLines = append(contextLines, line)
				if len(contextLines) > *context {
					contextLines = contextLines[1:]
				}
			}

			// Печать совпадающих строк
			if match {
				if *before > 0 {
					for _, beforeLine := range beforeLines {
						printLine(beforeLine)
					}
				}

				if *context > 0 {
					for _, contextLine := range contextLines {
						printLine(contextLine)
					}
				}

				printLine(line)

				if *after > 0 {
					afterCounter := 0
					for scanner.Scan() {
						afterLine := scanner.Text()
						lineNum++
						if *lineNum {
							fmt.Printf("%d:", lineNum)
						}
						fmt.Println(afterLine)
						afterCounter++
						if afterCounter >= *after {
							break
						}
					}
				}

				// Флаг "count"
				if *count {
					break
				}
			}

			// Сброс счётчиков для опций "before" и "context"
			if match {
				beforeCounter = 0
				contextCounter = 0
				beforeLines = []string{}
				contextLines = []string{}
			} else {
				beforeCounter++
				contextCounter++
				if *before > 0 && beforeCounter > *before {
					beforeCounter = *before
				}
				if *context > 0 && contextCounter > *context {
					contextCounter = *context
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Ошибка при чтении файла %s: %v\n", file, err)
		}
	}
}
