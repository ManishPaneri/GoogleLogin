package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type GoogleUserTable struct {
	Id            int    `orm:"column(id);auto"`
	Name          string `orm:"column(name);size(225)"`
	Email         string `orm:"column(email);size(225)"`
	Role          string `orm:"column(role)"`
	Picture       string `orm:"column(picture)"`
	Hd            string `orm:"column(hd)"`
	VerifiedEmail int8   `orm:"column(verified_email)"`
	AuthId        string `orm:"column(auth_id);size(256)"`
}

func (t *GoogleUserTable) TableName() string {
	return "Google_User_Table"
}

func init() {
	orm.RegisterModel(new(GoogleUserTable))
}

// AddGoogleUserTable insert a new GoogleUserTable into database and returns
// last inserted Id on success.
func AddGoogleUserTable(m *GoogleUserTable) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetGoogleUserTableById retrieves GoogleUserTable by Id. Returns error if
// Id doesn't exist
func GetGoogleUserTableById(id int) (v *GoogleUserTable, err error) {
	o := orm.NewOrm()
	v = &GoogleUserTable{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//GetAllGoogleUserTable retrieves all GoogleUserTable matches certain condition. Returns empty list if
// no records exist
func GetAllGoogleUserTable() (v []GoogleUserTable, err error) {
	o := orm.NewOrm()
	v = []GoogleUserTable{}
	_, err = o.QueryTable(new(GoogleUserTable)).All(&v)
	fmt.Println(v)
	return v, err
}

//GetGoogleUserTableByAuthId get user row by authID
func GetGoogleUserTableByAuthId(authid string) (v *GoogleUserTable, err error) {

	o := orm.NewOrm()
	v = &GoogleUserTable{}
	err = o.QueryTable(new(GoogleUserTable)).Filter("authid", authid).One(v)
	if err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateGoogleUserTableById updates GoogleUserTable by Id and returns error if
// the record to be updated doesn't exist
func UpdateGoogleUserTableById(m *GoogleUserTable) (err error) {
	o := orm.NewOrm()
	v := GoogleUserTable{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteGoogleUserTable deletes GoogleUserTable by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGoogleUserTable(id int) (err error) {
	o := orm.NewOrm()
	v := GoogleUserTable{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GoogleUserTable{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
