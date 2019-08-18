package dao

import (
	"strings"

	log "github.com/cihub/seelog"

	"github.com/satori/go.uuid"

	"cache"
	. "global"
	"model"
)

func getCurrentVersion() (int, error) {
	var version int64
Retry:
	if err := DB().
		Raw("SELECT current FROM node_versions WHERE id = 1").
		Scan(&version).
		Error; err != nil {
		return int64(0), err
	}
	return int(version), nil
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
		Exec("UPDATE node_versions "+
			"set current = current + 1 "+
			"WHERE id = 1 and current = ?",
			version).RowsAffected == int64(0) {
		goto Retry
	}
	return nil
}

func casResource() (bool, error) {
	version, err := getCurrentVersion()
	if err != nil {
		return false, err
	} else if cache.IsReData || (cache.Tree != nil && cache.Tree.Resource.Version == version) {
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
	var version int64
Retry:
	if err := DB().
		Raw("SELECT current FROM edge_versions WHERE id = 1").
		Scan(&version).
		Error; err != nil {
		return int64(0), err
	}
	return int(version), nil
}

func casEdgeVersion() error {
	// select current as c from version where id = 1;
	// update current set current = current + 1 where id = 1 and current = c;

Retry:
	version, err := getCurrentEdgeVersion()
	if err != nil {
		return err
	}

	if DB().
		Exec("UPDATE edge_versions "+
			"set current = current + 1 "+
			"WHERE id = 1 and current = ?",
			version).RowsAffected == int64(0) {
		goto Retry
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
	if err != nil {
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
	user := &User{
		Name:   name,
		CnName: cnName,
	}

	return user.create()
}

func createGroup(name, cnName string) error {
	group := &Group{
		Name:   name,
		CnName: cnName,
	}

	return group.create()
}

func addUserToGroup(userId, groupId interface{}) error {
	ug := &UserGroup{
		UserID:  userId.(int),
		GroupID: groupId.(int),
	}
	return ug.create()
}

func addUserToNode(userId, nodeId interface{}, permission int) error {
	un := &UserPermission{
		ReadWriteMask: permission,
		NodeID:        nodeId.(int),
		UserID:        userId.(int),
	}
	return un.create()
}

func addGroupToNode(groupId, nodeId interface{}, permission int) error {
	gn := &UserPermission{
		ReadWriteMask: permission,
		NodeID:        nodeId.(int),
		GroupID:       groupId.(int),
	}
	return gn.create()
}

func getUserById(id interface{}) (*DBUser, error) {
	user := new(DBUser)
	if result := DB().
		First(user); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return user, nil
}

func getUsersByIds(ids []int) ([]*DBUser, error) {
	var users []*DBUser
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

func getGroupById(id interface{}) (*DBGroup, error) {
	group := new(DBGroup)
	if result := DB().
		First(group); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return group, nil
}

func getGroupsByIds(ids []int) ([]*DBGroup, error) {
	var groups []*DBGroup
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

func getUserByName(name string) (*DBUser, error) {
	user := new(DBUser)
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

func getUserGroups(userId interface{}) ([]*DBGroup, error) {
	var userGroups []*DBUserGroup
	if result := DB().
		Where("`group_id` = ?", groupId).
		Find(userGroups); result.Error != nil {
		if result.RecordNotFound() {
			return []*Group{}, nil
		}
		return nil, result.Error
	}

	groupIds := make([]interface{}, len(userGroups))
	groupQuery := make([]string, len(userGroups))
	for i := range userGroups {
		groupIds[i] = userGroups[i].UserID
		groupQuery[i] = "?"
	}

	var groups []*DBGroup
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

func getNodeUserPermission(userId, nodeId interface{}) (*DBUserPermission, error) {
	userPermission := new(DBUserPermission)
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

func getNodeGroupPermission(groupId, nodeId interface{}) (*DBGroupPermission, error) {
	groupPermission := new(DBGroupPermission)
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
	var groupPermissions []*DBGroupPermission
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
	return nodeIds
}

func getUserNodeIds(userId interface{}) ([]int, error) {
	var userPermissions []*DBUserPermission
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
	return nodeIds
}

func getGroupsNodeIds(groups []*DBGroup) ([]int, error) {
	if groups == nil || len(groups) == 0 {
		return []int{}, nil
	}
	groupIds := make([]interface{}, len(groups))
	queryArray := make([]string, len(groups))
	for i := range groups {
		groupIds[i] = groups[i].ID
		queryArray[i] = "?"
	}

	var groupPermissions []*DBGroupPermission
	if result := DB().
		Where("`group_id` in (" +
			strings.Join(queryArray, ","+
				")", groupIds...)).
		Find(&groupPermissions); result.Error != nil {
		if result.RecordNotFound() {
			return []int{}, nil
		}
		return nil, err
	}
	nodeIds := make([]int, len(groupPermissions))
	for i := range groupPermissions {
		nodeIds[i] = groupPermissions[i].NodeID
	}
	return nodeIds, nil
}

// 获取user以及其所在组的所有可读节点，用map当set来使用
func getUserNodes(userId int) (map[int]struct{}, error) {
	groups, err := getUserGroups(userId)
	if err != nil {
		return nil, err
	}

	groupNodeIds, err := getGroupsNodeIds(groups)
	if err != nil {
		return nil, err
	}
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
	node, _ := cache.GetTreeNode(nodeId)
	if node == nil || node.Level == 0 {
		return []interface{}{}
	}
	permissions := make([]interface{}, node.Level)
	for tn := node; tn.Parent != 0; {
		permissions = append(permissions,
			tn.ID)
		tp, err := cache.GetTreeNode(tn.Parent)
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

func (n *model.DBResourceNode) create() error {
	return DB().Create(n).Error
}

func (n *model.DBResourceNode) update() error {
	return DB().Update(n).Error
}

func (n *model.DBResourceNode) delete() error {
	return DB().Delete(n).Error
}

func (u *model.DBUser) create() error {
	return DB().Create(u).Error
}

func (g *model.DBGroup) create() error {
	return DB().Create(g).Error
}

func (ug *DBUserGroup) create() error {
	return DB().Create(ug).Error
}

func (gn *model.DBGroupPermission) create() error {
	return DB().Create(gn).Error
}

func (un *DBUserPermission) create() error {
	return DB().Create(un).Error
}

func (u *DBUser) update() error {
	return DB().Save(u).Error
}

func (g *DBGroup) update() error {
	return DB().Save(g).Error
}

func (ug *DBUserGroup) update() error {
	return DB().Save(ug).Error
}

func (gp *DBGroupPermission) update() error {
	return DB().Save(gp).Error
}

func (up *DBUserPermission) update() error {
	return DB().Save(up).Error
}

func (rr *DBResourceRelationship) create() error {
	return DB().Create(rr).Error
}

func (rr *DBResourceRelationship) delete() error {
	return DB().Delete(rr).Error
}

// 获取指定节点的所有相关节点ID
func getAllRelationshipEdges(nodeId interface{}) ([]*ResourceEdge, error) {
	var rrs []*DBResourceRelationship
	if result := DB().
		Where("`source_resource_node_id` = ? OR"+
			"`target_desource_node_id = ?`",
			nodeId, nodeId).
		Find(&rrs); result.Error != nil {
		if result.RecordNotFound() {
			return []*ResourceEdge{}, nil
		}
		return nil, err
	}
	var results []*ResourceEdge
	for i := range rrs {
		results = append(results,
			&ResourceEdge{
				Source: rrs.SourceResourceNodeID,
				Target: rrs.SourceResourceNodeID,
			})
	}
	return results, nil
}

// 获取边列表中的所有节点
func getAllRelationshipNodes(rrs []*ResourceEdge) ([]*ResourceNode, error) {
	// 最多2倍的ID个数
	tmp := make(map[int]struct{}, len(rrs)*2)
	for i := range rrs {
		tmp[rrs.Source] = struct{}{}
		tmp[rrs.Target] = struct{}{}
	}

	casResource()
	index := makeResourceIndex(cache.ResourceNodes.Data)
	var results []*ResourceNode
	for k, _ := range tmp {
		results = append(results,
			index[k])
	}
	return results, nil
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
		Delete(DBResourceNode{},
			"`id` = ?", nodeId).
		Error
}

func deleteResourceRelationshipByNodeId(nodeId interface{}) error {
	return DB().
		Delete(DBResourceRelationship{},
			"source_resource_node_id = ? OR "+
				"target_resource_node_id = ?",
			nodeId, nodeId).
		Error
}
