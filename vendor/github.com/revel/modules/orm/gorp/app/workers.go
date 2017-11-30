package gorp

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// The worker container
type DbWorkerContainer struct {
	SharedWorker
	mutex            sync.Mutex
	Workers          []*DbWorker
	NumWorkers       int
	LongWorkTimeout  int64
	StartWorkTimeout int64
	Db               *DbGorp
}

// The timeoutInfo for monitoring long running processes
type timeoutInfo struct {
	worker  *DbWorker
	started time.Time
	ended   time.Time
	state   WorkerPhase
}
type DbWorker struct {
	Id int
	Db *DbGorp
	SharedWorker
	WorkUnit       int
	SharedData     map[string]interface{}
	TimeInfo       *timeoutInfo
	TimeoutChannel chan *timeoutInfo
}
type SharedWorker struct {
	workInfo       DbWorkInfo
	InputChannel   chan interface{}
	OutputChannel  chan interface{}
	ControlChannel chan func() (WorkerPhase, *DbWorker)
}
type DbWorkInfo interface {
	Status(phase WorkerPhase, worker *DbWorker)
	Work(value interface{}, worker *DbWorker)
}
type DbCallbackImplied struct {
	StatusFn func(phase WorkerPhase, worker *DbWorker)
	WorkFn   func(value interface{}, worker *DbWorker)
}
type WorkerPhase int

const (
	Start WorkerPhase = iota
	Stop
	StartJob
	EndJob
	JobLongrunning
)
// Creates a container to run the group of workers (up to a max of maxNumWorkers), does not return to all workers are completed)
// If returnResults is true then the task MUST write to the DbWorker.OutputChannel once for every task
func WorkParallel(db *DbGorp, tasks []func(worker *DbWorker), returnResults bool, maxNumWorkers int, timeouts int) (results []interface{}, err error) {
	if maxNumWorkers == 0 {
		maxNumWorkers = len(tasks)
	}

	// Create a container with no status callback
	container := NewDbWorker(db,
		MakeCallback(nil,
			func(value interface{}, worker *DbWorker) {
				task := value.(func(worker *DbWorker))
				task(worker)
			}), maxNumWorkers)
	err = container.Start()
	if err != nil {
		return
	}
	for _, task := range tasks {
		container.InputChannel <- task
	}
	if returnResults {
		for range tasks {
			result := <-container.OutputChannel
			results = append(results, result)
		}
	}

	container.Close(timeouts)
	return
}

// This creates a DbWorkerContainer with the number of working threads already started.
// Each working thread has their own database instance running.
func NewDbWorker(db *DbGorp, workInfo DbWorkInfo, numWorkers int) (container *DbWorkerContainer) {

	container = &DbWorkerContainer{
		SharedWorker: SharedWorker{
			InputChannel:   make(chan interface{}, numWorkers),
			OutputChannel:  make(chan interface{}, numWorkers),
			ControlChannel: make(chan func() (WorkerPhase, *DbWorker), numWorkers),
			workInfo:       workInfo,
		},
		NumWorkers:       numWorkers,
		Db:               db,
		StartWorkTimeout: 0,
		LongWorkTimeout:  0,
	}
	return
}
func (container *DbWorkerContainer) Start() (err error) {
	for x := 0; x < container.NumWorkers; x++ {
		go startWorker(container, container.Db, x)
	}
	// Make sure all containers are running before returning
	for x := 0; x < container.NumWorkers; x++ {
		if container.StartWorkTimeout > 0 {
			select {
			case result := <-container.ControlChannel:
				state, source := result()
				if state != Start {
					container.Close(5)
					err = fmt.Errorf("Failed to start workers %d", source)
					return
				}
			case <-time.After(time.Second * time.Duration(container.StartWorkTimeout)):
				container.Close(5)
				err = fmt.Errorf("Failed to start worker timeout")
				return
			}
		} else {
			result := <-container.ControlChannel
			state, source := result()
			if state != Start {
				container.Close(5)
				err = fmt.Errorf("Failed to start workers %d", source)
				return
			}
		}
	}
	return
}

