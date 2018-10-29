package homepage

import (
	"fmt"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/user"
)

func FetchPersonalizedInfo(email string, name string, mongoSession *mgo.Session) user.UserInfo {

	db := mongoSession.DB("competition_analysis")
	matchingReport := db.C("user")
	var info []user.UserInfo
	matchingReport.Find(bson.M{"email": email}).All(&info)
	fmt.Println("info:", info)
	if info != nil {
		matchedNumStr := strconv.Itoa(info[0].ThisWeekSKU)
		info[0].Message = "You matched " + matchedNumStr + " SKUs in the last 7 days. "
		if info[0].ThisWeekSKU >= 30 {
			info[0].Message = info[0].Message + "Great job, thanks! :)"
		} else if info[0].ThisWeekSKU > 0 {
			info[0].Message = "You matched only " + matchedNumStr + " SKUs in the last 7 days. "
		} else {
			info[0].Message = "You didn't match any SKUs in the last 7 days. "
		}
		return info[0]
	}
	// happens when user doesn't have any information
	info0 := user.UserInfo{
		Name:                name,
		Email:               email,
		TotalSKU:            0,
		TotalSupplier:       0,
		ThisWeekSKU:         0,
		LastWeekSKU:         0,
		TotalNotCompetitive: 0,
		Message:             "There is no historical information about you, wait 24 hours after your first matching and if it continues contact us.",
	}

	return info0

}
