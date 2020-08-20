package casinogrpc

import (
	"cash-server/db"
	"cash-server/pkg/util"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

//CasinoItem 給前端Item的格式
type CasinoItem struct {
	Platform  string
	ProductID string
	ItemType  int32
	ItemName  string
	Price     int32
}

//GrpcCasino Grpc 控制
var GrpcCasino *grpc.ClientConn

//GrpcCasinoCannot 連線
func GrpcCasinoCannot() grpc.ClientConnInterface {
	addr := "34.92.237.98:30001"
	//TODO
	// 正式環境
	// if configs.GetGlobalConfig().RunMode == "release" {
	// 	addr = configs.GetGlobalConfig().Casino.Proip
	// }
	GrpcCasino, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		util.Error("[GRPC-Casino] Can not connect to gRPC server: %v", err)
	}
	return GrpcCasino
}

//VetifyUserGUID 驗證玩家ID GUID
func VetifyUserGUID(ID string) (msg string) {
	util.Test(fmt.Sprint("[GRPC-Casino] 驗證玩家ID GUID >  ", ID))
	var newids Ids
	newids.MyID = ID
	conn := GrpcCasinoCannot()
	c := NewMemberClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a, err := c.GetAccountByID(ctx, &newids)
	if err != nil {
		util.Error("[GRPC-Casino] Could not get nonce: %v", err)
	}
	if a.Message != "" {
		util.Test("[GRPC-Casino] 驗證玩家ID資料 > No Account ")
		return "Error"
	}
	util.Test(fmt.Sprint("[GRPC-Casino] 驗證玩家ID GUID >  ", a.Account.GUID))
	return a.Account.GUID
}

//VetifyUserName 驗證玩家ID名稱
func VetifyUserName(ID string) (msg string) {
	util.Test(fmt.Sprint("[GRPC-Casino] 驗證玩家ID 名稱 >  ", ID))
	var newids Ids
	newids.MyID = ID
	conn := GrpcCasinoCannot()
	c := NewMemberClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a, err := c.GetAccountByID(ctx, &newids)
	if err != nil {
		util.Error("[GRPC-Casino] Could not get nonce: %v", err)
	}
	if a.Message != "" {
		util.Test("[GRPC-Casino] 驗證玩家ID名稱 > No Account ")
		return "Error"
	}
	util.Test(fmt.Sprint("[GRPC-Casino] 驗證玩家ID 名稱 >  ", a.Account.Name))
	return a.Account.Name
}

//VetifyItem 檢查商店資料驗證
func VetifyItem(itemID string, ItemPrice int32) bool {
	util.Test(fmt.Sprint("[GRPC-Casino] 檢查商店資料驗證 >  ID :", itemID, " price : ", ItemPrice))
	var itemList DbStrInput
	var newitems DbIAPItem
	newitems.ProductID = itemID
	newitems.Price = ItemPrice
	conn := GrpcCasinoCannot()
	c := NewMemberClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a, err := c.DbGetIAPItem(ctx, &itemList)
	if err != nil {
		util.Error("[GRPC-Casino] Could not get nonce: %v", err)
	}

	for _, itemobj := range filterItem(a.GetList()) {
		//util.Test(fmt.Sprint(itemobj))
		if itemobj.GetProductID() == itemID && itemobj.GetPrice() == ItemPrice {
			util.Test(fmt.Sprint("[GRPC-Casino] 檢查商店資料驗證 >  成功"))
			return true
		}
	}

	util.Test(fmt.Sprint("[GRPC-Casino] 檢查商店資料驗證 >  錯誤"))
	return false
}

//GetItem 查詢商店清單
func GetItem() []CasinoItem {
	util.Test("[GRPC-Casino] 查詢商店清單")
	var itemList DbStrInput
	var casinoItem []CasinoItem
	conn := GrpcCasinoCannot()
	c := NewMemberClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a, err := c.DbGetIAPItem(ctx, &itemList)
	if err != nil {
		util.Error("[GRPC-Casino] Could not get nonce: %v", err)
	}
	itemListJSON := db.Struct2JSON(removeDuplicateElement(filterItem(a.GetList())))
	util.Test(fmt.Sprint("[GRPC-Casino] 查詢商店清單 > ", itemListJSON))
	json.Unmarshal([]byte(itemListJSON), &casinoItem)
	util.Test(fmt.Sprint("[GRPC-Casino] 查詢商店清單 (縮小) > ", db.Struct2JSON((casinoItem))))
	return casinoItem
}

//filterItem ([]CasinoItem) 篩選出Mycard的資料
func filterItem(orig []*DbIAPItem) []*DbIAPItem {
	new := orig[:0]
	for _, xCasinoItem := range orig {
		if xCasinoItem.Platform == "mycard" {
			new = append(new, xCasinoItem)
		}
	}
	return new
}

//filterItemTest ([]CasinoItem)  去除slice中重复的元素
func removeDuplicateElement(orig []*DbIAPItem) []*DbIAPItem {
	result := make([]*DbIAPItem, 0, len(orig))
	temp := map[*DbIAPItem]struct{}{}
	for _, item := range orig {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

//SendItemBuy 購買 - 送出購買要求
func SendItemBuy(userid string, prodid string) (bool, int32) {
	util.Test("[GRPC-Casino] 送出購買要求")
	var dbIAPBuyInput DbIAPBuyInput
	dbIAPBuyInput.Platform = "mycard"
	dbIAPBuyInput.Quantity = 1
	util.Test(fmt.Sprint("[GRPC-Casino] Userid : ", userid))

	dbIAPBuyInput.GUID = userid
	dbIAPBuyInput.ProductID = prodid
	conn := GrpcCasinoCannot()
	c := NewMemberClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a, err := c.DbIAPBuy(ctx, &dbIAPBuyInput)
	if err != nil {
		util.Error("[GRPC-Casino] Could not get nonce: %v", err)
	}
	util.Test(fmt.Sprint("[GRPC-Casino] Game 資料庫建單 :", a.Success, " / Game 資料庫單號 :", a.Value))
	return a.GetSuccess(), a.GetValue()
}

//SendItemResult 購買 - 送出購買要求(回傳訂單序號)
func SendItemResult(userid string, prodid string) bool {
	util.Test("[GRPC-Casino] 回傳訂單序號")
	var cq DbIAPBuyResultInput
	cq.Platform = "mycard"
	cq.Quantity = 1
	cq.GUID = userid
	cq.ProductID = prodid
	cq.Status = 1
	cq.Content = "11122233"
	conn := GrpcCasinoCannot()
	c := NewMemberClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a, err := c.DbIAPBuyResult(ctx, &cq)
	if err != nil {
		util.Error("[GRPC-Casino] Could not get nonce: %v", err)
	}
	util.Test(fmt.Sprint(a))
	return true
}
