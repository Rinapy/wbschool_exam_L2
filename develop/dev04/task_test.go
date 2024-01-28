package main

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	// Тестовые данные
	words := []string{"cat", "dog", "act", "god", "tac", "good"}

	// Ожидаемый результат
	expected := map[string][]string{
		"act": {"cat", "act", "tac"},
		"dgo": {"dog", "god"},
	}

	// Вызов функции для получения результата
	result := FindAnagrams(words)

	// Проверка соответствия ожидаемого результата и полученного результата
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, result)
	}
}

func TestFindAnagrams_OneWord(t *testing.T) {
	// Тестовые данные с одним словом
	words := []string{"hello"}

	// Ожидаемый результат - пустая мапа
	expected := map[string][]string{}

	// Вызов функции для получения результата
	result := FindAnagrams(words)

	// Проверка соответствия ожидаемого результата и полученного результата
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, result)
	}
}

func TestFindAnagrams_NoAnagrams(t *testing.T) {
	// Тестовые данные без анаграмм
	words := []string{"apple", "pear", "banana"}

	// Ожидаемый результат - пустая мапа
	expected := map[string][]string{}

	// Вызов функции для получения результата
	result := FindAnagrams(words)

	// Проверка соответствия ожидаемого результата и полученного результата
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, result)
	}
}
