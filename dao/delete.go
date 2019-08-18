package dao

import (
	"cache"
	. "global"
	"model"
)

func PreDeleteNode(nodeId interface{}) (*cache.CacheGraph, error) {
	return GetGraphByNodeId(nodeId)
}

func DeleteNode(nodeId interface{}) error {
	if err := deleteResourceRelationshipByNodeId(nodeId); err != nil {
		return err
	}
	casEdgeVersion()
	return deleteNodeById(nodeId)
}

func DeleteResourceRelationship(srcNodeId, tarNodeId interface{}) error {
	rr := new(DBResourceRelationship)
	if err := DB().
		Where(
			"source_resource_node_id = ? AND "+
				"target_resource_node_id = ?",
			srcNodeId,
			tarNodeId).
		First(rr).
		Error; err != nil {
		return err
	}
	casEdgeVersion()
	return rr.delete()
}
