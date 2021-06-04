package helpers

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strings"
)

func EditParamGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	name := r.URL.Query().Get("name")
	c := session.DB("gp").C("gplist")
	lstgp := ListGP{}
	defer session.Close()
	//var depend []string
	err := c.Find(bson.M{"name": name}).One(&lstgp)
	if err != nil {
		log.Panic(err)
	}
	if r.Method == "POST" {
		r.ParseForm()
		lstgp.Description = r.FormValue("gpinfo")
		lstgp.Dependency = r.Form["gpdepend"]
	}
	//fmt.Println(depend)
	err = c.Update(bson.M{"name": lstgp.Name}, bson.M{"$set": bson.M{"description": lstgp.Description, "dependency": lstgp.Dependency}})
	if err != nil {
		log.Panic(err)
	}
	//c.Insert(lstgp)
	http.Redirect(w, r, "/edit?name="+name, 301)
}

func UpdateRule(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	name := r.URL.Query().Get("name")
	gpname := r.URL.Query().Get("gpname")
	scope := r.URL.Query().Get("class")
	fmt.Println(name, "\t", gpname, "\t", scope)
	c := session.DB("gp").C("gpsel")
	var tmp AllPoliciesBson
	var tmpVals Values
	err := c.Find(bson.M{"gpname": gpname, "name": name, "class": scope}).One(&tmp)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	//fmt.Println(tmp)
	if r.Method == "POST" {
		r.ParseForm()
		var mVN []string
		var mV []string
		var mN []string
		for _, val := range tmp.Values {
			for key, value := range r.Form {
				if key == "name" || key == "gpname" || key == "class" {
					continue
				}
				if (val.ValueName == key) && (val.Manual != true) {
					val.SelectedValue = strings.Join(value, "|")
					//fmt.Println(key, "\t", val.SelectedValue)
					err := c.Update(bson.M{"gpname": gpname, "name": name, "class": scope, "values.valueName": key}, bson.M{"$set": bson.M{"values.$.selectedvalue": val.SelectedValue}})
					if err != nil {
						log.Println(err)
					}
				}
				if val.Manual == true {
					switch key {
					case "manualValueName":
						mVN = value
					case "manualValue":
						mV = value
					case "manualDescription":
						mN = value
					default:
						fmt.Println("key: ", key, "\tvalue: ", value)
					}
				}
			}
		}
		fmt.Println(mVN)
		for _, x := range tmp.Values {
			if x.Manual == true {
				tmpVals = x
				break
			}
		}
		if len(mVN) > 0 {
			for _, val := range tmp.Values {
				if val.Manual == true {
					for i, _ := range mVN {
						err = c.Update(bson.M{"gpname": gpname, "name": name, "class": scope}, bson.M{"$push": bson.M{"values": tmpVals}})
						if err != nil {
							log.Println(err)
						}
						err = c.Update(bson.M{"gpname": gpname, "name": name, "class": scope, "values.notes": "new"}, bson.M{"$push": bson.M{"values": tmp.Values, "values.$.valueName": mVN[i], "values.$.selectedvalue": mV[i], "values.$.notes": mN[i]}})
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
			err = c.Update(bson.M{"gpname": gpname, "name": name, "class": scope}, bson.M{"$pull": bson.M{"values.valueName": "manual"}})
			if err != nil {
				log.Println(err)
			}
		}
	}
	w.Write([]byte("saved"))
}
