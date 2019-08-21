package dao

import (
	"testing"

	"github.com/Madongming/resource-tree/model"
)

func TestGetGroupUsers(t *testing.T) {
	type args struct {
		groupId interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.DBUser
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get group of ID 1.",
			args: args{
				groupId: 1,
			},
			want: []*model.DBUser{
				&model.DBUser{
					Name:   "mdm",
					CnName: "马东明",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGroupUsers(tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGroupUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got[0].Name != tt.want[0].Name ||
				got[0].CnName != tt.want[0].CnName {
				t.Errorf("GetGroupUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTreeNodeUsers(t *testing.T) {
	type args struct {
		nodeId interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.DBUser
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get root Node.",
			args: args{
				nodeId: 1,
			},
			want: []*model.DBUser{
				&model.DBUser{
					ID: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Get App2 Node.",
			args: args{
				nodeId: 5,
			},
			want: []*model.DBUser{
				&model.DBUser{
					ID: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTreeNodeUsers(tt.args.nodeId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTreeNodeUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got[0].ID != tt.want[0].ID {
				t.Errorf("GetTreeNodeUsers() = %#v, want %#v", got[0], tt.want[0])
			}
		})
	}
}

func TestGetTreeNodeGroups(t *testing.T) {
	type args struct {
		nodeId interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.DBGroup
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get group of node(ID:8).",
			args: args{
				nodeId: 8,
			},
			want: []*model.DBGroup{
				&model.DBGroup{
					ID: 2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTreeNodeGroups(tt.args.nodeId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTreeNodeGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got[0].ID != tt.want[0].ID {
				t.Errorf("GetTreeNodeGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTree(t *testing.T) {
	type args struct {
		userId interface{}
		isFull []bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "User mdm3(ID:3), Get full tree.",
			args: args{
				userId: 3,
				isFull: []bool{true},
			},
		},
		{
			name: "User mdm3(ID:3), Get self tree.",
			args: args{
				userId: 3,
				isFull: []bool{},
			},
		},
		{
			name: "User mdm4(ID:4), Get full tree.",
			args: args{
				userId: 4,
				isFull: []bool{true},
			},
		},
		{
			name: "User mdm4(ID:4), Get self tree.",
			args: args{
				userId: 4,
				isFull: []bool{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTree(tt.args.userId, tt.args.isFull...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Name: %s\nGet tree is %s", tt.name, got)
		})
	}
}

func TestGetNodeGraph(t *testing.T) {
	type args struct {
		nodeId interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get node app1 ID:4 graph(App1-4,ID:4-7)",
			args: args{
				nodeId: 4,
			},
			wantErr: false,
		},
		{
			name: "Get node app3 ID:6 graph(App1-4,ID:4-7)",
			args: args{
				nodeId: 6,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeGraph(tt.args.nodeId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeGraph() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Name: %s\n%s", tt.name, got)
		})
	}
}

func TestGetUserPermission(t *testing.T) {
	type args struct {
		userId int
		nodeId int
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get the permission of user mdm(ID:1) and node backend(ID:3)",
			args: args{
				userId: 1,
				nodeId: 3,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserPermission(tt.args.userId, tt.args.nodeId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGroupPermission(t *testing.T) {
	type args struct {
		groupId int
		nodeId  int
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get the permission of group test2(ID:2) and node others(ID:8)",
			args: args{
				groupId: 2,
				nodeId:  8,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGroupPermission(tt.args.groupId, tt.args.nodeId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGroupPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGroupPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}
