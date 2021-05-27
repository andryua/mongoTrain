package helpers

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func DeleteGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")
	gpname := r.URL.Query().Get("gpname")
	fmt.Println(name, " and ", gpname)
	gplst := session.DB("gp").C("gplist")
	gpsel := session.DB("gp").C("gpsel")
	defer session.Close()
	if gpname != "" && id != "" {
		err := gpsel.Remove(bson.M{"id": id, "gpname": gpname})
		if err != nil {
			fmt.Printf("remove fail %v\n", err)
		}
		http.Redirect(w, r, "/edit?name="+gpname, 301)
	}
	if name != "" {
		err := gplst.Remove(bson.M{"name": name})
		_, err = gpsel.RemoveAll(bson.M{"gpname": name})
		if err != nil {
			fmt.Printf("remove fail %v\n", err)
		}
		http.Redirect(w, r, "/", 301)
	}
}
