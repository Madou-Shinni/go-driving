package service

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/global"
	sysService "github.com/flipped-aurora/gin-vue-admin/server/service"
	"gorm.io/gorm"
	"strconv"
)

var (
	ErrorMenuExits          = errors.New("存在重复name，请修改name")
	ErrorMenuAuthorityExits = errors.New("已拥有菜单权限")
)

var (
	systemService = sysService.ServiceGroupApp.SystemServiceGroup
)

// AddBaseMenu 添加菜单
func AddBaseMenu(menu system.SysBaseMenu) (id uint, err error) {
	if !errors.Is(global.GlobalConfig.DB.Where("name = ?", menu.Name).First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return 0, ErrorMenuExits
	}

	err = global.GlobalConfig.DB.Create(&menu).Error
	if err != nil {
		return
	}

	return menu.ID, nil
}

// SetMenuAuthority 设置菜单权限
func SetMenuAuthority(menuid uint, authorityid uint) error {
	if !errors.Is(global.GlobalConfig.DB.Where("sys_base_menu_id = ? AND sys_authority_authority_id =?", menuid, authorityid).First(&system.SysAuthorityMenu{}).Error, gorm.ErrRecordNotFound) {
		return ErrorMenuAuthorityExits
	}

	return global.GlobalConfig.DB.Create(&system.SysAuthorityMenu{MenuId: strconv.FormatUint(uint64(menuid), 10), AuthorityId: strconv.FormatUint(uint64(authorityid), 10)}).Error
}

// AddApiAuthority 添加并设置Api权限
func AddApiAuthority(authorityId uint, apis []system.SysApi) error {
	// 添加api权限
	db := global.GlobalConfig.DB
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, api := range apis {
			if !errors.Is(tx.Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
				// 已存在api
				continue
			}

			err := tx.Create(&api).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 获取所有的PolicyPath
	maps := systemService.GetPolicyPathByAuthorityId(authorityId)
	// 添加api
	for i := range apis {
		maps = append(maps, request.CasbinInfo{Path: apis[i].Path, Method: apis[i].Method})
	}
	// 更新Casbin
	err = systemService.UpdateCasbin(authorityId, maps)
	if err != nil {
		return err
	}

	return nil
}
