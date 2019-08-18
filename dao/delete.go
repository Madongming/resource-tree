package dao

import (
	. "global"
	"model"
)

func PreDeleteNode(nodeId interface{}) (*model.Tree, error) {
	return GetGraphByNodeId(nodeId)
}

func DeleteNode(nodeId interface{}) error {
	if err := deleteResourceRelationshipByNodeId(nodeId); err != nil {
		return err
	}
	return deleteNodeById(nodeId)
}

func DeleteResourceRelationship(srcNodeId, tarNodeId interface{}) error {

}
