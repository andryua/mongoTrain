package helpers

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ListGP struct {
	_id         bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name"`
	Type        string        `bson:"type"`
	Description string        `bson:"description"`
	Dependency  string        `bson:"dependency"`
}

func DownloadGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	name := r.URL.Query().Get("name")
	c := session.DB("gp").C("gpsel")
	defer session.Close()
	gplst := []AllPoliciesBson{}
	s := ""
	err := c.Find(bson.M{"gpname": name}).All(&gplst)
	if err != nil {
		fmt.Printf("find from delete fail %v\n", err)
	}
	for _, data := range gplst {
		scope := ""
		if strings.ToLower(data.Class) == "machine" {
			scope = "Computer"
		} else {
			scope = data.Class
		}
		for _, val := range data.Values {
			if val.SelectedValue == "" {
				val.SelectedValue = "none"
			}
			s += scope + "\n" + val.Key + "\n" + val.ValueName + "\n" + val.Type + ":" + val.SelectedValue + "\n\n"
		}
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+name+".log")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(s))
}

func SendId(w http.ResponseWriter, r *http.Request) {
	ids := ""
	name := ""
	s := []string{}
	session := ConnectDB()
	c := session.DB("gp").C("gpall")
	rec := session.DB("gp").C("gpsel")
	f := session.DB("gp").C("gplist")
	defer session.Close()
	var rr = AllPoliciesBson{}
	if r.Method == "POST" {
		r.ParseForm()
		ids = r.FormValue("ids")
		name = r.FormValue("gpname")
		s = strings.Split(ids, ",")
	}
	for _, id := range s {
		id_int, err := strconv.Atoi(id)
		err = c.Find(bson.M{"id": id_int}).One(&rr)
		if err != nil {
			log.Fatal("find in db:", err)
		}
		tmp := ListGP{}
		err = f.Find(bson.M{"name": name}).One(&tmp)
		rr.GpType = tmp.Type
		rr.GpName = tmp.Name
		exist, _ := rec.Find(bson.M{"name": rr.Name, "class": rr.Class, "gpname": rr.GpName, "gptype": rr.GpType}).Count()
		if exist == 0 {
			//fmt.Println(rr)
			err = rec.Insert(&rr)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	//fmt.Println(name)
	http.Redirect(w, r, "/edit?name="+name, 301)
}
