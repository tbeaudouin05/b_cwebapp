package useraccess

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func AddUserLogin(email string, name string, mongoSession *mgo.Session) {

	db := mongoSession.DB("competition_analysis")
	collection := db.C("user")
	//queries
	_, err := collection.Upsert(bson.M{"email": email}, bson.M{"$set": bson.M{"name": name, "access": "normal"}})
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
