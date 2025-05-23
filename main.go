package main

import (
	"github.com/deatil/lakego-admin/bootstrap"
)

//go:generate swag init -o=./swagger --parseDependency --parseDepth=1

// @title lakego-admin 后端 API 文档

// @version 1.3.1
// @description lakego-admin 是基于 gin、JWT 和 RBAC 的 go 后台管理系统
// @termsOfService https://github.com/deatil

// @license.name Apache2
// @license.url https://github.com/deatil/lakego-admin/blob/main/LICENSE

// @contact.name deatil
// @contact.url https://github.com/deatil
// @contact.email deatil@github.com

// @host 127.0.0.1:8080
// @BasePath /admin-api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	bootstrap.Execute()
}
