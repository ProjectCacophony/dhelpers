package dhelpers

import (
	"context"
	"time"

	"github.com/bsm/redis-lock"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// Job defines setting for a job
type Job struct {
	// Name should be unique, prefixed by module
	Name string
	// Cron is a cron expression https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format, https://crontab.guru/
	Cron string
	// AtLaunch if set to true will start the Job at launch
	AtLaunch bool
	Job      func()
}

func jobLockKey(jobName string) (key string) {
	return "project-d:job:" + jobName + ":status"
}

// JobStart returns true and a locker if the Job has been started successfully, returns false if the Job is already running
// after timeout the locks unlock itself, use locker.Lock() to renew a lock
func JobStart(jobName string, timeout time.Duration) (start bool, locker *lock.Locker, err error) {
	locker = lock.New(cache.GetRedisClient(), jobLockKey(jobName), &lock.Options{
		LockTimeout: timeout,
		RetryCount:  0,
		RetryDelay:  100 * time.Millisecond,
	})

	// lock locker
	start, err = locker.LockWithContext(context.Background())
	return start, locker, err
}

// JobErrorHandler handles errors at jobs, defer to this: defer JobErrorHandler(jobName)
func JobErrorHandler(jobName string) {
	err := recover()
	if err != nil {
		// handle errors
		HandleJobErrorWith("Worker", jobName, err.(error), SentryErrorHandler)
	}
}
