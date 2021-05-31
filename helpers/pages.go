package helpers

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
)

func GPList(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	var lstgp = []ListGP{}
	c := session.DB("gp").C("gplist")
	defer session.Close()
	n, err := c.Count()
	if err != nil {
		log.Fatal(err)
	}
	if n < 1 {
		http.Redirect(w, r, "/addgp", 301)
	} else {
		err := c.Find(bson.M{}).All(&lstgp)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(lstgp)
		t := template.Must(template.ParseFiles("./templates/index.html"))
		var v = make(map[string]interface{})
		v["GPList"] = lstgp
		t.ExecuteTemplate(w, "index", v)
	}
}

func EditGP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	//fmt.Println(name)
	session := ConnectDB()
	c := session.DB("gp").C("gpsel")
	gpl := session.DB("gp").C("gplist")
	var sh = []AllPoliciesBson{}
	var rules = []AllPoliciesBson{}
	var gplst = ListGP{}
	err := gpl.Find(bson.M{"name": name}).One(&gplst)
	err = c.Find(bson.M{"gpname": name}).All(&sh)
	if len(gplst.Dependency) > 0 {
		for _, dependence := range gplst.Dependency {
			rules = nil
			err = c.Find(bson.M{"gpname": dependence}).All(&rules)
			if err != nil {
				fmt.Printf("fail %v\n", err)
			}
			sh = append(sh, rules...)
		}
	}
	//fmt.Println(sh)
	if err != nil {
		log.Fatal("find in db:", err)
	}
	defer session.Close()
	t := template.Must(template.ParseFiles("./templates/edit.html"))
	var v = make(map[string]interface{})
	v["Gpname"] = name
	v["Rules"] = sh
	t.ExecuteTemplate(w, "editgp", v)
}

func ShowAddGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	var lstgp = []ListGP{}
	c := session.DB("gp").C("gplist")
	defer session.Close()
	err := c.Find(bson.M{}).All(&lstgp)
	if err != nil {
		log.Fatal(err)
	}
	var v = make(map[string]interface{})
	v["GPList"] = lstgp
	t := template.Must(template.ParseFiles("./templates/addgp.html"))
	err = t.ExecuteTemplate(w, "addgp", v)
	if err != nil {
		log.Fatal(err)
	}
}

func GPTree(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	t := template.Must(template.ParseFiles("./templates/gptree.html"))
	var v = make(map[string]interface{})
	v["Name"] = name
	t.ExecuteTemplate(w, "gptree", v)
}
