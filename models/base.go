package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id uint `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
}

type BaseModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type BaseModelNoUpdateTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type ControlBy struct {
	CreateBy uint `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy uint `json:"updateBy" gorm:"index;comment:更新者"`
}

// SetCreateBy 设置创建人id
func (e *ControlBy) SetCreateBy(createBy uint) {
	e.CreateBy = createBy
}

// SetUpdateBy 设置修改人id
func (e *ControlBy) SetUpdateBy(updateBy uint) {
	e.UpdateBy = updateBy
}