func (container *DbWorkerContainer) Close(timeouts int) (totalWork int, err error) {
	close(container.InputChannel)
	for x := 0; x < len(container.Workers); x++ {
		// Allow close to continue even if a worker does not respond
		if timeouts > 0 {
			select {
			case result := <-container.ControlChannel:
				_, worker := result()
				totalWork += worker.WorkUnit
			case <-time.After(time.Second * time.Duration(timeouts)):
			}
		} else {
			_, worker := (<-container.ControlChannel)()
			totalWork += worker.WorkUnit
		}
	}
	close(container.OutputChannel)
	return
}

// Called by using "go" to invoke, creates a DBWorker, and starts a watchdog channel
func startWorker(container *DbWorkerContainer, db *DbGorp, id int) {
	newDb, _ := db.CloneDb(true)
	worker := &DbWorker{
		Db:         newDb,
		Id:         id,
		SharedData: map[string]interface{}{},
		SharedWorker: SharedWorker{
			workInfo:       container.workInfo,
			InputChannel:   container.InputChannel,
			OutputChannel:  container.OutputChannel,
			ControlChannel: container.ControlChannel,
		},
	}
	// Close the database after worker has ended (Start returned
	defer worker.Db.Close()
	container.mutex.Lock()
	container.Workers = append(container.Workers, worker)
	container.mutex.Unlock()
	// Only monitor jobs if Status function defined and a timeout is also defined
	if worker.workInfo.Status != nil && container.LongWorkTimeout > 0 {
		worker.TimeoutChannel = make(chan *timeoutInfo)
		go worker.TimeInfo.start(worker.TimeoutChannel, container.LongWorkTimeout)
	}
	worker.start()
}

// Starts the worker, continues running until inputchannel is closed
func (worker *DbWorker) start() {
	if worker.workInfo.Status != nil {
		worker.workInfo.Status(Start, worker)
	}
	worker.ControlChannel <- func() (WorkerPhase, *DbWorker) { return Start, worker }
	for job := range worker.InputChannel {
		worker.invoke(job)
	}
	if worker.workInfo.Status != nil {
		worker.workInfo.Status(Stop, worker)
	}
	worker.ControlChannel <- func() (WorkerPhase, *DbWorker) { return Stop, worker }
	if worker.TimeoutChannel != nil {
		close(worker.TimeoutChannel)
	}
}

// Wrapper to prevent panics from disturbing the channel
func (worker *DbWorker) invoke(job interface{}) {
	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1024)
			count := runtime.Stack(trace, true)
			moduleLogger.Error("Recover from panic: ", "error", err)
			moduleLogger.Error("Stack", "size", count, "trace", string(trace))
		}
	}()
	// Setup the timeout information
	if worker.TimeoutChannel != nil {
		worker.TimeInfo = &timeoutInfo{worker: worker, started: time.Now(), state: StartJob}
		worker.TimeoutChannel <- worker.TimeInfo
	}
	worker.workInfo.Work(job, worker)
	if worker.TimeoutChannel != nil {
		worker.TimeInfo.state = EndJob
		worker.TimeoutChannel <- worker.TimeInfo
	}
}

// A function to return an object that is a valid DbCallback
func MakeCallback(status func(phase WorkerPhase, worker *DbWorker), work func(value interface{}, worker *DbWorker)) DbWorkInfo {
	return &DbCallbackImplied{StatusFn: status, WorkFn: work}
}

// Call the status function if available
func (dbCallback *DbCallbackImplied) Status(phase WorkerPhase, worker *DbWorker) {
	if dbCallback.StatusFn != nil {
		dbCallback.StatusFn(phase, worker)
	}
}

// Calls the work function
func (dbCallback *DbCallbackImplied) Work(value interface{}, worker *DbWorker) {
	dbCallback.WorkFn(value, worker)
}

//Starts the timeout worker
func (_ *timeoutInfo) start(TimeoutChannel chan *timeoutInfo, timeout int64) {
	for j := range TimeoutChannel {
		j.started = time.Now()
		j.state = StartJob
		j.worker.workInfo.Status(j.state, j.worker)
		for {
			select {
			case complete, ok := <-TimeoutChannel:
				if !ok {
					// Channel closed returning...
					return
				}
				// Received new State, record and loop
				complete.worker.workInfo.Status(complete.state, complete.worker)
				break
			case <-time.After(time.Second * time.Duration(timeout)):
				j.worker.workInfo.Status(JobLongrunning, j.worker)
			}
		}
	}
}
