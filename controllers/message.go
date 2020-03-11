package controllers

import (
	knot "github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type MessageController struct {
}

func (w *MessageController) AjaxExample(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "message/goroutine.html"

	return nil
}
func (w *MessageController) AjaxKnockoutExample(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "message/ajax-knockout.html"

	return nil
}

func (w *MessageController) SampleJsonAjax(r *knot.WebContext) interface{} {
	// set output type ke knot.OutputJson
	r.Config.OutputType = knot.OutputJson

	sample := []toolkit.M{
		toolkit.M{"name": "noval"},
		toolkit.M{"name": "agung"},
	}

	// langsung kembalikan objek yg ingin di-return kan sebagai json
	return sample

}
