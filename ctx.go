package cronv

import (
	"fmt"
	"github.com/tkmgo/cronexpr"
	"os"
	"strings"
	"time"
)

type Cronv struct {
	Crontab         *Crontab
	expr            *cronexpr.Expression
	startTime       time.Time
	durationMinutes float64
}

func NewCronv(line string, startTime time.Time, durationMinutes float64) (*Cronv, error) {
	crontab, err := parseCrontab(line)
	if err != nil {
		return nil, err
	}

	expr, err := cronexpr.Parse(crontab.Schedule.toCrontab())
	if err != nil {
		return nil, err
	}

	cronv := &Cronv{
		Crontab:         crontab,
		expr:            expr,
		startTime:       startTime,
		durationMinutes: durationMinutes,
	}
	return cronv, nil
}

type EnvDef struct {
	Name string
	Val  string
}

func maybeParseEnv(line string) (*EnvDef, bool) {
	if kv := strings.SplitN(line, "=", 2); len(kv) == 2 { // K=V => [K V]
		return &EnvDef{
			Name: kv[0],
			Val:  kv[1],
		}, true
	}
	return nil, false
}

type Exec struct {
	Start time.Time
	End   time.Time
}

func (self *Cronv) iter() <-chan *Exec {
	ch := make(chan *Exec)
	eneTime := self.startTime.Add(time.Duration(self.durationMinutes) * time.Minute)
	next := self.expr.Next(self.startTime)
	go func() {
		for next.Equal(eneTime) || eneTime.After(next) {
			ch <- &Exec{
				Start: next,
				End:   next.Add(time.Duration(1) * time.Minute),
			}
			next = self.expr.Next(next)
		}
		close(ch)
	}()
	return ch
}

type CronvCtx struct {
	Opts            *Command
	TimeFrom        time.Time
	TimeTo          time.Time
	CronEntries     []*Cronv
	EnvEntries      []*EnvDef
	durationMinutes float64
}

func NewCtx(opts *Command) (*CronvCtx, error) {
	timeFrom, err := opts.toFromTime()
	if err != nil {
		return nil, err
	}

	durationMinutes, err := opts.toDurationMinutes()
	if err != nil {
		return nil, err
	}

	return &CronvCtx{
		Opts:            opts,
		TimeFrom:        timeFrom,
		TimeTo:          timeFrom.Add(time.Duration(durationMinutes) * time.Minute),
		CronEntries:     []*Cronv{},
		durationMinutes: durationMinutes,
	}, nil
}

func (self *CronvCtx) AppendNewLine(line string) (bool, error) {
	trimed := strings.TrimSpace(line)
	if len(trimed) == 0 || string(trimed[0]) == "#" {
		return false, nil
	}

	cronv, err := NewCronv(trimed, self.TimeFrom, self.durationMinutes)
	if err != nil {
		switch err.(type) {
		case *InvalidTaskError:
			// Try to parse line as env-def (if line is not a cron entry)
			if envDef, ok := maybeParseEnv(line); ok {
				self.EnvEntries = append(self.EnvEntries, envDef)
			}

			// Pass (not a cron entry)
			return false, nil
		default:
			return false, fmt.Errorf("Failed to analyze cron '%s': %s", line, err)
		}
	}

	self.CronEntries = append(self.CronEntries, cronv)
	return true, nil
}

func (self *CronvCtx) Dump() (string, error) {
	output, err := os.Create(self.Opts.OutputFilePath)
	if err != nil {
		return "", err
	}
	makeTemplate().Execute(output, self)
	return self.Opts.OutputFilePath, nil
}
