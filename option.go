package cronv

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	optDateFormat        = "2006/01/02"
	optTimeFormat        = "15:04"
	optDefaultDuration   = "6h"
	optDefaultOutputPath = "./crontab.html"
	optDefaultTitle      = "Cron Tasks"
	optDefaultWidth      = 100
)

type Option struct {
	OutputFilePath string `short:"o" long:"output" description:"path to .html file to output"`
	Duration       string `short:"d" long:"duration" description:"duration to visualize in N{suffix} style. e.g.) 1d(day)/1h(hour)/1m(minute)"`
	FromDate       string `long:"from-date" description:"start date in the format '2006/01/02' to visualize"`
	FromTime       string `long:"from-time" description:"start time in the format '15:04' to visualize"`
	Title          string `short:"t" long:"title" description:"title/label of output"`
	Width          int    `short:"w" long:"width" description:"Table width of output"`
}

func (o *Option) toFromTime() (time.Time, error) {
	return time.Parse(fmt.Sprintf("%s %s", optDateFormat, optTimeFormat),
		fmt.Sprintf("%s %s", o.FromDate, o.FromTime))
}

func (o *Option) toDurationMinutes() (float64, error) {
	length := len(o.Duration)
	if length < 2 {
		return 0, errors.Errorf("invalid duration format: '%s'", o.Duration)
	}

	duration, err := strconv.Atoi(string(o.Duration[:length-1]))
	if err != nil {
		return 0, errors.Errorf("invalid duration format: '%s', %s", o.Duration, err)
	}

	unit := string(o.Duration[length-1])
	switch strings.ToLower(unit) {
	case "d":
		return float64(duration * 24 * 60), nil
	case "h":
		return float64(duration * 60), nil
	case "m":
		return float64(duration), nil
	}

	return 0, errors.Errorf("invalid duration format: '%s', '%s' is not in d/h/m", o.Duration, unit)
}

func NewCronvOption(t time.Time) *Option {
	now := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	return &Option{
		OutputFilePath: optDefaultOutputPath,
		Duration:       optDefaultDuration,
		FromDate:       now.Format(optDateFormat),
		FromTime:       now.Format(optTimeFormat),
		Title:          optDefaultTitle,
		Width:          optDefaultWidth,
	}
}
