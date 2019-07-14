package model

import (
	log "github.com/cihub/seelog"

	"global"
)

func getCurrentVersion() (int64, error) {
	var version int64
Retry:
	if err := DB().
		Raw("SELECT current FROM version WHERE id = 1").
		Scan(&version).
		Error; err != nil {
		return int64(0), err
	}
	return version, nil
}

func casVersion() error {
	// select current as c from version where id = 1;
	// update current set current = current + 1 where id = 1 and current = c;

Retry:
	version, err := getCurrentVersion()
	if err != nil {
		return err
	}

	if DB().
		Exec("UPDATE version "+
			"set current = current + 1 "+
			"WHERE id = 1 and current = ?",
			version).RowsAffected < int64(1) {
		goto Retry
	}
	return nil
}
