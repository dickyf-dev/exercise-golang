package controllers

import (
	"fmt"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"time"
)

type GoroutineController struct {
}

func (w *GoroutineController) Goroutine(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "goroutine/goroutine.html"

	return nil
}
func (w *GoroutineController) Knockout(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "datetime/ajax-knockout.html"

	return nil
}

func (w *GoroutineController) GetJson(r *knot.WebContext) interface{} {
	// set output type ke knot.OutputJson
	r.Config.OutputType = knot.OutputJson

	return time.Now().Format("2006-01-02T15:04Z")
}
func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}
func (w *GoroutineController) Save(r *knot.WebContext) interface{} {
// set output type ke knot.OutputJson
r.Config.OutputType = knot.OutputJson

	payload := struct {
		First   string
		Second   string
	}{}
	e := r.GetPayload(&payload)
	if e != nil {
		toolkit.Println(e.Error())
	}
	// process to change date
	go say(payload.First)
	say(payload.Second)
	//return
	return nil
}