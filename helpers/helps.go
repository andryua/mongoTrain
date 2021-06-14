package helpers

import (
	"gopkg.in/mgo.v2"
	"time"
)

func ConnectDB() mgo.Session {
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

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func RemoveIndex(s []AllPoliciesBson, index int) []AllPoliciesBson {
	if index == len(s)-1 {
		return s[:index]
	} else {
		return append(s[:index], s[index+1:]...)
	}
}

func RemoveDuplicateInt(intSlice []int) []int {
	keys := make(map[int]bool)
	var list []int

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func RemoveDuplicateStr(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
