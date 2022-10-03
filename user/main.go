package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	"user/domain/repository"
	service2 "user/domain/service"
	"user/handler"
	user "user/proto"
)

func main() {
	//服务参数设置
	srv := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("latest"),
	)
	//初始化服务
	srv.Init()

	//创建数据库连接
	// dsn := "root:123456@tcp(127.0.0.1:8000)/micro?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/micro?charset=utf8mb4&parseTime=True")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.SingularTable(true)

	//只执行一次，数据表初始化
	rp := repository.NewUserRepository(db)
	rp.InitTable()

	//创建服务实例
	userDataService := service2.NewUserDataService(repository.NewUserRepository(db))

	//注册Handler
	err = user.RegisterUserHandler(srv.Server(), &handler.User{UserDataService: userDataService})
	if err != nil {
		fmt.Println(err)
	}

	//Run service
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
