package models

import (
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Username string `json:"username" gorm:"unique;not null;type:varchar(32);comment:用户名"`
	Name string `json:"name" gorm:"type:varchar(32);comment:姓名"`
	Avatar string `json:"avatar" gorm:"type:varchar(255);comment:头像"`
	// sha256 加密，固定64位
	Password string `json:"-" gorm:"not null;type:char(64);comment:密码"`
	IsSuper bool `json:"isSuper" gorm:"not null;type:bool;comment:是否为超级管理员"`
	Status bool `json:"status" gorm:"not null;type:bool;comment:是否启用"`
	Mobile string `json:"mobile" gorm:"type:char(11);comment:'手机'"`
	Email string `json:"email" gorm:"type:varchar(255);comment:'邮箱'"`
	// 未使用 gorm 自带的一对多实现
	// 采用了在一对多的两个实体间，在“多”的实体表(即用户)中新增一个字段，该字段是“一”实体表(即部门)的主键。
	DeptId uint `json:"deptId" gorm:"comment:部门ID"`
	RoleId uint `json:"roleId" gorm:"comment:角色ID"`
	DeptIds []uint `json:"-" gorm:"-"`
	RoleIds []uint `json:"-" gorm:"-"`
	// 权限标识
	Permissions []string `json:"permissions" gorm:"-"`
	ControlBy
	BaseModelTime
}

func (User) TableName() string {
	// 默认表名是 users
	return "user"
}

func (u *User) AfterFind(_ *gorm.DB) error {
	u.DeptIds = []uint{u.DeptId}
	u.RoleIds = []uint{u.RoleId}
	return nil
}

func (u *User) GetPermissions() {
	if u.IsSuper {
		u.Permissions = []string{"admin"}
		return
	}
	if u.RoleId == 0 {
		return
	}
	var role Role
	DB.Preload("Permissions").First(&role, u.RoleId)
	for _, v := range role.Permissions {
		u.Permissions = append(u.Permissions, v.Sign)
	}
}

type LoginLog struct {
	BaseModel
	// IPV6最大长度 39 位
	Ip string `json:"ip" gorm:"type:varchar(40);comment:IP地址"`
	Location string `json:"location" gorm:"type:varchar(64);comment:地点"`
	Client string `json:"client" gorm:"type:varchar(255);comment:客户端"`
	Os string `json:"os" gorm:"type:varchar(255);comment:操作系统"`
	Status bool `json:"status" gorm:"type:bool;comment:成功或者失败"`
	Kind string `json:"kind" gorm:"type:varchar(32);comment:日志类型"`
	Detail string `json:"detail" gorm:"type:varchar(255);comment:操作详情"`
	DeptId uint `json:"deptId" gorm:"comment:部门ID"`
	ControlBy
	BaseModelNoUpdateTime
}

func (LoginLog) TableName() string {
	return "login_log"
}
