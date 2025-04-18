package model

import (
	"sort"
	"strings"

	"xorm.io/xorm"
)

type Permission struct {
	ID          int64  `json:"id" xorm:"id pk autoincr notnull unique index"`
	Name        string `json:"name" xorm:"name"`
	Code        string `json:"code" xorm:"code"`
	Type        string `json:"type" xorm:"type null "`
	Parentid    int64  `json:"parent_id" xorm:"parentId null "`
	Path        string `json:"path" xorm:"path null "`
	Redirect    string `json:"redirect" xorm:"redirect null "`
	Icon        string `json:"icon" xorm:"icon null "`
	Component   string `json:"component" xorm:"component null "`
	Layout      string `json:"layout" xorm:"layout null "`
	Keepalive   bool   `json:"keep_alive" xorm:"keepAlive null tinyint(1)"`
	Method      string `json:"method" xorm:"method null "`
	Description string `json:"description" xorm:"description null "`
	Show        bool   `json:"show" xorm:"show notnull default(1) tinyint(1)"` // 是否展示在页面菜单
	Enable      bool   `json:"enable" xorm:"enable notnull default(1) tinyint(1)"`
	Order       int64  `json:"order" xorm:"order"`
}

// TableName 表名称
func (*Permission) TableName() string {
	return "permission"
}

type PermissionModel struct {
	db *xorm.Engine
	Permission
}

// MenuNode 菜单节点结构，用于构建树形结构
type MenuNode struct {
	Permission
	Children []*MenuNode `json:"children"`
}

func NewPermission(db *xorm.Engine) *PermissionModel {
	return &PermissionModel{db: db}
}
func (p *PermissionModel) GetAllPermissions() ([]Permission, error) {
	var permissions []Permission
	err := p.db.Find(&permissions)
	return permissions, err
}

// GetMenuTree 获取菜单树
func (p *PermissionModel) GetMenuTree() ([]*MenuNode, error) {
	// 获取所有权限
	permissions, err := p.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return p.buildMenuTree(permissions, 0), nil
}

// GetMenuTreeByType 根据类型获取菜单树
func (p *PermissionModel) GetMenuTreeByType(menuType string, parentId int) ([]Permission, error) {
	var permissions []Permission
	err := p.db.Where("type = ? AND enable = 1 and parentId=?", strings.ToUpper(menuType), parentId).Find(&permissions)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return permissions, nil
}

// buildMenuTree 构建菜单树
func (p *PermissionModel) buildMenuTree(permissions []Permission, parentID int64) []*MenuNode {
	var nodes []*MenuNode

	// 找出当前层级的所有节点
	for _, perm := range permissions {
		if perm.Parentid == parentID {
			node := &MenuNode{
				Permission: perm,
				Children:   p.buildMenuTree(permissions, perm.ID), // 递归构建子节点
			}
			nodes = append(nodes, node)
		}
	}

	// 根据 Order 字段排序
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Order < nodes[j].Order
	})

	return nodes
}

// GetMenuByID 根据ID获取菜单及其所有子菜单
func (p *PermissionModel) GetMenuByID(id int64) (*MenuNode, error) {
	// 获取指定ID的菜单
	var menu Permission
	has, err := p.db.ID(id).Get(&menu)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	// 获取所有权限
	permissions, err := p.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	// 构建当前菜单的节点
	node := &MenuNode{
		Permission: menu,
		Children:   p.buildMenuTree(permissions, menu.ID),
	}

	return node, nil
}

// AddMenu 添加菜单
func (p *PermissionModel) AddMenu(menu *Permission) (int64, error) {
	return p.db.Insert(menu)
}

// UpdateMenu 更新菜单
func (p *PermissionModel) UpdateMenu(menu *Permission) (int64, error) {
	return p.db.ID(menu.ID).Update(menu)
}

// DeleteMenu 删除菜单及其子菜单
func (p *PermissionModel) DeleteMenu(id int64) error {
	// 开启事务
	session := p.db.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	// 获取所有子菜单ID
	childIDs, err := p.getChildMenuIDs(id)
	if err != nil {
		session.Rollback()
		return err
	}

	// 删除所有子菜单
	if len(childIDs) > 0 {
		_, err = session.In("id", childIDs).Delete(&Permission{})
		if err != nil {
			session.Rollback()
			return err
		}
	}

	// 删除当前菜单
	_, err = session.ID(id).Delete(&Permission{})
	if err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}

// getChildMenuIDs 递归获取所有子菜单ID
func (p *PermissionModel) getChildMenuIDs(parentID int64) ([]int64, error) {
	var permissions []Permission
	err := p.db.Where("parent_id = ?", parentID).Find(&permissions)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, perm := range permissions {
		ids = append(ids, perm.ID)
		// 递归获取子菜单ID
		childIDs, err := p.getChildMenuIDs(perm.ID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, childIDs...)
	}

	return ids, nil
}

// 根据组ID获取权限列表
func (p *PermissionModel) GetPermissionsByRoleID(roleID int64) ([]*MenuNode, error) {
	var permissions []Permission

	db := p.db.Table("permission").Where("permission.type=?", "MENU")
	if roleID != -1 {
		db = db.Join("INNER", "rule", "permission.id = rule.v1")
		db = db.Where("rule.v0=?", roleID).Where("rule.ptype='p' and rule.v2 = 'role'")
	}
	err := db.Find(&permissions)
	return p.buildMenuTree(permissions, 0), err
}

// 根据组ID获取权限列表
func (p *PermissionModel) GetPermissionsTreeAll(roleID int64) ([]*MenuNode, error) {
	var permissions []Permission

	db := p.db.Table("permission")
	if roleID != -1 {
		db = db.Join("INNER", "rule", "permission.id = rule.v1")
		db = db.Where("rule.v0=?", roleID).Where("rule.ptype='p' and rule.v2 = 'role'")
	}
	err := db.Find(&permissions)
	return p.buildMenuTree(permissions, 0), err
}

// 根据ID和type=API
func (p *PermissionModel) GetPermissionsByRoleIDAndType(roleIDs []int64, menuType string) ([]Permission, error) {
	var permissions []Permission
	p.db.Where("type =? AND enable = 1", strings.ToUpper(menuType)).In("id", roleIDs).Find(&permissions)
	return permissions, nil
}
