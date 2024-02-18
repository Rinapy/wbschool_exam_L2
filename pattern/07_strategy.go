package main

import "fmt"

/*
	Реализовать паттерн «стратегия».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/Strategy_pattern
*/
type SortStrategy interface {
	Sort(nums []int) []int
}

type Sorter struct {
	strategy SortStrategy
}

func (s *Sorter) SetStrategy(strategy SortStrategy) {
	s.strategy = strategy
}

func (s *Sorter) SortNumbers(nums []int) []int {
	return s.strategy.Sort(nums)
}

type QuickSort struct {
	arr []int
}

func (q *QuickSort) Sort(nums []int) []int {
	q.arr = append(q.arr, nums...)
	low := 0
	high := len(q.arr) - 1
	q.quickSort(low, high)
	return q.arr
}

func (q *QuickSort) quickSort(low int, high int) {

	if low < high {
		// Выбираем опорный элемент
		pivotIndex := q.partition(low, high)

		// Рекурсивно сортируем подмассивы до и после опорного элемента
		q.quickSort(low, pivotIndex-1)
		q.quickSort(pivotIndex+1, high)
	}
}

func (q *QuickSort) partition(low, high int) int {
	// Используем последний элемент в качестве опорного
	pivot := q.arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		// Если текущий элемент меньше или равен опорному, меняем их местами
		if q.arr[j] <= pivot {
			i++
			q.arr[i], q.arr[j] = q.arr[j], q.arr[i]
		}
	}

	// Помещаем опорный элемент в его окончательную позицию
	q.arr[i+1], q.arr[high] = q.arr[high], q.arr[i+1]

	return i + 1
}

type QuickSortReversed struct {
	arr []int
}

func (q *QuickSortReversed) Sort(nums []int) []int {
	q.arr = append(q.arr, nums...)
	low := 0
	high := len(q.arr) - 1
	q.quickSort(low, high)
	return q.arr
}

func (q *QuickSortReversed) quickSort(low int, high int) {

	if low < high {
		// Выбираем опорный элемент
		pivotIndex := q.partition(low, high)

		// Рекурсивно сортируем подмассивы до и после опорного элемента
		q.quickSort(low, pivotIndex-1)
		q.quickSort(pivotIndex+1, high)
	}
}

func (q *QuickSortReversed) partition(low, high int) int {
	// Используем последний элемент в качестве опорного
	pivot := q.arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		// Если текущий элемент больше или равен опорному, меняем их местами
		if q.arr[j] >= pivot {
			i++
			q.arr[i], q.arr[j] = q.arr[j], q.arr[i]
		}
	}

	// Помещаем опорный элемент в его окончательную позицию
	q.arr[i+1], q.arr[high] = q.arr[high], q.arr[i+1]

	return i + 1
}

func main() {
	sorter := &Sorter{}
	numbers := []int{5, 2, 8, 1, 9}
	sorter.SetStrategy(&QuickSort{})
	up := sorter.SortNumbers(numbers)
	fmt.Println(up)
	sorter.SetStrategy(&QuickSortReversed{})
	down := sorter.SortNumbers(numbers)
	fmt.Println(down)
	fmt.Println("Исходный массив не изменён", numbers)
}
