# Timecalc

Timecalc is a calculator for time values, i.e dates and durations.
It can be used to convert values between RFC3339 dates and timestamps, add dates and durations, filter and truncate dates.
Such functionality is usually available either using the python repl or a bunch of programs like date(1). Timecalc
is an attempt to simplify and gather such calculations is a simple program.

## Usage

```
$ ./timecalc
2026-10-28T14:42:39+02:00

$ ./timecalc print now
2026-10-28T14:42:39+02:00

$ ./timecalc -l utc print now
2026-10-28T12:42:45Z

$ ./timecalc -l utc add now 3h
2026-10-28T15:42:56Z

$ ./timecalc -l utc -o unix add now 3h
1793202188

$ ./timecalc -o unix 2026-10-28T00:00:00Z
1793145600

$ ./timecalc -o unixmilli 2026-10-28T00:00:00Z
1793145600000

$ ./timecalc -o unixmilli 2026-10-28T00:03:12Z
1793145792000

$ ./timecalc 1793145792000
2026-10-28T02:03:12+02:00

$ ./timecalc after now tomorrow
2026-10-29T00:00:00+02:00

$ ./timecalc -l utc after now tomorrow
2026-10-29T00:00:00Z

$ ./timecalc -l utc before now yesterday
2026-10-27T00:00:00Z

$ ./timecalc -l utc trunc now 5m
2026-10-28T12:40:00Z

$ ./timecalc diff now 15:00:00 16:00:00
14m50.907596542s
1h14m50.907596542s

$ ./timecalc gaps now 15:00:00 16:00:00
14m46.570367047s
1h0m0s
```

## Install

Timecalc is written in go. You need the go SDK which you can install with your OS package manager or download it from [here](https://go.dev/dl/).

To install timecalc, simply

```
go install github.com/anastasop/timecalc@latest
```

## License

Licensed under the 3-Clause BSD License.

## Bugs

Still unsure if i should keep the lisp style commands or use infix expressions.
