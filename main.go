package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"
)

type TemporalValue struct {
	t time.Time
	d time.Duration
}

func (v TemporalValue) IsTime() bool {
	return !v.t.IsZero()
}

func (v TemporalValue) IsDuration() bool {
	return v.t.IsZero()
}

var (
	durationRE = regexp.MustCompile(`^[-+]?([0-9]*(\.[0-9]*)?[a-z]+)+$`)

	nowLocationArg = flag.String("l", "Local", "user location: UTC or Local or a TZ")
	outFormat      = flag.String("o", "", "output format: epoch,epochmilli,utc,local or a TZ")

	now time.Time
	nowLocation *time.Location
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: timecalc <cmd> temporal0 temporals...

Timecalc does calculations (add, sub etc) with temporals i.e dates and durations
using the go time package.

A temporal is a date or a duration and must be parseable by the go time package.
The special values are now, tomorrow, yesterday, today.

The available commands are:

add: temporals can be dates or durations and adds each one to temporal0.
diff: temporals must be dates substracts each one from temporal0.
round: temporals must be durations and rounds temporal0 accordingly.
trunc: temporals must be durations and truncates temporal0 accordingly.
after: temporals must be dates, returns those after temporal0.
before: temporals must be dates, returns those before temporal0.
gaps: temporals must be dates, prints the gaps between successive temporals.
print: converts and prints all temporals.

The semantics of the commands are in accordance with the go time package.

Flags:`)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	var err error
	if strings.ToUpper(*nowLocationArg) == "LOCAL" {
		nowLocation, err = time.LoadLocation("Local")
	} else if strings.ToUpper(*nowLocationArg) == "UTC" {
		nowLocation, err = time.LoadLocation("UTC")
	} else {
		nowLocation, err = time.LoadLocation(*nowLocationArg)
	} 
        if err != nil {
		log.Fatalf("unknown time location %s", *nowLocationArg)
        }
        now = time.Now().In(nowLocation)

	if flag.NArg() == 0 {
		output([]TemporalValue{TemporalValue{t: now}})
		os.Exit(0)
	}

	var cmdExec func([] TemporalValue) []TemporalValue
	cmd := flag.Arg(0)
	args := flag.Args()[1:]
	switch cmd {
	case "add": 
		cmdExec = add
	case "diff":
		cmdExec = diff 
	case "round":
		cmdExec = round
	case "trunc":
		cmdExec = trunc
	case "before":
		cmdExec = before
	case "after":
		cmdExec = after
	case "gaps":
		cmdExec = gaps 
	case "print", "conv":
		cmdExec = print 
	default:
		cmdExec = print
		args = flag.Args()
	}

	var temps []TemporalValue
	for _, s := range args {
		temps = append(temps, parseTemporalValue(s))
	}

	if len(temps) > 0 {
		output(cmdExec(temps))
	}
}

func add(temps []TemporalValue) []TemporalValue {
	var res []TemporalValue

	a := temps[0]
	for _, b := range temps[1:] {
		if a.IsTime() && b.IsDuration() {
			res = append(res, TemporalValue{t: a.t.Add(b.d)})
		} else if a.IsDuration() && b.IsTime() {
			res = append(res, TemporalValue{t: b.t.Add(a.d)})
		} else if a.IsDuration() && b.IsDuration() {
			res = append(res, TemporalValue{d: a.d + b.d})
		} else {
			log.Fatal("add is not applicable for dates")
		}
	}

	return res
}

func diff(temps []TemporalValue) []TemporalValue {
	var res []TemporalValue

	a := temps[0]
	for _, b := range temps[1:] {
		if a.IsTime() && b.IsTime() {
			res = append(res, TemporalValue{d: b.t.Sub(a.t)})
		} else {
			log.Fatal("diff applies only for dates")
		}
	}

	return res 
}

func round(temps []TemporalValue) []TemporalValue {
	var res []TemporalValue

	a := temps[0]
	for _, b := range temps[1:] {
		if a.IsTime() && b.IsDuration() {
			res = append(res, TemporalValue{t: a.t.Round(b.d)})
		} else if a.IsDuration() && b.IsTime() {
			res = append(res, TemporalValue{t: b.t.Round(a.d)})
		} else {
			log.Fatal("round applies only for durations")
		}
	}

	return res 
}

func trunc(temps []TemporalValue) []TemporalValue {
	var res []TemporalValue

	a := temps[0]
	for _, b := range temps[1:] {
		if a.IsTime() && b.IsDuration() {
			res = append(res, TemporalValue{t: a.t.Truncate(b.d)})
		} else if a.IsDuration() && b.IsTime() {
			res = append(res, TemporalValue{t: b.t.Truncate(a.d)})
		} else {
			log.Fatal("trunc applies only for durations")
		}
	}

	return res 
}

func before(temps []TemporalValue) []TemporalValue {
	var res []TemporalValue

	a := temps[0]
	for _, b := range temps[1:] {
		if a.IsTime() && b.IsTime() {
			if b.t.Before(a.t) {
				res = append(res, b)
			}
		} else {
			log.Fatal("before applies only to dates")
		}
	}

	return res 
}

func after(temps []TemporalValue) []TemporalValue {
	var res []TemporalValue

	a := temps[0]
	for _, b := range temps[1:] {
		if a.IsTime() && b.IsTime() {
			if b.t.After(a.t) {
				res = append(res, b)
			}
		} else {
			log.Fatal("after applies only to dates")
		}
	}

	return res 
}

func gaps(temps []TemporalValue) []TemporalValue {
	var res []TemporalValue

	a := temps[0]
	for _, b := range temps[1:] {
		if a.IsTime() && b.IsTime() {
			res = append(res, TemporalValue{d: b.t.Sub(a.t)})
			a = b
		} else {
			log.Fatal("gaps applies only for dates")
		}
	}

	return res
}

func print(temps []TemporalValue) []TemporalValue {
	return temps
}

func newTemporalValueFromTime(t time.Time) TemporalValue {
	return TemporalValue{t: t}
}

func newTemporalValueFromDuration(d time.Duration) TemporalValue {
	return TemporalValue{d: d}
}

func parseTemporalValue(s string) TemporalValue {
	// try parse keywords
	if s == "now" {
		return TemporalValue{t: now}
	}
	if s == "yesterday" {
		t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return TemporalValue{t: t.Add(time.Duration(-24) * time.Hour)}
	}
	if s == "today" {
		t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return TemporalValue{t: t}
	}
	if s == "tomorrow" {
		t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return TemporalValue{t: t.Add(time.Duration(24) * time.Hour)}
	}

	// try parse as duration
	if durationRE.MatchString(s) {
		d, err := time.ParseDuration(s)
		if err != nil {
			log.Fatalf("failed to parse duration %s: %v", s, err)
		}
		return TemporalValue{d: d}
	}

	// try parse as timestamp
	num := int64(0)
	for _, b := range []byte(s) {
		if '0' <= b && b <= '9' {
			num = num*10 + int64(b-'0')
		} else {
			num = 0
			break
		}
	}
	if num > 10000000000 {
		return newTemporalValueFromTime(time.UnixMilli(num))
	} else if num > 0 {
		return newTemporalValueFromTime(time.Unix(num, 0))
	}

	// try standard formats
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return newTemporalValueFromTime(t)
	}
	if t, err := time.Parse(time.DateOnly, s); err == nil {
		return newTemporalValueFromTime(
			time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location()))
	}
	if t, err := time.Parse(time.TimeOnly, s); err == nil {
		return newTemporalValueFromTime(
			time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), now.Location()))
	}

	log.Fatalf("argument %s can't be parsed as time", s)
	return TemporalValue{} // silence compiler
}

func output(s []TemporalValue) {
	var loc *time.Location
	outMilli := false
	outEpoch := false

	switch strings.ToLower(*outFormat) {
	case "":
		loc = nowLocation
	case "epoch", "unix", "u":
		outEpoch = true
	case "epochmilli", "unixmilli", "um":
		outMilli = true
	case "local":
		loc, _ = time.LoadLocation("Local")
	case "gmt", "utc":
		loc, _ = time.LoadLocation("UTC")
	default:
		l, err := time.LoadLocation(*outFormat)
		if err != nil {
                	log.Fatalf("unknown time location %s", *outFormat)
		}
		loc = l
	}	

	for _, t := range s {
		if t.IsDuration() {
			fmt.Println(t.d.String())
		} else if t.IsTime() {
			if outMilli {
				fmt.Println(t.t.UnixMilli())
			} else if outEpoch {
				fmt.Println(t.t.Unix())
			} else {
				fmt.Println(t.t.In(loc).Format(time.RFC3339))
			}
		}
	}
}
