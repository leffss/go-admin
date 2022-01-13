package main

import (
	"github.com/leffss/go-admin/cmd"
)

// @title go-admin
// @version 1.0.0
// @description api 文档
// @termsOfService https://github.com/leffss/go-admin
// @license.name MIT
// @license.url https://github.com/leffss/go-admin/blob/master/LICENSE

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
