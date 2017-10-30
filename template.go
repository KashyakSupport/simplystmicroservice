package hellodatastore

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	tpl      *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() {
		t.tpl = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.tpl.Execute(w, nil)
}
