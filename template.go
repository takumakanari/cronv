package cronv

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

const TEMPLATE = `
<html>
<head>
<title>{{.Opts.Title}} | {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</title>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
</head>
<body>
  <div class="container-fluid">
    <h1>
			{{.Opts.Title}}&nbsp;<small class="text-muted">From {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</small>
		</h1>

    <br>

    {{if .Extras}}
      <h3>Extra</h3>
      <div id="cronv-extra" style="width:{{.Opts.Width}}%;">
        <dl class="row">
          {{range $index, $extra := .Extras}}
            <dt class="col-sm-1">{{$extra.Label}}</dt>
            <dd class="col-sm-11">{{$extra.Job}}</dd>
          {{end}}
        </dl>
      </div>
      <hr>
    {{end}}

    <h3>Timeline<h3>
    <div id="cronv-timeline" style="height:100%; width:{{.Opts.Width}}%;">
      <b>Loading...</b>
    </div>
  </div>
  <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
  <script type="text/javascript">
     google.charts.load("current", {packages:["timeline"]});
     google.charts.setOnLoadCallback(function() {
       var container = document.getElementById('cronv-timeline');
       var chart = new google.visualization.Timeline(container);
       var dataTable = new google.visualization.DataTable();
        dataTable.addColumn({ type: 'string', id: 'job' });
        dataTable.addColumn({ type: 'string', id: 'dummy bar label' });
        dataTable.addColumn({ type: 'string', role: 'tooltip' });
        dataTable.addColumn({ type: 'date', id: 'Start' });
        dataTable.addColumn({ type: 'date', id: 'End' });

				var tasks = {};
				{{ $timeFrom := .TimeFrom }}
				{{ $timeTo := .TimeTo }}
				{{range $index, $cronv := .CronEntries}}
					{{ $job := JSEscapeString $cronv.Crontab.Job }}
					tasks['{{$job}}'] = [];
				{{end}}
				{{range $index, $cronv := .CronEntries}}
					{{ $job := JSEscapeString $cronv.Crontab.Job }}
					{{if IsRunningEveryMinutes $cronv.Crontab }}
						tasks['{{$job}}'].push(['{{$job}}', '', 'Every minutes {{$job}}', {{NewJsDate $timeFrom}}, {{NewJsDate $timeTo}}]);
					{{else}}
						{{range CronvIter $cronv}}tasks['{{$job}}'].push(['{{$job}}', '', '{{DateFormat .Start "15:04"}} {{$job}}', {{NewJsDate .Start}}, {{NewJsDate .End}}]);{{end}}
					{{ end }}
				{{end}}

				var taskByJobCount = [];
				for (var k in tasks) taskByJobCount.push({name: k, size: tasks[k].length});
				taskByJobCount.sort(function(a, b) {
					if (a.size == b.size) return 0;
					return a.size > b.size ? -1 : 1;
				});

				var rows = [];
				for (var i = 0; i < taskByJobCount.length; i++) {
					jobs = tasks[taskByJobCount[i].name];
					var jl = jobs.length;
					for (var j = 0; j < jl; j++) rows.push(jobs[j]);
				}

        if (rows.length > 0) {
          dataTable.addRows(rows);
          chart.draw(dataTable, {
            timeline: {
              colorByRowLabel: true
            },
            avoidOverlappingGridLines: false
          });
        } else {
          container.innerHTML = '<div class="alert alert-success"><strong>Woops!</strong> There is no data!</div>';
        }
     });
  </script>
</body>
</html>
`

func makeTemplate() *template.Template {
	funcMap := template.FuncMap{
		"CronvIter": func(cronv *Cronv) <-chan *Exec {
			return cronv.iter()
		},
		"JSEscapeString": func(v string) string {
			return template.JSEscapeString(strings.TrimSpace(v))
		},
		"NewJsDate": func(v time.Time) string {
			return fmt.Sprintf("new Date(%d,%d,%d,%d,%d)", v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute())
		},
		"DateFormat": func(v time.Time, format string) string {
			return v.Format(format)
		},
		"IsRunningEveryMinutes": func(c *Crontab) bool {
			return c.isRunningEveryMinutes()
		},
	}
	return template.Must(template.New("").Funcs(funcMap).Parse(TEMPLATE))
}
