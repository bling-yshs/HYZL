package main

import (
	"time"
)

func scheduleList() {
	go schedule(3*time.Hour, giteeAPI.updateReleaseDataFromAPI, true)
	go schedule(3*time.Hour, getAnnouncement, true)
	go schedule(3*time.Hour, update, true)
}

func schedule(delay time.Duration, fn func(), runRightNow bool) {
	if runRightNow {
		fn()
	}
	for {
		timer := time.NewTimer(delay)
		<-timer.C
		fn()
	}
}
