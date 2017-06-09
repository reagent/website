package ui

import (
	"bytes"
	"html/template"

	"github.com/milehighgophers/website/data"
)

const (
	indexTemplateStr = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Mile High Gopher Events</title>
		<link rel="stylesheet" href="/assets/styles.css">
	</head>
	<body>
	<img src="/assets/logo.png">
	{{range $key, $value := .UpcomingEvents}}
		<h1>{{$key}}</h1>
		<ul>
		{{range $value}}
			<li>{{.HumanTime}} -- {{.Name}}</li>
		{{else}}
			<div><strong>No Events</strong></div>
		{{end}}
		</ul>
	{{end}}
	</body>
</html>
`
)

var indexTemplate = template.Must(template.New("index").Parse(indexTemplateStr))

// Render will turn meetup event data into something to write out.
func Render(events map[string][]data.Event) []byte {
	index := &indexPage{
		events: events,
	}
	buf := &bytes.Buffer{}
	indexTemplate.Execute(buf, index)
	return buf.Bytes()
}

type indexPage struct {
	events map[string][]data.Event
}

func (p *indexPage) UpcomingEvents() map[string][]data.Event {
	threeEvents := make(map[string][]data.Event)
	for k, v := range p.events {
		threeEvents[k] = v[0:3]
	}
	return threeEvents
}
