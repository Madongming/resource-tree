package model

import (
	. "github.com/Madongming/resource-tree/global"
)

func makeTables() error {
	tables := []interface{}{
		&DBResourceNode{},
		&DBUser{},
		&DBGroup{},
		&DBUserGroup{},
		&DBUserPermission{},
		&DBGroupPermission{},
		&DBResourceRelationship{},
		&NodeVersion{},
		&EdgeVersion{},
	}
	err := DB().Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(tables...).Error
	if err != nil {
		return err
	}
	if err := DB().Create(&NodeVersion{Current: 1}).Error; err != nil {
		return err
	}
	if err := DB().Create(&EdgeVersion{Current: 1}).Error; err != nil {
		return err
	}
	return nil
}
