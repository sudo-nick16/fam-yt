package workerpool

type Task interface {
	Execute() error
	Failed(error)
	Success()
}

type Worker struct {
	taskQueue *chan Task
}

func NewWorker(taskQueue *chan Task) Worker {
	return Worker{
		taskQueue,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			task := <-*w.taskQueue
			err := task.Execute()
			if err != nil {
				task.Failed(err)
				continue
			}
			task.Success()
		}
	}()
}

type WorkerPool struct {
	taskQueue chan Task
	workers   []Worker
	quit      bool
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	taskQueue := make(chan Task)
	workers := make([]Worker, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workers[i] = NewWorker(&taskQueue)
	}
	return &WorkerPool{
		taskQueue: taskQueue,
		workers:   workers,
		quit:      false,
	}
}

func (wp *WorkerPool) AddTask(t Task) {
	wp.taskQueue <- t
}

func (wp *WorkerPool) Start() {
	for i := 0; i < len(wp.workers); i++ {
		wp.workers[i].Start()
	}
}

func (wp *WorkerPool) Stop() {
	// TODO: Implement a way to stop the workers
}
