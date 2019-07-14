package cache

import (
	"model"
)

func init() {
	if Cache == nil {
		fetchCache()
	}
}

func CasResource() error {
	version, err := model.GetResourceVersion()
	if err != nil {
		return err
	} else if Cache != nil && Cache.Version == version {
		return nil
	}
	if err := fetchCache(version); err != nil {
		return err
	}
	return nil
}

func FindTree(nodeId int) (*ResourceTree, error) {
	if tree, found := Cache.Index[nodeId]; !found {
		return nil, nil
	}

	return tree, nil
}

func fetchCache(opts ...int64) error {
	var err error
	var version int64
	if opts == nil || len(opts) == 0 {
		version, err = model.GetResourceVersion()
		if err != nil {
			return err
		}
	} else {
		version = opts[0]
	}
	allNodes, err := model.GetAllNodes()
	if Cache, err = makeNodeArray2Tree(allNodes, version); err != nil {
		return err
	}
	return nil
}
