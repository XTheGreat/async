# Async goroutine sychronization & concurrency tool

How to use:
1. Job: safely run goroutines, supports graceful shutdown. (Credits: @kai.pei)
```
Example:

job := MakeJob()

// register in graceful/background listener
gracefulServer.Go(job) or background.Observe(job)

// invoke in any function
job.Go(func(){

	err := SomeFunc(arg1, arg2, ...)
	if err == nil {
		// error handling
	}
})

```

2. Threading: handles safe concurrent read & execution of data in multiple goroutine threads, and supports graceful shutdown.
```
Example:

th, err := MakeThreading(Util{
	Read:         func(offset, block int) (interface{}, error) { return []int{}, nil },
	Execute:      func(data interface{}, block int) error { return nil },
	MakeInterval: func() (idx, limit int, err error) { return 0, 100, nil },
	BlockSize:    10,
	Nodes:        5,
})
if err != nil {
	// error handling
}

// register in graceful/background listener
gracefulServer.Go(th) or background.Observe(th)


// run tasks
th.Go()
```

3. Background: handles graceful shutdown of background jobs when server is not running. If http server is not running, graceful cannot be used to shutdown running jobs. Instead, use the background listener.
```
Example:

// main function
background = NewListener()
defer background.Quit()


// function call
background.Observe(Job1)
background.Observe(Job2)
background.Observe(Cron1)
background.Observe(Cron2)


// notes: if signal termination is invoked at any time, the running jobs/cron will terminate gracefully.
```