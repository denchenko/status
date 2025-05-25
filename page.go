package status

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
)

var (
	//go:embed page.tmpl
	defaultTemplateContent string
	defaultTemplate        = template.Must(template.New("page").Parse(defaultTemplateContent))
)

// Page represents a status page that can be served via HTTP
type Page struct {
	title string
	tmpl  *template.Template
	hc    *HealthChecker
	links []Link
}

// PageOption is a function that configures a Page
type PageOption func(*Page)

// WithTitle sets the title of the status page
func WithTitle(title string) PageOption {
	return func(p *Page) {
		p.title = title
	}
}

// WithTemplate sets a custom HTML template for the status page
func WithTemplate(tmpl *template.Template) PageOption {
	return func(p *Page) {
		p.tmpl = tmpl
	}
}

// WithHealthChecker sets the health checker for the status page
func WithHealthChecker(hc *HealthChecker) PageOption {
	return func(p *Page) {
		p.hc = hc
	}
}

// WithLink adds a navigation link to the status page
func WithLink(name, url string) PageOption {
	return func(p *Page) {
		p.links = append(p.links, Link{
			Name: name,
			URL:  url,
		})
	}
}

// NewPage creates a new status page with the given options
func NewPage(opts ...PageOption) *Page {
	p := &Page{
		title: "System Status",
		tmpl:  defaultTemplate,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Link represents a navigation link in the status page
type Link struct {
	Name string
	URL  string
}

// PageData contains the data that will be rendered in the status page template
type PageData struct {
	Title         string
	Version       string
	Revision      string
	CommitDate    string
	HealthResults []HealthCheckResult
	Links         []Link
}

// Handler returns an HTTP handler that serves the status page
func (p *Page) Handler() http.HandlerFunc {
	version, revision, commitDate := retrieveBuildInfo()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var healthResults []HealthCheckResult
		if p.hc != nil {
			var err error
			healthResults, err = p.hc.Check(r.Context())
			if err != nil {
				http.Error(w, fmt.Sprintf("Error checking health: %v", err), http.StatusInternalServerError)
				return
			}
		}

		data := PageData{
			Title:         p.title,
			Version:       version,
			Revision:      revision,
			CommitDate:    commitDate,
			HealthResults: healthResults,
			Links:         p.links,
		}

		w.Header().Set("Content-Type", "text/html")
		if err := p.tmpl.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Error executing template :%v", err), http.StatusInternalServerError)
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
