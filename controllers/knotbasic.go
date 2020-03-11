package controllers

import (
	knot "github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"time"
)

type KnotBasicController struct {
}

func (w *KnotBasicController) ResAjaxs(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "knotbasic/knotbasic.html"

	return nil
}
func (w *KnotBasicController) ResKnockout(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "knotbasic/ajax-knockout.html"

	return nil
}

func (w *KnotBasicController) ResJson(r *knot.WebContext) interface{} {
	// set output type ke knot.OutputJson
	r.Config.OutputType = knot.OutputJson
	sample := []toolkit.M{
		{"date": time.Now().Format("2006-01-02T15:04:05Z")},
	}
	return sample
}
