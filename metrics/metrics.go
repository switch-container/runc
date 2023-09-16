package metrics

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type TimeMetric struct {
	start   map[string]time.Time
	elapsed map[string]time.Duration
}

var Timer *TimeMetric

func (timer *TimeMetric) StartTimer(name string) error {
	if _, ok := timer.start[name]; ok {
		return fmt.Errorf("duplicated timer %s", name)
	}
	timer.start[name] = time.Now()
	return nil
}

func (timer *TimeMetric) FinishTimer(name string) error {
	start, ok := timer.start[name]
	if !ok {
		return fmt.Errorf("%s timer does not start", name)
	}
	timer.elapsed[name] = time.Since(start)
	return nil
}

func (timer *TimeMetric) Report() {
	defer func() {
		// clear all timer once report
		timer.start = make(map[string]time.Time)
		timer.elapsed = make(map[string]time.Duration)
	}()
	for name := range timer.start {
		if _, ok := timer.elapsed[name]; !ok {
			logrus.Errorf("timer-%s-not-finish", name)
			return
		}
	}

  entry := "RUNC-METRIC: "

	for name, e := range timer.elapsed {
		entry += fmt.Sprintf("%s=%s, ", name, e.String())
	}

	logrus.Info(entry)
}

func (timer *TimeMetric) Clean() {
	timer.start = make(map[string]time.Time)
	timer.elapsed = make(map[string]time.Duration)
}

func init() {
	Timer = &TimeMetric{
		start:   make(map[string]time.Time),
		elapsed: make(map[string]time.Duration),
	}
}
