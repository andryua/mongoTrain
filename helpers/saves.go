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
	Dependency  []string      `bson:"dependency"`
}

func AddGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	c := session.DB("gp").C("gplist")
	lstgp := ListGP{}
	defer session.Close()
	//var depend []string
	if r.Method == "POST" {
		r.ParseForm()
		lstgp.Name = r.FormValue("gpname")
		lstgp.Type = r.FormValue("gptype")
		lstgp.Description = r.FormValue("gpinfo")
		lstgp.Dependency = r.Form["gpdepend"]
	}
	//fmt.Println(depend)
	c.Insert(lstgp)
	http.Redirect(w, r, "/", 301)
}

func DownloadGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	name := r.URL.Query().Get("name")
	c := session.DB("gp").C("gpsel")
	gpl := session.DB("gp").C("gplist")
	defer session.Close()
	gplst := []AllPoliciesBson{}
	gpdep := []AllPoliciesBson{}
	gpls := ListGP{}
	s := ""
	err := gpl.Find(bson.M{"name": name}).One(&gpls)
	err = c.Find(bson.M{"gpname": name}).All(&gplst)
	if err != nil {
		fmt.Printf("fail %v\n", err)
	}
	if len(gpls.Dependency) > 0 {
		for _, dependence := range gpls.Dependency {
			gpdep = nil
			err = c.Find(bson.M{"gpname": dependence}).All(&gpdep)
			if err != nil {
				fmt.Printf("fail %v\n", err)
			}
			gplst = append(gplst, gpdep...)
		}
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

/*
func CopyRule(w http.ResponseWriter, r *http.Request) {
	//зміна не лише ід, але і назви гп і її типу
	session := ConnectDB()
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Printf("atoi fail: %v\n", err)
	}
	gpname := r.URL.Query().Get("gpname")
	gpsel := session.DB("gp").C("gpsel")
	defer session.Close()
	cpDoc := AllPoliciesBson{}
	err = gpsel.Find(bson.M{"id": id, "gpname": gpname}).All(&cpDoc)
	if err != nil {
		fmt.Printf("find fail: %v\n", err)
	}
	if gpname != "" {
		cpDoc.GpName = gpname
		if cpDoc.ID == -1 {
			cpDoc.ID--
		} else {
			cpDoc.ID = -1
		}
		err = gpsel.Insert(cpDoc)
		if err != nil {
			fmt.Printf("insert fail: %v\n", err)
		}
		http.Redirect(w, r, "/edit?name="+gpname, 301)
	}
}
*/
