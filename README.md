# go-admin
基于 gin + gorm 实现的后台管理框架，功能还在初始阶段；

会持续更新，逐步实现 RBAC 权限管理，最终结合前端 vue 实现类似我的 python 项目的堡垒机: https://github.com/leffss/devops

# 运行
首先根据实际情况修改配置文件 `conf/app.ini`

生成 swag 最新 api 文档
```
go get -u github.com/swaggo/swag/cmd/swag
swag init
```
- 项目根目录执行

初始化数据库
```
go run main.go init -p 123456
```
- -p 指定 admin 超级管理员密码

运行
```
go run main.go
```

或者编译后运行
```
go build -o go-admin main.go
chmod 755 go-admin
./go-admin
```

查看 swag api 文档

http://127.0.0.1:8080/api/swagger/index.html

# MIT License
```
Copyright (c) 2021-2022 leffss.
```
