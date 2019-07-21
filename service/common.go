package service

import (
	"model"
)

func changeUserssModel2Service(modelUsers []*model.User) []*User {
	if modelUsers == nil || len(modelUsers) == 0 {
		return []*User{}
	}
	users := make([]*User, len(modelUsers))
	for i := range modelUsers {
		users[i] = &User{
			Base.ID:       int64(modelUsers[i].ID),
			Base.RootNode: int64(modelUsers[i].RootNode),
			Base.Name:     modelUsers.Name,
			Base.CnName:   modelUsers.CnName,
		}
	}
	return users, nil
}

func changeGroupsModel2Service(modelGroups []*model.Group) []*Group {
	if modelGroups == nil || len(modelGroups) == 0 {
		return []*Group{}
	}
	groups := make([]*Group, len(modelGroups))
	for i := range modelGroups {
		groups[i] = &Group{
			Base.ID:       int64(modelGroups[i].ID),
			Base.RootNode: int64(modelGroups[i].RootNode),
			Base.Name:     modelGroups.Name,
			Base.CnName:   modelGroups.CnName,
		}
	}
	return groups, nil
}
