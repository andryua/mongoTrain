package main

import (
	"fmt"
	"log"
	"mongoTrain/helpers"
	"net/http"
)

var result []helpers.AllPolicies

func admjson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(helpers.Treegen(result)))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	dataPolicies, lang, dataCat, cataloguesName, present := helpers.ParseFiles()
	//fmt.Printf("%v\n",dataCat)
	cataloguePath := helpers.CategoriesPath(dataCat, cataloguesName)
	result = helpers.PoliciesParse(dataPolicies, lang, cataloguePath, present)
	/*
		jsonRes, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		file1 := "gpo.json"
		ioutil.WriteFile(file1, jsonRes, 0777)
		jsonTree := helpers.Treegen(result)
		file2 := "gpoTree.json"
		ioutil.WriteFile(file2, []byte(jsonTree), 0777)
	*/
	session := helpers.ConnectDB()
	c := session.DB("gp").C("gpall")
	n, err := c.Count()
	defer session.Close()
	if n > 0 {
		err := c.DropCollection()
		if err != nil {
			fmt.Println(err)
		}
	}
	helpers.AllgpToBson(c, result)

	http.HandleFunc("/admjson", admjson)
	http.HandleFunc("/add", helpers.GPTree)
	http.HandleFunc("/addgp", helpers.AddGP)
	http.HandleFunc("/showaddgp", helpers.ShowAddGP)
	http.HandleFunc("/showeditgp", helpers.ShowEditGP)
	http.HandleFunc("/", helpers.GPList)
	http.HandleFunc("/sendids", helpers.SendId)
	http.HandleFunc("/edit", helpers.EditGP)
	http.HandleFunc("/deletegp", helpers.DeleteGP)
	http.HandleFunc("/deleterule", helpers.DeleteRule)
	http.HandleFunc("/download", helpers.DownloadGP)
	http.HandleFunc("/editgp", helpers.EditParamGP)
	http.HandleFunc("/update", helpers.UpdateRule)
	http.HandleFunc("/copy", helpers.CopyRule)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))) //погашення папки (щоб при роботі сервера він знав де брати файли для вебу)

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Println("Error on ListenAndServe:\n")
		log.Fatal(err)
	}
}
