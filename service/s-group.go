package service

import (
	"cash-server/db"
	"cash-server/model"
	"fmt"
	"regexp"
)

//GroupOrderGet 查詢群組order  下筆 subid
func GroupOrderGet(groupid int, StageType int) string {
	GroupOrderUpdate(groupid)
	fid := model.GroupExist(groupid).GroupName
	subid := fmt.Sprintf("%0*d", 9, model.GroupExist(groupid).GroupOrder)
	if StageType == 0 {
		return fid + "x" + subid
	}
	return fid + "_" + subid
}

//GroupOrderUpdate 更新群組order  subid
func GroupOrderUpdate(groupid int) {
	model.GroupOrderUpdate(groupid)
}

//GroupAdd 群組新增
func GroupAdd(g db.PlatformGroup) {
	if match, _ := regexp.MatchString("([a-zA-Z0-9]+)", g.GroupName); match && model.GroupNameCheck(g.GroupName) {
		model.GroupAdd(g)
	}

}

//GroupAuthAdd 群組權限新增
func GroupAuthAdd(ga db.PlatformGroupsAuth) {
	model.GroupAuthAdd(ga)
}
