package helpers

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	//"strconv"
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
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
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
	id := r.URL.Query().Get("id")
	//fmt.Println(id)
	c := session.DB("gp").C("gpsel")
	var tmp AllPoliciesBson
	//var tmpVals Values
	err := c.Find(bson.M{"id": id}).One(&tmp)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	//fmt.Println(tmp)
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		var mVN []string
		var mV []string
		var mN []string
		for _, val := range tmp.Values {
			for key, value := range r.Form {
				if key == "id" {
					continue
				}
				if (val.ValueName == key) && (val.Manual != true) {
					val.SelectedValue = strings.Join(value, "|")
					//fmt.Println(key, "\t", val.SelectedValue)
					err := c.Update(bson.M{"id": id, "values.valueName": key}, bson.M{"$set": bson.M{"values.$.selectedvalue": val.SelectedValue}})
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
		if len(mVN) > 0 {
			fmt.Println(mVN)
			/*	for _, x := range tmp.Values {
					if x.Manual == true {
						tmpVals = x
						break
					}
				}
			*/for _, val := range tmp.Values {
				if val.Manual == true {
					/*for i, _ := range mVN {
						tmpVals.ValueName = "manual-"+strconv.Itoa(i)
						err = c.Update(bson.M{"id": id}, bson.M{"$push": bson.M{"values": tmpVals}})
						if err != nil {
							log.Println(err)
						}
					}*/
					for i, _ := range mVN {
						j, err := c.Find(bson.M{"id": id, "values.valueName": mVN[i]}).Count()
						if err != nil {
							log.Println(err)
						}
						if j > 0 {
							err = c.Update(bson.M{"id": id, "values.valueName": mVN[i]}, bson.M{"$set": bson.M{"values.$.selectedvalue": mV[i], "values.$.notes": mN[i]}})
							if err != nil {
								log.Println(err)
							}
						} else {
							val.ValueName = mVN[i]
							val.SelectedValue = mV[i]
							val.Notes = mN[i]
							err = c.Update(bson.M{"id": id}, bson.M{"$push": bson.M{"values": val}})
							if err != nil {
								log.Println(err)
							}
						}
					}
				}
			}
			err = c.Update(bson.M{"id": id}, bson.M{"$pull": bson.M{"values": bson.M{"values.valueName": "manual"}}})
			if err != nil {
				log.Println(err)
			}
		}
	}
	w.Write([]byte("saved"))
}
