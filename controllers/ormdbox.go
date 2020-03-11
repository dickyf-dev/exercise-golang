package controllers
//
import (
	"github.com/eaciit/knot/knot.v1"
	//db "github.com/eaciit/dbox"
)
type OrmDboxController struct {
}
//func PrepareConnection() (db.IConnection, error) {
//	ci := &db.ConnectionInfo{"localhost:27017", "belajar_golang", "", "", nil}
//	c, e := db.NewConnection("mongo", ci)
//
//	if e != nil {
//		return nil, e
//	}
//
//	e = c.Connect()
//	if e != nil {
//		return nil, e
//	}
//
//	return c, nil
//}
func (w *OrmDboxController) OrmDbox(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "ormdbox/ormdbox.html"

	return nil
}
//func (w *OrmDboxController) GetJson(r *knot.WebContext) interface{} {
//	// set output type ke knot.OutputJson
//	r.Config.OutputType = knot.OutputJson
//	//return nil
//	payload := struct {
//		Name     string
//		Birthday time.Time
//		Parents  []string
//	}{}
//	e := r.GetPayload(&payload)
//	if e != nil {
//		toolkit.Println(e.Error())
//	}
//	ctx, _ := prepareContext()
//	ageNow := time.Now().Year() - payload.Birthday.Year()
//	u := models.NewDataUserModel()
//	u.ID = toolkit.RandomString(10)
//	u.Name = payload.Name
//	u.Birthday = payload.Birthday
//	u.Age = ageNow
//	u.Parents = payload.Parents
//	e = ctx.Save(u)
//	if e != nil {
//		return e.Error()
//	}
//	return "success"
//}
