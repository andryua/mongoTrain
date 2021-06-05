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
	id := r.URL.Query().Get("id")

	var tmp AllPoliciesBson
	gpsel := session.DB("gp").C("gpsel")
	err := gpsel.Find(bson.M{"id": id}).One(&tmp)
	//fmt.Println(tmp)
	defer session.Close()
	err = gpsel.Remove(bson.M{"id": id})
	if err != nil {
		fmt.Printf("remove fail: %v\n", err)
	}
	http.Redirect(w, r, "/edit?name="+tmp.GpName, 301)
}
