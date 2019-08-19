package dao

import (
	. "github.com/Madongming/resource-tree/global"
	"github.com/Madongming/resource-tree/model"
)

func PreDeleteNode(nodeId interface{}) (*model.Graph, error) {
	return GetNodeGraph(nodeId)
}

func DeleteNode(nodeId interface{}) error {
	if err := deleteResourceRelationshipByNodeId(nodeId); err != nil {
		return err
	}
	casEdgeVersion()
	return deleteNodeById(nodeId)
}

func DeleteResourceRelationship(srcNodeId, tarNodeId interface{}) error {
	rr := new(model.DBResourceRelationship)
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
	return rr.Delete()
}
