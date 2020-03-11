package models

import (
	"time"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	)
type DataUserModel struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID        string    `json:"ID" bson:"_id"`
	Name      string    `json:"Name" bson:"Name"`
	Age 	  int		`json:"Age" bson:"Age"`
	Birthday  time.Time `json:"Birthday" bson:"Birthday"`
	Parents   []string  `json:"Parents" bson:"Parents"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
}
func NewDataUserModel() *DataUserModel {
	m := new(DataUserModel)
	m.ID = tk.RandomString(32)
	return m

}
func (u *DataUserModel) TableName() string {
	return "datausers"
}

func (u *DataUserModel) RecordID() interface{} {
	return u.ID

}