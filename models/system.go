package models

import (
	"errors"
)

type Department struct {
	BaseModel
	Name string `json:"name" gorm:"unique;not null;type:varchar(32);comment:部门名称"`
	ParentId uint `json:"parentId" gorm:"comment:上级部门"`
	ControlBy
	BaseModelTime
}

func (Department) TableName() string {
	return "department"
}

type Permission struct {
	BaseModel
	Name string `json:"name" gorm:"unique;not null;type:varchar(32);comment:'权限名称'"`
	// 前端验证用
	Sign string `json:"sign" gorm:"unique;not null;type:varchar(32);comment:'权限标识'"`
	Menu bool `json:"menu" gorm:"not null;type:bool;comment:'是否为左侧菜单'"`
	Method string `json:"method" gorm:"not null;type:varchar(8);comment:'HTTP请求方法'"`
	Path string `json:"path" gorm:"not null;type:varchar(256);comment:'HTTP请求路径正则'"`
	ParentId uint `json:"parentId" gorm:"comment:上级权限"`
	BaseModelTime
	Roles []Role `json:"roles" gorm:"many2many:role_permission"`
}

func (Permission) TableName() string {
	return "permission"
}

func GetAllPermission() *[]Permission {
	var perms []Permission
	DB.Find(&perms)
	return &perms
}

func GetPermissionByUser(user *User) *[]Permission {
	var perms []Permission
	if user.IsSuper {
		return GetAllPermission()
	}
	if user.RoleId == 0 {
		return nil
	}
	var role Role
	DB.Preload("Permissions").First(&role, user.RoleId)
	for _, v := range role.Permissions {
		perms = append(perms, v)
	}
	return &perms
}

type Tree struct {
	Id uint `json:"id"`
	Label string `json:"label"`
	Children []*Tree `json:"children"`
}

func GetPerm(from uint, permissions *[]Permission) (perms []Tree) {
	var child = Tree{0, "所有权限", []*Tree{}}
	err := getTreeNode(from, &child, permissions)
	if err != nil {
		return
	}
	perms = append(perms, child)
	return perms
}

func getTreeNode(pid uint, treeNode *Tree, permissions *[]Permission) error {
	var perms []Permission
	//DB.Where("parent_id = ?", pid).Find(&perms)
	for _, v := range *permissions {
		if v.ParentId == pid {
			perms = append(perms, v)
		}
	}
	if len(perms) == 0 {
		return errors.New("not record")
	}
	for i := 0; i < len(perms); i++ {
		child := Tree{perms[i].Id, perms[i].Name, []*Tree{}}
		treeNode.Children = append(treeNode.Children, &child)
		_ = getTreeNode(perms[i].Id, &child, permissions)
	}
	return nil
}

type Role struct {
	BaseModel
	Name string `json:"name" gorm:"unique;not null;type:varchar(32);comment:角色名称"`
	DeptId uint `json:"deptId" gorm:"comment:部门ID"`
	ControlBy
	BaseModelTime
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permission"`
}

func (Role) TableName() string {
	return "role"
}
