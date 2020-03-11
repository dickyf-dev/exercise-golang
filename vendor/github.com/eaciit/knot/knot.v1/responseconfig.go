package knot

import "github.com/eaciit/toolkit"

type OutputType int

const (
	OutputNone     OutputType = 1
	OutputTemplate OutputType = 2
	OutputHtml     OutputType = 10
	OutputJson     OutputType = 100
	OutputByte     OutputType = 1000
)

func (o OutputType) String() string {
	if o == OutputNone {
		return "None"
	} else if o == OutputTemplate {
		return "Template"
	} else if o == OutputHtml {
		return "HTML"
	} else if o == OutputJson {
		return "JSON"
	} else if o == OutputByte {
		return "Byte"
	}
	return "N/A"
}

type ResponseConfig struct {
	AppName          string
	App              *App
	ControllerName   string
	MethodName       string
	ViewName         string
	OutputType       OutputType
	IgnoreValidation bool
	LayoutTemplate   string
	ViewsPath        string
	IncludeFiles     []string
	NoLog            bool
	Headers          map[string]string

	data toolkit.M
}

func NewResponseConfig() *ResponseConfig {
	c := new(ResponseConfig)
	c.Headers = map[string]string{}
	c.IncludeFiles = []string{}
	c.OutputType = DefaultOutputType
	return c
}

func (r *ResponseConfig) Data(key string, def interface{}) interface{} {
	if r.data == nil {
		r.data = toolkit.M{}
	}
	return r.data.Get(key, def)
}

func (r *ResponseConfig) SetData(key string, value interface{}) {
	if r.data == nil {
		r.data = toolkit.M{}
	}
	r.data.Set(key, value)
}
