package useraccess

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"
)

func GetUserAccess(user *useraccess.User, mongoSession *mgo.Session) {

	db := mongoSession.DB("competition_analysis")
	collection := db.C("user")

	collection.Find(bson.M{"email": user.Email}).One(&user)

}
