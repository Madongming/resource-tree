package cache

import (
	"model"
)

func init() {
	if Cache == nil {
		Cache, _ = fetchCache()
	}
}

func casResource() (err error) {
	version, err := model.GetResourceVersion()
	if err != nil {
		return err
	} else if Cache.Version == version {
		return nil
	}
	if Cache, err = model.GetResource(); err != nil {
		return err
	}
	return nil
}
