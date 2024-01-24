package main

import (
	"errors"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// ErrIncorrectString описывает ошибку, возникающую при передачи некорректной строки.
var ErrIncorrectString = errors.New("incorrect string passed")

/*
UnzipStr функция производит распаковку переданной строки Примеры:
Примеры работы функции:
- "a4bc2d5e" => "aaaabccddddde", nil
- "abcd" => "abcd", nil
- "45" => "", ErrIncorrectString
- "" => "", nil
- qwe\4\5 => "qwe45", nil
- qwe\45 => "qwe44444", nil
- qwe\\5 => "qwe\\\\\", nil
*/
func UnzipStr(s string) (string, error) {
	var sb strings.Builder
	var prevRune rune
	var isEscaped bool
	backSlash := '\\'

	for i, r := range s {
		isDigit := unicode.IsDigit(r)
		if isDigit && i == 0 {
			return "", ErrIncorrectString
		} // Если первая в строке руна цифра - кидаемся ошибкой
		isBackSlash := r == backSlash
		isUnescapedPrevDigit := unicode.IsDigit(prevRune) && !isEscaped // Если встречаем цифру которая не была экранирована ставим false
		twoUnescapedDigitsInRow := isDigit && isUnescapedPrevDigit      // Если встречаем цифру и прошлая цифра не была экранирована делаем True

		if twoUnescapedDigitsInRow {
			return "", ErrIncorrectString
		} // Если две цифры в строке и первая без экрана кидаемся ошибкой

		// Обработаем слеш
		if isBackSlash && prevRune != backSlash {
			prevRune = r
			continue
		} //Если текущая руна слеш и прошлая руна не слеш обновляем prevRune и идём дальше

		if isBackSlash && prevRune == backSlash {
			sb.WriteRune(r)
			prevRune = -1
			continue
		} // Если два слеша подряд обнуляем prevRune чтобы не попасть на это условием в след цикле и записываем этот слеш ибо экранирован

		if isDigit && prevRune == backSlash {
			sb.WriteRune(r)
			prevRune = r
			isEscaped = true
			continue
		} // Если текущая руна цифра и прошлая слеш пишем цифру и ставим isEscaped в true чтобы не схватить ошибку

		if !isDigit {
			sb.WriteRune(r)
		} // Если текущая руна не цифра пишем её и идём дальше - выражения выше гарантируют что мы не встретим тут слеш

		if isDigit {
			count := int(r - '1')
			if prevRune == -1 {
				prevRune = backSlash
			} // Если прошлая руна -1 то был слеш, возвращаяем его чтобы распакавать в дальнейшем
			for i := 0; i < count; i++ {
				sb.WriteRune(prevRune)
			}
		}
		prevRune = r
		isEscaped = false
	}
	return sb.String(), nil
}
