package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"time"
)

type DateTimeController struct {
}

func (w *DateTimeController) ResAjax(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "datetime/datetime.html"

	return nil
}
func (w *DateTimeController) Knockout(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "datetime/ajax-knockout.html"

	return nil
}

func (w *DateTimeController) GetJson(r *knot.WebContext) interface{} {
	// set output type ke knot.OutputJson
	r.Config.OutputType = knot.OutputJson

	return time.Now().Format("2006-01-02T15:04Z")
}

func (w *DateTimeController) Save(r *knot.WebContext) interface{} {
// set output type ke knot.OutputJson
r.Config.OutputType = knot.OutputJson

	payload := struct {
		Time   string `json:"time"`
	}{}
	e := r.GetPayload(&payload)
	if e != nil {
		toolkit.Println(e.Error())
	}
	toolkit.Println(payload.Time)
	// process to change date
	sample := []toolkit.M{
		{"dateNow": time.Now().Format("2006-01-02T15:04Z")},
		{"dateAfter" : time.Now().Add(12*time.Hour).Format("2006-01-02T15:04Z")},
	}
	return sample
}