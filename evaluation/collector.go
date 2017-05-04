package evaluation

import (
	"github.com/creichlin/pentaconta/logger"
	"strings"
	"sync"
	"time"
)

type Status struct {
	Samples  int
	Services map[string]*ServiceStats
}

type ServiceStats struct {
	Errors       int
	Crashes      int
	Terminations int
	Logs         int
}

type Collector struct {
	logs         chan logger.Log
	seconds      []map[string]*ServiceStats
	secondsMutex sync.Mutex
}

func EvaluationCollector() *Collector {
	seconds := []map[string]*ServiceStats{}
	for i := 0; i < 120; i++ {
		seconds = append(seconds, map[string]*ServiceStats{})
	}

	collector := &Collector{
		logs:         make(chan logger.Log, 100),
		seconds:      seconds,
		secondsMutex: sync.Mutex{},
	}

	go collector.start()
	go collector.clear()
	return collector
}

func (c *Collector) Status(seconds int) *Status {
	results := map[string]*ServiceStats{}
	samples := 0
	now := time.Now().Unix()

	c.secondsMutex.Lock()
	defer c.secondsMutex.Unlock()

	for i := now - int64(seconds); i < now; i++ {
		services := c.seconds[modI(i, len(c.seconds))]
		if len(services) > 0 {
			samples++
		}
		for key, stats := range services {
			target := results[key]
			if target == nil {
				target = &ServiceStats{}
				results[key] = target
			}
			target.Crashes += stats.Crashes
			target.Errors += stats.Errors
			target.Logs += stats.Logs
			target.Terminations += stats.Terminations
		}
	}
	return &Status{
		Samples:  samples,
		Services: results,
	}
}

// clears the time slot in the future which might contain old data
// does it 3 times every second which is dirty but works well if not
// under heavy load
func (c *Collector) clear() {
	for {
		c.secondsMutex.Lock()
		nextSlot := modI(time.Now().Unix()+1, len(c.seconds))
		c.seconds[nextSlot] = map[string]*ServiceStats{}
		c.secondsMutex.Unlock()
		time.Sleep(time.Millisecond * 300)
	}
}

func (c *Collector) start() {
	for lg := range c.logs {
		c.secondsMutex.Lock()
		slot := modI(lg.Time.Unix(), len(c.seconds))
		stats := c.seconds[slot][lg.Service]
		if stats == nil {
			stats = &ServiceStats{}
			c.seconds[slot][lg.Service] = stats
		}
		switch lg.Level {
		case logger.PENTACONTA:
			if strings.HasPrefix(lg.Message, "Terminated service with exit status") {
				stats.Crashes++
			} else if lg.Message == "Started service" {
				stats.Terminations++
			}
		case logger.STDERR:
			stats.Errors++
		case logger.STDOUT:
			stats.Logs++
		}
		c.secondsMutex.Unlock()
	}
}

func (c *Collector) Log(lg logger.Log) {
	c.logs <- lg
}

func modI(val int64, mod int) int {
	res := int(val % int64(mod))
	if res < 0 {
		res += mod
	}
	return res
}
