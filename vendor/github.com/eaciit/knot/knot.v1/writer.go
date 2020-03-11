package knot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eaciit/toolkit"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func (r *WebContext) Write(data interface{}) error {
	if int(r.Config.OutputType) == 0 {
		r.Config.OutputType = DefaultOutputType
	}

	if r.Config.OutputType == OutputTemplate {
		return r.WriteError(r.WriteTemplate(data))
	}

	if r.Config.OutputType == OutputJson {
		return r.WriteError(r.WriteJson(data))
	}

	if r.Config.OutputType == OutputByte || r.Config.OutputType == OutputHtml {
		fmt.Fprint(r.Writer, data)
		return nil
	}

	return nil
}

func (r *WebContext) WriteCookie() error {
	for _, c := range r.Cookies() {
		http.SetCookie(r.Writer, c)
	}
	return nil
}

func (r *WebContext) WriteTemplate(data interface{}) error {
	var e error
	w := r.Writer
	cfg := r.Config

	if cfg.ViewName == "" {
		cfg.ViewName = strings.Join([]string{
			strings.ToLower(cfg.ControllerName),
			strings.ToLower(cfg.MethodName)}, "/") + ".html"
	}

	//viewFile := cfg.ViewsPath + cfg.ViewName

	//w.Header().Set("Content-Type", "text/html")
	if cfg.ViewName != "" {
		useLayout := false
		viewsPath := cfg.ViewsPath
		fixLogicalPath(&viewsPath, true, true)
		viewFile := viewsPath
		if cfg.LayoutTemplate != "" {
			useLayout = true
			//viewFile += cfg.LayoutTemplate
			viewFile = filepath.Join(viewFile, cfg.LayoutTemplate)
		} else {
			//viewFile += cfg.ViewName
			viewFile = filepath.Join(viewFile, cfg.ViewName)
		}
		if useLayout {
			buf := bytes.Buffer{}
			e = r.writeToTemplate(&buf, data, cfg.ViewName)
			if e != nil {
				return e
			}
			e = r.writeToTemplate(w, struct{ Content interface{} }{
				template.HTML(string(buf.Bytes()))}, cfg.LayoutTemplate)
			if e != nil {
				return e
			}
		} else {
			e = r.writeToTemplate(w, data, cfg.ViewName)
		}
		if e != nil {
			return e
		}
	} else {
		return fmt.Errorf("No template define for %s", strings.ToLower(r.Request.URL.String()))
	}
	return nil
}

func (r *WebContext) writeToTemplate(w io.Writer, data interface{}, templateFile string) error {
	cfg := r.Config
	viewsPath := cfg.ViewsPath
	viewFile := filepath.Join(viewsPath, templateFile)
	bs, e := ioutil.ReadFile(viewFile)
	if e != nil {
		return e
	}
	t, e := template.New("main").Funcs(template.FuncMap{
		"BaseUrl": func() string {
			base := "/"
			if cfg.AppName != "" {
				base += strings.ToLower(cfg.AppName)
			}
			if base != "/" {
				base += "/"
			}
			return base
		},
		"UnescapeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"NoCacheUrl": func(s string) string {
			concatenator := "?"
			if strings.Contains(s, "?") {
				concatenator = `&`
			}

			randomString := toolkit.RandomString(32)
			noCachedUrl := fmt.Sprintf("%s%snocache=%s", s, concatenator, randomString)
			return noCachedUrl
		},
	}).Parse(string(bs))
	if e != nil {
		return e
	}

	for _, includeFile := range cfg.IncludeFiles {
		if includeFile != cfg.LayoutTemplate && includeFile != templateFile {
			//includeFilePath := viewsPath + includeFile
			includeFilePath := filepath.Join(viewsPath, includeFile)
			_, e = t.New(includeFile).ParseFiles(includeFilePath)
			if e != nil {
				return e
			}
		}
	}
	e = t.Execute(w, data)
	if e != nil {
		return e
	}
	return nil
}

func (r *WebContext) WriteJson(data interface{}) error {
	w := r.Writer
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(r.Writer).Encode(data)
}

func (r *WebContext) WriteError(e error) error {
	if e != nil {
		errorString := e.Error()
		hr := r.Request
		r.Server.Log().Error(fmt.Sprintf("%s %s Error: %s", hr.URL.String(), hr.RemoteAddr, errorString))
		fmt.Fprint(r.Writer, errorString)
	}
	return e
}
