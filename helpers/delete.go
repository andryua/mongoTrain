package helpers

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func DeleteGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	name := r.URL.Query().Get("name")
	gplst := session.DB("gp").C("gplist")
	gpsel := session.DB("gp").C("gpsel")
	defer session.Close()
	err := gplst.Remove(bson.M{"name": name})
	_, err = gpsel.RemoveAll(bson.M{"gpname": name})
	if err != nil {
		fmt.Printf("remove fail %v\n", err)
	}
	http.Redirect(w, r, "/", 301)
}

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	name := r.URL.Query().Get("name")
	gpname := r.URL.Query().Get("gpname")
	scope := r.URL.Query().Get("class")
	gpsel := session.DB("gp").C("gpsel")
	defer session.Close()
	err := gpsel.Remove(bson.M{"name": name, "gpname": gpname, "class": scope})
	if err != nil {
		fmt.Printf("remove fail: %v\n", err)
	}
	http.Redirect(w, r, "/edit?name="+gpname, 301)
}
