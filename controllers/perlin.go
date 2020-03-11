package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"strconv"
	"sync"
)

type PerlinController struct {
}

func (w *PerlinController) Perlin(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "perlin/perlin.html"
	return nil
}

func jumlah(x, y float64, data *[]float64, wg *sync.WaitGroup) {
	hasil := x + y
	*data = append(*data, hasil)
	wg.Done()
}

func (w *PerlinController) Create(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := struct {
		First   string
		Second  string
	}{}
	e := r.GetPayload(&payload)
	if e != nil {
		toolkit.Println(e.Error())
	}
	toolkit.Println(payload.Second)
	toolkit.Println(payload.First)
	result := []float64{}

	startPoint, _ := strconv.ParseFloat(payload.First, 64)
	endPoint, _ := strconv.ParseFloat(payload.Second, 64)

	total := int(startPoint * endPoint)

	var wg sync.WaitGroup

	wg.Add(total)
	for i := 0.0; i < startPoint; i++ {
		for j := 0.0; j < endPoint; j++ {
			go jumlah(i, j, &result, &wg)
		}
	}

	wg.Wait()
	return result
}