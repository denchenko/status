package status

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
)

//go:embed page.tmpl
var tpl []byte

type Page struct {
	tmpl *template.Template
	hc   *HealthChecker
	urls []navURL
}

func NewPage() (*Page, error) {
	tmpl, err := template.New("page").Parse(string(tpl))
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
	Name string
	URL  string
}

func (p *Page) WithURL(name, url string) {
	p.urls = append(p.urls, navURL{
		Name: name,
		URL:  url,
	})
}

func (p *Page) Handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version, revision, commitDate := retrieveBuildInfo()

		var healthResults []HealthCheckResult
		if p.hc != nil {
			var err error
			healthResults, err = p.hc.Check(r.Context())
			if err != nil {
				http.Error(w, fmt.Sprintf("Error checking health: %v", err), http.StatusInternalServerError)
				return
			}
		}

		data := struct {
			Title         string
			Heading       string
			Message       string
			Version       string
			Revision      string
			CommitDate    string
			HealthResults []HealthCheckResult
			URLs          []navURL
		}{
			Title:         "System Status",
			Heading:       "System Status Dashboard",
			Message:       "Current status of all system components and dependencies.",
			Version:       version,
			Revision:      revision,
			CommitDate:    commitDate,
			HealthResults: healthResults,
			URLs:          p.urls,
		}

		w.Header().Set("Content-Type", "text/html")
		if err := p.tmpl.Execute(w, data); err != nil {
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
