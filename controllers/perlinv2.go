package controllers

import (
	"fmt"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"strconv"
	"time"
)

type Perlinv2Controller struct {
}
type Perlin struct {
	Login 		string			`json:"login"`
	Name    	string    		`json:"name"`
	Company 	string			`json:"company"`
	CreateAt	time.Time		`json:"create_at"`
	UpdateAt	time.Time		`json:"update_at"`
}
func (w *Perlinv2Controller) Perlinv2(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "perlinv2/goroutine.html"

	return nil
}
func (w *Perlinv2Controller) Knockout(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "datetime/ajax-knockout.html"

	return nil
}

func (w *Perlinv2Controller) GetJson(r *knot.WebContext) interface{} {
	// set output type ke knot.OutputJson
	r.Config.OutputType = knot.OutputJson

	return time.Now().Format("2006-01-02T15:04Z")
}
func (w *Perlinv2Controller) Save(r *knot.WebContext) interface{} {
// set output type ke knot.OutputJson
r.Config.OutputType = knot.OutputJson
	payload := struct {
		X   string
		Y	string
	}{}
	e := r.GetPayload(&payload)
	if e != nil {
		toolkit.Println(e.Error())
	}
	var ix float64
	var iy float64
	var y, _ = strconv.ParseFloat(payload.Y, 100)
	var x, _ = strconv.ParseFloat(payload.X, 100)

	//var n0, n1, ix0, ix1, value int

	//var lerp = (1.0 - 11) * 30 + 11 * 20
	//var dx = x - ix
	//var dy = y - iy
	//var dot = (dx * ix * iy) + (dy * ix * iy)

	var x0 = x;
	var x1 = x0 + 1;
	var y0 = y;
	var y1 = y0 + 1;

	var sx = x - x0
	var sy = y - y0

	var n0 = ((x-ix) * ix * iy) + ((y-iy) * ix * iy)
	var n1 = ((x1-ix) * ix * iy) + ((y-iy) * ix * iy)
	var ix0 = (100 - sx) * n0 + sx * n1

	var n00 = ((x-ix) * ix * y1) + ((y-y1) * ix * y1)
	var n11 = ((x1-ix) * ix * y1) + ((y-y1) * ix * y1)
	var ix1 = (1.0 - sx) * n00 + sx * n11

	var value = (1.0 - sx) * ix0 + sy * ix1
	fmt.Println(value)
	return value
}