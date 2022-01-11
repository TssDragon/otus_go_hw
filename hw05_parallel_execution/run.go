package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Воркер для выполнения задач.
func worker(taskChannel <-chan Task, wg *sync.WaitGroup, errorsCountDown *int32) {
	defer wg.Done()

	for {
		if t, ok := <-taskChannel; ok && atomic.LoadInt32(errorsCountDown) >= 0 {
			if t() != nil {
				atomic.AddInt32(errorsCountDown, -1)
			}
		} else {
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Если "нельзя" создавать воркеров или нет задач - то ничего и не будем делать
	if n <= 0 || len(tasks) == 0 {
		return nil
	}

	// Если передано отрицательно число ошибок (равно как и 0), то считаем, что ошибки не допустимы
	// Но дадим таскам шанс, выставив счетчик в ноль
	if m < 0 {
		m = 0
	}

	// Синхронизация количества ошибок
	errorsCounter := int32(m)

	// Создадим канал для тасков
	tasksCh := make(chan Task, len(tasks))

	// WG для ожидания завершения горутин
	wg := sync.WaitGroup{}

	// Запустим N горутин
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(tasksCh, &wg, &errorsCounter)
	}

	// Отправляем задачи в канал на исполнение
	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	// Дожидаемся завершения горутин
	wg.Wait()

	if errorsCounter < 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
