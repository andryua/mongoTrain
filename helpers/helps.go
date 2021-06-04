package helpers

import (
	"gopkg.in/mgo.v2"
	"time"
)

func ConnectDB() mgo.Session {
	MongoDBHosts := "127.0.0.1"
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Username: "admin",
		Password: "password",
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
