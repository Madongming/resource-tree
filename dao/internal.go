package dao

import (
	"fmt"
	"strings"

	"github.com/Madongming/wheelmaking/data-structure/queue"
	log "github.com/cihub/seelog"

	"github.com/Madongming/resource-tree/cache"
	. "github.com/Madongming/resource-tree/global"
	"github.com/Madongming/resource-tree/model"
)

func getCurrentVersion() (int, error) {
	currentVersion := new(model.NodeVersion)
	if err := DB().
		First(currentVersion, 1).
		Error; err != nil {
		return 0, err
	}
	return currentVersion.Current, nil
}

func casVersion() error {
	// select current as c from version where id = 1;
	// update current set current = current + 1 where id = 1 and current = c;

Node_Retry:
	version, err := getCurrentVersion()
	if err != nil {
		return err
	}

	if DB().
		Exec("UPDATE node_versions "+
			"set current = current + 1 "+
			"WHERE id = 1 and current = ?",
			version).RowsAffected == int64(0) {
		goto Node_Retry
	}
	return nil
}

func casResource() (bool, error) {
	version, err := getCurrentVersion()
	if err != nil {
		return false, err
	} else if cache.IsReData || (cache.Tree != nil && cache.Tree.Version == version) {
		// 如果资源正在重建，或者当前资源的版本未更新，则跳过更新资源。
		return false, nil
	}
	data, err := GetAllNodes()
	if err != nil {
		return false, err
	}
	if err := cache.Tree.Set(version, data); err != nil {
		return false, err
	}
	return true, nil
}

func getCurrentEdgeVersion() (int, error) {
	currentVersion := new(model.EdgeVersion)
	if err := DB().
		First(currentVersion, 1).
		Error; err != nil {
		return 0, err
	}
	return currentVersion.Current, nil
}

func casEdgeVersion() error {
	// select current as c from version where id = 1;
	// update current set current = current + 1 where id = 1 and current = c;

Edge_Retry:
	version, err := getCurrentEdgeVersion()
	if err != nil {
		return err
	}
	if DB().
		Exec("UPDATE edge_versions "+
			"set current = current + 1 "+
			"WHERE id = 1 and current = ?",
			version).RowsAffected == int64(0) {
		goto Edge_Retry
	}
	return nil
}

// Fetch all node of the tree. For be used by cache model.
func GetAllNodes() ([]*model.DBResourceNode, error) {
	var resourceNodes []*model.DBResourceNode
	if result := DB().
		Find(&resourceNodes); result.Error != nil {
		if result.RecordNotFound() {
			return []*model.DBResourceNode{}, nil
		}
	}
	return resourceNodes, nil
}

// 获取树，不分权限
func getAllTree() (*model.Tree, error) {
	return cache.Tree.Get()
}

