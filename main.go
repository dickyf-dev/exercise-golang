package main

import (
	"knot_example/controllers"
	"net/http"
	"os"

	knot "github.com/eaciit/knot/knot.v1"
)

var (
	appViewsPath = (func(dir string, _ error) string { return dir + "/views/" }(os.Getwd()))
)

func main() {
	app := knot.NewApp("")
	app.ViewsPath = appViewsPath

	app.Register(new(controllers.MessageController))
	app.Register(new(controllers.KnotBasicController))
	app.Register(new(controllers.DateTimeController))
	app.Register(new(controllers.OrmDboxController))
	app.Register(new(controllers.FetchGithubController))
	app.Register(new(controllers.FileUploadController))
	app.Register(new(controllers.PerlinController))
	app.Register(new(controllers.Perlinv2Controller))
	app.Register(new(controllers.GoroutineController))

	app.LayoutTemplate = "_template.html"
	app.DefaultOutputType = knot.OutputTemplate

	knot.RegisterApp(app)
	otherRoutes := map[string]knot.FnContent{
		"/": func(r *knot.WebContext) interface{} {
			http.Redirect(r.Writer, r.Request, "/fetchgithub/fetchgithub", http.StatusTemporaryRedirect)
			return true
		},
	}

	knot.StartAppWithFn(app, "localhost:8999", otherRoutes)

}
