package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io/ioutil"
	"log"
	"mongoTrain/helpers"
	"net/http"
	"strconv"
	"strings"

	//"os"
	"time"
)

type ListGP struct {
	_id         bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name"`
	Type        string        `bson:"type"`
	Description string        `bson:"description"`
	Dependency  string        `bson:"dependency"`
}

func connectDB() mgo.Session {
	MongoDBHosts := "127.0.0.1"
	dialInfo := &mgo.DialInfo{
		Addrs: []string{MongoDBHosts},
		//Username: "admin",
		//Password: "password",
		Timeout:  60 * time.Second,
		Database: "admin",
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	return *session
}

func admjson(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("gpoTree.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(content)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GPTree(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	t := template.Must(template.ParseFiles("./templates/gptree.html"))
	var v = make(map[string]interface{})
	v["Name"] = name
	t.ExecuteTemplate(w, "gptree", v)
}

func EditGP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	//fmt.Println(name)
	session := connectDB()
	c := session.DB("gp").C("gpsel")
	var sh = []helpers.AllPoliciesBson{}
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

func SendId(w http.ResponseWriter, r *http.Request) {
	ids := ""
	name := ""
	s := []string{}
	session := connectDB()
	c := session.DB("gp").C("gpall")
	rec := session.DB("gp").C("gpsel")
	f := session.DB("gp").C("gplist")
	defer session.Close()
	var rr = helpers.AllPoliciesBson{}
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

func GPList(w http.ResponseWriter, r *http.Request) {
	session := connectDB()
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

func AddGP(w http.ResponseWriter, r *http.Request) {
	session := connectDB()
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

func main() {
	session := connectDB()
	c := session.DB("gp").C("gpall")
	n, err := c.Count()
	defer session.Close()
	if n > 0 {
	} else {
		helpers.GetAllgp(c)
	}

	http.HandleFunc("/admjson", admjson)
	http.HandleFunc("/add", GPTree)
	http.HandleFunc("/addgp", AddGP)
	http.HandleFunc("/showaddgp", ShowAddGP)
	http.HandleFunc("/", GPList)
	http.HandleFunc("/sendids", SendId)
	http.HandleFunc("/edit", EditGP)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))) //погашення папки (щоб при роботі сервера він знав де брати файли для вебу)

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Println("Error on ListenAndServe:\n")
		log.Fatal(err.Error())
	}
}
