package permission

import "github.com/dysodeng/app/internal/infrastructure/shared/model"

// Admin 管理员
type Admin struct {
	model.PrimaryKeyID
	Username     string `gorm:"index:admin_username_idx,unique;type:varchar(50);not null;default:'';comment:用户名" json:"username"`
	SafePassword string `gorm:"type:varchar(150);not null;default:;'';comment:登录密码" json:"safe_password"`
	RealName     string `gorm:"type:varchar(50);not null;default:'';comment:姓名" json:"real_name"`
	Telephone    string `gorm:"type:varchar(20);not null;default:'';comment:手机号" json:"telephone"`
	Remark       string `gorm:"type:varchar(50);not null;default:'';comment:备注" json:"remark"`
	IsSuper      uint8  `gorm:"not null;default:0;comment:是否超级管理员 0-否 1-是" json:"is_super"`
	Status       uint8  `gorm:"not null;default:0;comment:状态 0-禁用 1-启用" json:"status"`
	model.Time
}

func (Admin) TableName() string {
	return "ams_admin"
}

// Permission 权限节点
type Permission struct {
	model.PrimaryKeyID
	Identify string `gorm:"index:permission_idx,unique;type:varchar(100);not null;default:'';comment:权限唯一标识" json:"identify"`
	Name     string `gorm:"type:varchar(100);not null;default:'';comment:权限名称" json:"name"`
	ParentID uint64 `gorm:"not null;default:0;comment:权限父级节点ID" json:"parent_id"`
	Sort     uint   `gorm:"not null;default:0;comment:排序值,越小越靠前" json:"sort"`
	model.Time
}

func (Permission) TableName() string {
	return "ams_admin_permissions"
}

// AdminHasPermission 管理员拥有的权限
type AdminHasPermission struct {
	model.PrimaryKeyID
	AdminID      uint64 `gorm:"index:admin_has_perm_idx,unique;not null;default:0;comment:管理员ID" json:"admin_id"`
	PermissionID uint64 `gorm:"index:admin_has_perm_idx,unique;not null;default:0;comment:权限节点ID" json:"permission_id"`
	model.Time
}

func (AdminHasPermission) TableName() string {
	return "ams_admin_has_permissions"
}
