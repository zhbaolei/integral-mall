package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/grpcx"
	"github.com/yakaa/log4g"
	"google.golang.org/grpc"
	"integral-mall/user/command/rpc/config"
	"integral-mall/user/logic"
	"integral-mall/user/model"
	"integral-mall/user/protos"
	"io/ioutil"
	"log"
)
var configFile=flag.String("f","config/config.json","user config")

func main() {
	flag.Parse()
	conf:=new(config.Config)
	bs,err:=ioutil.ReadFile(*configFile)
	if err!=nil {
		log.Fatal(err)
	}

	if err :=json.Unmarshal(bs,conf);err !=nil{
		log.Fatal(err)
	}

	if conf.Mode ==gin.DebugMode {
		log4g.Init(log4g.Config{Path:"logs"})
		gin.DefaultWriter=log4g.InfoLog
		gin.DefaultErrorWriter=log4g.ErrorLog
	}
	if err != nil {
		log.Fatal(err)
	}

	engine, err := xorm.NewEngine("mysql", conf.Mysql.DataSource)
	if err != nil {
		log.Fatal(err)
	}

	RedisClient := redis.NewClient(&redis.Options{Addr:conf.Redis.DataSource,Password:conf.Redis.Auth })

	userModel:=model.NewUserModel(engine,RedisClient,conf.Mysql.Table.User)

	userServerLogic := logic.NewUserRpcServerLogic(userModel)

	rpcServer, err := grpcx.MustNewGrpcxServer(conf.RpcServerConfig, func(server *grpc.Server) {
		protos.RegisterUserRpcServer(server, userServerLogic)
	})

	if err != nil {
		log.Fatal(err)
	}

	log4g.Error(rpcServer.Run())
}