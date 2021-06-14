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
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		lstgp.Name = r.FormValue("gpname")
		lstgp.Type = r.FormValue("gptype")
		lstgp.Description = r.FormValue("gpinfo")
		lstgp.Dependency = r.Form["gpdepend"]
	}
	//fmt.Println(depend)
	err := c.Insert(lstgp)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

func DownloadGP(w http.ResponseWriter, r *http.Request) {
	session := ConnectDB()
	name := r.URL.Query().Get("name")
	selectGP := session.DB("gp").C("gpsel")
	gpList := session.DB("gp").C("gplist")
	defer session.Close()
	var selectedGP []AllPoliciesBson
	var dependenciesGP []AllPoliciesBson
	var additionalSelectedGP []AllPoliciesBson
	var gpNew []AllPoliciesBson
	var additionalGP ListGP
	var currentGP ListGP
	var j []string
	var gpname_ID = make(map[string]string)
	s := ""
	err := gpList.Find(bson.M{"name": name}).One(&currentGP)
	err = selectGP.Find(bson.M{"gpname": name}).All(&selectedGP)
	if err != nil {
		fmt.Printf("fail %v\n", err)
	}
	switch currentGP.Type {
	case "default":
		if len(currentGP.Dependency) > 0 {
			for _, dependence := range currentGP.Dependency {
				dependenciesGP = nil
				err = selectGP.Find(bson.M{"gpname": dependence}).All(&dependenciesGP)
				if err != nil {
					fmt.Printf("fail %v\n", err)
				}
				selectedGP = append(selectedGP, dependenciesGP...)
			}
		} else {
			gpNew = selectedGP
		}
	case "users":
		if len(currentGP.Dependency) > 0 {
			for _, dependence := range currentGP.Dependency {
				dependenciesGP = nil
				err = selectGP.Find(bson.M{"gpname": dependence}).All(&dependenciesGP)
				if err != nil {
					fmt.Printf("fail %v\n", err)
				}
				selectedGP = append(selectedGP, dependenciesGP...)
				err := gpList.Find(bson.M{"name": dependence}).One(&additionalGP)
				if len(additionalGP.Dependency) > 0 {
					for _, depend := range additionalGP.Dependency {
						dependenciesGP = nil
						err = selectGP.Find(bson.M{"gpname": depend}).All(&additionalSelectedGP)
						if err != nil {
							fmt.Printf("fail %v\n", err)
						}
						selectedGP = append(selectedGP, additionalSelectedGP...)
					}
				}
			}
		} else {
			gpNew = selectedGP
		}
	case "main":
		gpNew = selectedGP
	}

	for _, data := range selectedGP {
		if data.GpName == name {
			gpname_ID[data.Name] = data.Class
		}
	}

	for key, value := range gpname_ID {
		for _, data := range selectedGP {
			if data.GpName != name && data.Name == key && data.Class == value {
				j = append(j, data.ID)
			}
		}
	}
	j = RemoveDuplicateStr(j)
	for _, data := range selectedGP {
		if !Contains(j, data.ID) && (data.GpType == "users" || data.GpType == "default") {
			gpNew = append(gpNew, data)
		}
	}
	for _, data := range gpNew {
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
			fmt.Println(val.SelectedValue)
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
	var s []string
	session := ConnectDB()
	c := session.DB("gp").C("gpall")
	rec := session.DB("gp").C("gpsel")
	f := session.DB("gp").C("gplist")
	defer session.Close()
	var rr = AllPoliciesBson{}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		ids = r.FormValue("ids")
		name = r.FormValue("gpname")
		s = strings.Split(ids, ",")
	}
	for _, id := range s {
		id_int, err := strconv.Atoi(id)
		err = c.Find(bson.M{"IDtmp": id_int}).One(&rr)
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
			//rr.ID = bson.NewObjectId().Hex()
			err = rec.Insert(&rr)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	//fmt.Println(name)
	http.Redirect(w, r, "/edit?name="+name, 301)
}

func CopyRule(w http.ResponseWriter, r *http.Request) {
	//зміна не лише ід, але і назви гп і її типу
	session := ConnectDB()
	id := r.URL.Query().Get("id")
	gpname := r.URL.Query().Get("gpname")
	gpsel := session.DB("gp").C("gpsel")
	gplist := session.DB("gp").C("gplist")
	defer session.Close()
	var cpRule AllPoliciesBson
	var currGP ListGP
	err := gpsel.Find(bson.M{"id": id}).One(&cpRule)
	err = gplist.Find(bson.M{"name": gpname}).One(&currGP)
	if err != nil {
		fmt.Printf("find fail: %v\n", err)
	}
	cpRule.GpName = gpname
	cpRule.GpType = currGP.Type
	cpRule.ID = bson.NewObjectId().Hex()
	//cpRule.Dependencies = currGP.Dependency
	for _, val := range cpRule.Values {
		val.SelectedValue = ""
		val.Notes = ""
	}
	err = gpsel.Insert(cpRule)
	if err != nil {
		fmt.Printf("copy rule fail: %v\n", err)
	}
	http.Redirect(w, r, "/edit?name="+gpname, 301)
}
