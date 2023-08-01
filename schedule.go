package main

import (
	"time"
)

func scheduleList() {
	go schedule(1*time.Hour-30*time.Second, giteeAPI.updateReleaseDataFromAPI, true)
	go schedule(1*time.Hour, wrapNoErrorFunc(getAnnouncement), true)
	go schedule(1*time.Hour, wrapNoErrorFunc(update), true)
}

func schedule(delay time.Duration, fn func() error, runRightNow bool) {
	if runRightNow {
		_ = fn()
	}
	for {
		timer := time.NewTimer(delay)
		<-timer.C
		_ = fn()
	}
}

func wrapNoErrorFunc(fn func()) func() error {
	return func() error {
		fn()
		return nil
	}
}
