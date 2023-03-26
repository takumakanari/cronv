# cronv

Visualize your cron schedules in crontab

![cronv output 1d](https://raw.github.com/wiki/takumakanari/cronv/images/outputs/cronv-1d.png)
![cronv output 30m](https://raw.github.com/wiki/takumakanari/cronv/images/outputs/cronv-30m.png)


## Install

## go install

```shell
go install -v github.com/takumakanari/cronv/cronv@0.4.5
```

```shell
$ cronv --help
```

## Basic usage

Cronv can parse your crontab from stdin like as follows:

```shell
$ crontab -l | cronv -o ./my_cron_schedule.html
```

You can also specify the duration to analysis job schedules.

In a case like the follows, the job schedules will be analyzed from now to 24 hours later:


```shell
$ crontab -l | cronv -o ./my_cron_schedule.html -d 24h
```



> Cronv can parse cron entry written in basic cron format.
You can see the basically crontab specofication in [https://en.wikipedia.org/wiki/Cron#CRON_expression](https://en.wikipedia.org/wiki/Cron#CRON_expression).


## Options

```shell
Application Options:
  -o, --output=    path to .html file to output (default: ./crontab.html)
  -d, --duration=  duration to visualize in N{suffix} style. e.g.)
                   1d(day)/1h(hour)/1m(minute) (default: 6h)
      --from-date= start date in the format '2006/01/02' to visualize (default:
                   2017/03/15)
      --from-time= start time in the format '15:04' to visualize (default:
                   19:28)
  -t, --title=     title/label of output (default: cron tasks)
  -w, --width=     Table width of output (default: 100)

Help Options:
  -h, --help       Show this help message
```

## Examples
Analyze crontab for 6 hours (by default) from now, , output html file to default path:
```shell
$ crontab -l | cronv
```

For 1 day from now, output html file to default path:

```shell
$ crontab -l | cronv -d 1d
```

For 12 hours from 21:00, today:

```shell
$ crontab -l | cronv --from-time 21:00 -d 12h
```

For 30 minuts from now, output html file to path/to/output.html:

```shell
$ crontab -l | cronv -d 30m -o path/to/output.html
```

For 2 hours from 2016/12/24 17:30, output html file to path/to/output2.html:

```shell
$ crontab -l | cronv --from-date '2016/12/24' --from-time 17:30 -d 2h -o path/to/output2.html
```

With original title/label:

```shell
$ crontab -l | cronv -d 1d -t "crontab@`hostname`"  # title/label of html file will be 'crontab@myhost'
```

With *width* to spread output table:

```shell
$ crontab -l | cronv -o path/to/output2.html -w 180 # table width be 180% of the screen width (100% by default)

$ crontab -l | cronv -o path/to/output2.html -w 75 # be 75% of the screen width
```


## Development

Using **dep**.

```shell
$ cd /path/to/cronv
$ dep ensure
$ crontab -l | go run cronv/main.go
```


## TODO

- Add output format/style other than HTML.
- Filter entries in output HTML file.


## Patch

Welcome!
