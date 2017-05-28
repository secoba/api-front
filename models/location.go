package models

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/orm"
	"time"
)

// Location  具体的一个api接口
type Location struct {
	ID       int64  `orm:"pk;auto;column(location_id)"`
	NodeId   int64     `orm:"column(node_id)"`
	
	Location string `orm:"size(255);column(location)"` // 转发路径  eg / ，/api/ ， /myapi/subscribe
	
	Status   int    `orm:"column(status)"`
	Name     string `orm:"size(64);column(name)"`
	Intro    string `orm:"size(10000);column(intro)"`

	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(ctime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(mtime)"`

	ServerNode *ServerNode `orm:"-"`

	Backends []*Backend `orm:"-"`
}

func (l *Location) TableName() string {
	return "api_location"
}

func (l *Location) String() string {
	bs, _ := json.Marshal(l)
	return string(bs)
}

func (l *Location) TableUnique() [][]string {
	return [][]string{
		{"NodeId", "Location"},
	}
}

func (g *Location) Create() (id int64, err error) {
	o := orm.NewOrm()
	return o.Insert(g)
}

func (g *Location) Read() error {
	o := orm.NewOrm()
	return o.Read(g)
}

func (g *Location) Delete() (num int64, err error) {
	o := orm.NewOrm()
	return o.Delete(g)
}

func (g *Location) Update() error {
	if g.ID < 1 {
		return errors.New("id is required")
	}
	o := orm.NewOrm()
	_, err := o.Update(g)
	if err != nil {
		return err
	}
	return nil
}
func (g *Location) GetBackends() []*Backend {
	cond := orm.NewCondition().And("location_id", g.ID)
	ls := ListAllBackend(cond)
	g.Backends = ls
	return ls
}

func (g *Location) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(g)
}

func ListAllLocation(cond *orm.Condition) []*Location {
	var ls []*Location
	query := new(Location).Query().SetCond(cond)
	query.All(&ls)
	return ls
}