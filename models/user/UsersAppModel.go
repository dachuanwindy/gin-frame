package user

import (
	"database/sql"
	"gin-frame/models/base"
	"github.com/gin-gonic/gin"
	"gin-frame/libraries/util"
)

func GetInfoByAppUid(c *gin.Context, AppUid int) map[string]interface{} {
	db := base.GetConn("ymt360")

	var data = make(map[string]interface{})
	var (
		appUid 		int
		appId		int
		lastToken	string
		securekey	string
		cid			int
		mobile		string
		status 		int
		fCode		int
	)

	query := "select app_uid, app_id, last_token, secure_key, customer_id, mobile, status, fcode from users_app where app_uid = ? limit 1"
	row := db.MasterDBQueryRowContext(c.Request.Context(), query, AppUid)
	err := row.Scan( &appUid, &appId, &lastToken, &securekey, &cid, &mobile, &status, &fCode)
	if err == sql.ErrNoRows {
		return data
	}
	util.Must(err)

	data["app_uid"] = appUid
	data["app_id"] = appId
	data["county_id"] = lastToken
	data["last_token"] = securekey
	data["secure_key"] = cid
	data["mobile"] = mobile
	data["status"] = status
	data["fcode"] = fCode

	return data
}
