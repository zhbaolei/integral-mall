package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/log4g"
	"integral-mall/user/command/api/config"
	"integral-mall/user/controller"
	"integral-mall/user/logic"
	"integral-mall/user/model"
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

	userLogic :=logic.NewUserLogic(userModel,RedisClient)

	userController:=controller.NewUserController(userLogic)
	r := gin.Default()

	userRouterGroup :=r.Group("/user")
	{
		userRouterGroup.POST("/register",userController.Register)
		userRouterGroup.POST("/login",userController.Login)
	}

	/*r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})*/
	log4g.Error(r.Run(conf.Port)) // listen and serve on 0.0.0.0:8080
}