func makeUserTreeIndex(userId interface{}) (*model.Tree, error) {
	// 获取用户所有被显示设定有权限的节点
	permissionSet, err := getUserNodes(userId.(int))
	if err != nil || len(permissionSet) == 0 {
		return nil, err
	}

	// 根据权限集合，从根开始，删掉没有权限的节点，直到遇到又权限的点
	tree, err := cache.NewTreeByPermission(permissionSet)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func getNodeById(nodeId interface{}) (*model.DBResourceNode, error) {
	node := new(model.DBResourceNode)
	if err := DB().
		First(node, nodeId).Error; err != nil {
		return nil, err
	}
	return node, nil
}

func createUser(name, cnName string) error {
	user := &model.DBUser{
		Name:   name,
		CnName: cnName,
	}

	return user.Create()
}

func createGroup(name, cnName string) error {
	group := &model.DBGroup{
		Name:   name,
		CnName: cnName,
	}

	return group.Create()
}

func addUserToGroup(userId, groupId interface{}) error {
	ug := &model.DBUserGroup{
		UserID:  userId.(int),
		GroupID: groupId.(int),
	}
	return ug.Create()
}

func addUserToNode(userId, nodeId interface{}, permission int) error {
	un := &model.DBUserPermission{
		ReadWriteMask: uint(permission),
		NodeID:        nodeId.(int),
		UserID:        userId.(int),
	}
	return un.Create()
}

func addGroupToNode(groupId, nodeId interface{}, permission int) error {
	gn := &model.DBGroupPermission{
		ReadWriteMask: uint(permission),
		NodeID:        nodeId.(int),
		GroupID:       groupId.(int),
	}
	return gn.Create()
}

func getUserById(id interface{}) (*model.DBUser, error) {
	user := new(model.DBUser)
	if result := DB().
		First(user); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return user, nil
}

func getUsersByIds(ids []int) ([]*model.DBUser, error) {
	var users []*model.DBUser
	if ids == nil || len(ids) == 0 {
		return nil, nil
	}
	idsIter := make([]interface{}, len(ids))
	idsQueryStr := make([]string, len(ids))
	for i := range ids {
		idsIter[i] = ids[i]
		idsQueryStr[i] = "?"
	}

	if result := DB().
		Where("`id` in ("+
			strings.Join(idsQueryStr, ",")+
			")",
			idsIter...).
		Find(&users); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return users, nil
}

func getGroupById(id interface{}) (*model.DBGroup, error) {
	group := new(model.DBGroup)
	if result := DB().
		First(group); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return group, nil
}

func getGroupsByIds(ids []int) ([]*model.DBGroup, error) {
	var groups []*model.DBGroup
	if ids == nil || len(ids) == 0 {
		return nil, nil
	}
	idsIter := make([]interface{}, len(ids))
	idsQueryStr := make([]string, len(ids))
	for i := range ids {
		idsIter[i] = ids[i]
		idsQueryStr[i] = "?"
	}

	if result := DB().
		Where("`id` in ("+
			strings.Join(idsQueryStr, ",")+
			")",
			idsIter...).
		Find(&groups); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return groups, nil
}

func getUserByName(name string) (*model.DBUser, error) {
	user := new(model.DBUser)
	if result := DB().
		Where("`name` = ?", name).
		First(user); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return user, nil
}

func getUserGroups(userId interface{}) ([]*model.DBGroup, error) {
	var userGroups []*model.DBUserGroup
	if result := DB().
		Where("`user_id` = ?", userId).
		Find(&userGroups); result.Error != nil {
		if result.RecordNotFound() {
			return []*model.DBGroup{}, nil
		}
		return nil, result.Error
	}

	groupIds := make([]interface{}, len(userGroups))
	groupQuery := make([]string, len(userGroups))
	for i := range userGroups {
		groupIds[i] = userGroups[i].GroupID
		groupQuery[i] = "?"
	}

	var groups []*model.DBGroup
	if result := DB().
		Where("`id` in ("+
			strings.Join(groupQuery, ",")+
			")", groupIds...).
		Find(&groups); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return groups, nil
}

func getNodeUserPermission(userId, nodeId interface{}) (*model.DBUserPermission, error) {
	userPermission := new(model.DBUserPermission)
	if result := DB().
		Where("`user_id` = ? AND "+
			"`node_id` = ?",
			userId, nodeId).
		First(userPermission); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return userPermission, nil
}

func getNodeGroupPermission(groupId, nodeId interface{}) (*model.DBGroupPermission, error) {
	groupPermission := new(model.DBGroupPermission)
	if result := DB().
		Where("`group_id` = ? AND "+
			"`node_id` = ?",
			groupId, nodeId).
		First(groupPermission); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return groupPermission, nil
}

func getGroupNodeIds(groupId interface{}) ([]int, error) {
	var groupPermissions []*model.DBGroupPermission
	if result := DB().
		Where("`group_id` = ?", groupId).
		Find(&groupPermissions); result.Error != nil {
		if result.RecordNotFound() {
			return []int{}, nil
		}
		return nil, result.Error
	}
	nodeIds := make([]int, len(groupPermissions))
	for i := range groupPermissions {
		nodeIds[i] = groupPermissions[i].NodeID
	}
	return nodeIds, nil
}

func getUserNodeIds(userId interface{}) ([]int, error) {
	var userPermissions []*model.DBUserPermission
	if result := DB().
		Where("`user_id` = ?", userId).
		Find(&userPermissions); result.Error != nil {
		if result.RecordNotFound() {
			return []int{}, nil
		}
		return nil, result.Error
	}
	nodeIds := make([]int, len(userPermissions))
	for i := range userPermissions {
		nodeIds[i] = userPermissions[i].NodeID
	}
	return nodeIds, nil
}

func getGroupsNodeIds(groups []*model.DBGroup) ([]int, error) {
	if groups == nil || len(groups) == 0 {
		return []int{}, nil
	}
	groupIds := make([]interface{}, len(groups))
	queryArray := make([]string, len(groups))
	for i := range groups {
		groupIds[i] = groups[i].ID
		queryArray[i] = "?"
	}

	var groupPermissions []*model.DBGroupPermission
	if result := DB().
		Where("`group_id` in ("+
			strings.Join(queryArray, ",")+
			")", groupIds...).
		Find(&groupPermissions); result.Error != nil {
		if result.RecordNotFound() {
			return []int{}, nil
		}
		return nil, result.Error
	}
	nodeIds := make([]int, len(groupPermissions))
	for i := range groupPermissions {
		nodeIds[i] = groupPermissions[i].NodeID
	}
	return nodeIds, nil
}

// 获取user以及其所在组的所有可读节点，用map当set来使用
func getUserNodes(userId interface{}) (map[int]struct{}, error) {
	groups, err := getUserGroups(userId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Get userid:%d's groups is %v\n", userId, groups)
	groupNodeIds, err := getGroupsNodeIds(groups)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Get userid:%d's groupNodeIds is %v\n", userId, groupNodeIds)

	userNodeIds, err := getUserNodeIds(userId)
	if err != nil {
		return nil, err
	}

	// 最多为两者之和
	resultSet := make(map[int]struct{}, len(groupNodeIds)+len(userNodeIds))

	for i := range groupNodeIds {
		resultSet[groupNodeIds[i]] = struct{}{}
	}
	for i := range userNodeIds {
		resultSet[userNodeIds[i]] = struct{}{}
	}
	return resultSet, nil
}

func isUserHasNodePermission(userId, nodeId interface{}) (bool, error) {
	nodeIds := findAllNodesToRoot(nodeId)
	return isUserHaveOneNodePermisson(userId, nodeIds)
}

// 获取一个节点从下至上到root的所有节点
func findAllNodesToRoot(nodeId interface{}) []interface{} {
	id, ok := nodeId.(int)
	if !ok {
		log.Error(ERR_ASSERTION)
	}
	node, _ := cache.Tree.GetTreeNode(id)
	if node == nil ||
		node.Node == nil ||
		node.Node.Level == 0 {
		return []interface{}{}
	}
	permissions := make([]interface{}, node.Node.Level)
	for tn := node; tn.Node.Parent != 0; {
		permissions = append(permissions,
			tn.Node.ID)
		tp, err := cache.Tree.GetTreeNode(tn.Node.Parent)
		if err != nil {
			return []interface{}{}
		}
		tn = tp
	}
	return permissions
}

// 判断节点列表中是否存在一个节点，用户对其有权限
func isUserHaveOneNodePermisson(userId interface{}, nodeIds []interface{}) (bool, error) {
	permissionSet, err := getUserNodes(userId)
	if err != nil || len(permissionSet) == 0 {
		return false, err
	}
	for i := range nodeIds {
		_, found := permissionSet[nodeIds[i].(int)]
		if found {
			return true, nil
		}
	}
	return false, nil
}

func isRelationshipExist(nodeId1, nodeId2 interface{}) bool {
	var count int
	if err := DB().
		Where("(`source_resource_node_id` = ? AND `target_resource_node_id` = ?)"+
			" OR "+
			"(`target_resource_node_id` = ? AND `source_resource_node_id` = ?)",
			nodeId1, nodeId2, nodeId1, nodeId2).
		Count(&count).Error; err != nil || count != 0 {
		return false
	}
	return true
}

func getUserGraph(nodeId interface{}) ([]int, []*model.ResourceNode, []*model.ResourceEdge, error) {
	nodeSet, edgeSet, err := bfsGetNodesAndEdges(nodeId)
	if err != nil {
		return nil, nil, nil, err
	}
	nodeIds := make([]int, len(nodeSet))
	i := 0
	for k, _ := range nodeSet {
		nodeIds[i] = k
		i++
	}

	_, err = casResource()
	if err != nil {
		// 降级，使用缓存中的数据
		log.Error(err)
	}

	index := makeResourceIndex(cache.ResourceNodes.Data)
	nodes := make([]*model.ResourceNode, len(nodeSet))
	j := 0
	for k, _ := range nodeSet {
		nodes[j] = index[k]
		j++
	}

	edges := make([]*model.ResourceEdge, len(edgeSet))
	k := 0
	for _, v := range edgeSet {
		edges[k] = &model.ResourceEdge{
			Source: v[0],
			Target: v[1],
		}
		k++
	}
	return nodeIds, nodes, edges, nil
}

func bfsGetNodesAndEdges(nodeId interface{}) (map[int]struct{}, map[string][2]int, error) {
	q := queue.NewQueue()
	id, ok := nodeId.(int)
	if !ok {
		return nil, nil, ERR_ASSERTION
	}

	nodeSet := make(map[int]struct{})
	edgeSet := make(map[string][2]int)

	q.EnQueue(id)
	for q.Size() != 0 {
		curId, err := q.DeQueue()
		if err != nil {
			return nil, nil, err
		}
		nodeSet[curId] = struct{}{}
		directNodeSet, directEdgeSet := getDirectNodesAndEdges(curId)
		for k, _ := range directNodeSet {
			if _, found := nodeSet[k]; !found {
				q.EnQueue(k)
			}
		}
		for k, v := range directEdgeSet {
			edgeSet[k] = v
		}
	}
	return nodeSet, edgeSet, nil
}

// 获取直连的节点及边的集合
func getDirectNodesAndEdges(id int) (map[int]struct{}, map[string][2]int) {
	var rrs []*model.DBResourceRelationship
	if err := DB().
		Where("(`source_resource_node_id` = ?"+
			" OR "+
			"`target_resource_node_id` = ?)",
			id, id).
		Find(&rrs).Error; err != nil {
		return map[int]struct{}{}, map[string][2]int{}
	}

	nodeSet := make(map[int]struct{}, len(rrs)*2)
	edgeSet := make(map[string][2]int, len(rrs))
	for i := range rrs {
		if rrs[i].SourceResourceNodeID != id {
			nodeSet[rrs[i].SourceResourceNodeID] = struct{}{}
		}
		if rrs[i].TargetResourceNodeID != id {
			nodeSet[rrs[i].TargetResourceNodeID] = struct{}{}
		}
		key := fmt.Sprintf("%d-%d",
			rrs[i].SourceResourceNodeID,
			rrs[i].TargetResourceNodeID)
		edgeSet[key] = [2]int{
			rrs[i].SourceResourceNodeID,
			rrs[i].TargetResourceNodeID,
		}
	}
	return nodeSet, edgeSet
}

// 制作缓存中节点数据的索引
func makeResourceIndex(data []*model.ResourceNode) []*model.ResourceNode {
	indexLen := data[len(data)-1].ID + 1
	index := make([]*model.ResourceNode, indexLen)
	for i := range data {
		index[data[i].ID] = data[i]
	}
	return index
}

func deleteNodeById(nodeId interface{}) error {
	return DB().
		Delete(model.DBResourceNode{},
			"`id` = ?", nodeId).
		Error
}

func deleteResourceRelationshipByNodeId(nodeId interface{}) error {
	return DB().
		Delete(model.DBResourceRelationship{},
			"`source_resource_node_id` = ? OR "+
				"`target_resource_node_id` = ?",
			nodeId, nodeId).
		Error
}
