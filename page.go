package status

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"text/template"
)

const tpl = `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Heading}}</h1>
    <p>{{.Message}}</p>
</body>
</html>`

type Page struct {
	tmpl *template.Template
	hc   *HealthChecker
	urls []navURL
}

func NewPage() (*Page, error) {
	tmpl, err := template.New("page").Parse(tpl)
	if err != nil {
		return nil, fmt.Errorf("parsing html template: %w", err)
	}
	return &Page{
		tmpl: tmpl,
	}, nil
}

func (p *Page) WithHealthChecker(hc *HealthChecker) {
	p.hc = hc
}

type navURL struct {
	name string
	url  string
}

func (p *Page) WithURL(name, url string) {
	p.urls = append(p.urls, navURL{
		name: name,
		url:  url,
	})
}

func (p *Page) Handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		version, revision, commitDate := retrieveBuildInfo()
		data := struct {
			Title      string
			Heading    string
			Message    string
			Version    string
			Revision   string
			CommitDate string
		}{
			Title:      "Welcome Page",
			Heading:    "Hello, World!",
			Message:    "This is a simple Go HTML template example.",
			Version:    version,
			Revision:   revision,
			CommitDate: commitDate,
		}
		w.Header().Set("Content-Type", "text/html")
		err := p.tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})
}

const (
	vcsRevisionKey = "vcs.revision"
	vcsTimeKey     = "vcs.time"
)

func retrieveBuildInfo() (string, string, string) {
	var (
		version    = "unknown"
		revision   = "unknown"
		commitDate = "unknown"
	)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version, revision, commitDate
	}

	version = info.Main.Version

	for i := range info.Settings {
		switch info.Settings[i].Key {
		case vcsRevisionKey:
			revision = info.Settings[i].Value
		case vcsTimeKey:
			commitDate = info.Settings[i].Value
		}
	}

	return version, revision, commitDate
}
