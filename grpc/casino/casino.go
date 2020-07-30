package casinogrpc

import (
	"cash-server/pkg/util"
	"context"
	"time"

	"google.golang.org/grpc"
)

//GrpcCasino Grpc 控制
var GrpcCasino *grpc.ClientConn

//GrpcCasinoCannot 連線
func GrpcCasinoCannot(addr string) grpc.ClientConnInterface {
	GrpcCasino, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		//log.Fatalf("Can not connect to gRPC server: %v", err)
		util.Error("Can not connect to gRPC server: %v", err)
	}
	return GrpcCasino
}

//VetifyUserID 驗證玩家ID資料
func VetifyUserID(ID string) (msg string) {
	var newids Ids
	newids.MyID = ID
	addr := "35.194.245.46:30001"
	conn := GrpcCasinoCannot(addr)
	c := NewMemberClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a, err := c.GetAccountByID(ctx, &newids)
	if err != nil {
		//log.Fatalf("Could not get nonce: %v", err)
		util.Error("Could not get nonce: %v", err)
	}
	//fmt.Println(a)
	if a.Message != "" {
		return "NoAccount"
	}

	return a.Account.GUID
}
