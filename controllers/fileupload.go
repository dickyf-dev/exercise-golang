package controllers

import (
	"github.com/eaciit/knot/knot.v1"
)

type FileUploadController struct {
}

func (w *FileUploadController) FileUpload(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "fileupload/fileupload.html"

	return nil
}
func (w *FileUploadController) Knockout(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "datetime/ajax-knockout.html"

	return nil
}

func (w *FileUploadController) GetJson(r *knot.WebContext) interface{} {
	// set output type ke knot.OutputJson
	r.Config.OutputType = knot.OutputJson
	return nil
}