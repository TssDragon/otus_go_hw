package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Если канала нет, то вернем закрытый канал (по сути - пустоту)
	if in == nil {
		ch := make(chan interface{})
		close(ch)
		return ch
	}

	take := func(done In, value In) Out {
		takeChannel := make(chan interface{})
		go func() {
			defer close(takeChannel)
			for {
				// Если получили сигнал на завершение - сразу выходим
				select {
				case <-done:
					return
				default:
				}

				select {
				case <-done:
					return
				case val, ok := <-value:
					// Проверка на закрытый канал
					if !ok {
						return
					}
					takeChannel <- val
				}
			}
		}()
		return takeChannel
	}

	// Если не будет ни одного стэйджа (пустой массив или nil) - вернем то, что и получили на вход
	outChannel := in

	for _, stage := range stages {
		if stage == nil {
			return outChannel
		}
		outChannel = stage(take(done, outChannel))
	}
	return outChannel
}
