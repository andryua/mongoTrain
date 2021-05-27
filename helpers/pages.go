package helpers

import (
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
	var sh = []AllPoliciesBson{}
	err := c.Find(bson.M{"gpname": name}).All(&sh)
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

func AddGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	c := session.DB("gp").C("gplist")
	lstgp := ListGP{}
	if r.Method == "POST" {
		r.ParseForm()
		lstgp.Name = r.FormValue("gpname")
		lstgp.Type = r.FormValue("gptype")
		lstgp.Description = r.FormValue("gpinfo")
		lstgp.Dependency = r.FormValue("gpdepend")
	}
	c.Insert(lstgp)
	http.Redirect(w, r, "/", 301)
}

func ShowAddGP(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/addgp.html"))
	t.ExecuteTemplate(w, "addgp", "")
}

func GPTree(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	t := template.Must(template.ParseFiles("./templates/gptree.html"))
	var v = make(map[string]interface{})
	v["Name"] = name
	t.ExecuteTemplate(w, "gptree", v)
}
