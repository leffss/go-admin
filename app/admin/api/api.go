package api

import (
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/app/auth/service"
	"github.com/leffss/go-admin/models"
	"github.com/leffss/go-admin/pkg/app"
	"github.com/leffss/go-admin/pkg/codemsg"
	"github.com/leffss/go-admin/pkg/jwt"
	"github.com/leffss/go-admin/tasks"
	"github.com/unknwon/com"
	"net/http"
)

// Test
// @Summary 测试
// @Description 测试
// @Tags admin
// @Produce json
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Router /api/admin/test [get]
func Test(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	token := c.GetString("token")
	claims, _ := jwt.ParseToken(token)
	user := service.GetUserByUsername(claims.UserName)
	user.GetPermissions()

	//perms := models.GetPerm(3, models.GetAllPermission())
	perms := models.GetPerm(3, models.GetPermissionByUser(user))

	appG.Response(http.StatusOK, codemsg.Success, gin.H{
		"admin": user,
		"perms": perms,
	})
}

// TestTask
// @Summary 添加一个测试后台任务
// @Description 添加一个测试后台任务
// @Tags admin
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/admin/test/{id} [get]
func TestTask(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := com.StrTo(c.Param("id")).MustInt()

	res, _ := tasks.SendTask("adminTaskTest", id)

	appG.Response(http.StatusOK, codemsg.Success, gin.H{
		"admin": res,
	})
}
