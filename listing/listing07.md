Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select { 
			case v := <-a: // v, ok := <- a вместе с обработчиком ok решило бы вопрос c нулями
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
1, 2, 3, 4, 5, 6, 7, 8, 0 ...

Числа возвращаются последовательно благодаря задержке но не всегда.
Нули возвращаються из-за чтения с закрытого канала, именно 0 из-за типа данного канала
Нули возвращаються из-за чтения с закрытого канала, именно 0 из-за типа данного канала
```
