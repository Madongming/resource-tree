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
	}
	for i := range tables {
		err := DB().Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(tables[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
