package dao

import (
	"testing"
)

func TestCreateGroup(t *testing.T) {
	type args struct {
		name   string
		cnName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Add Group `test`, ID: 1",
			args: args{
				name:   "test",
				cnName: "测试",
			},
		},
		{
			name: "Add Group test2, ID: 2",
			args: args{
				name:   "test2",
				cnName: "测试2",
			},
		},
		{
			name: "Add Group test3, ID: 3",
			args: args{
				name:   "test3",
				cnName: "测试3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateGroup(tt.args.name, tt.args.cnName); (err != nil) != tt.wantErr {
				t.Errorf("CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddUser(t *testing.T) {
	type args struct {
		name   string
		cnName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Add user mdm. id: 1",
			args: args{
				name:   "mdm",
				cnName: "马东明",
			},
		},
		{
			name: "Add user mdm2. id: 2",
			args: args{
				name:   "mdm2",
				cnName: "马东明2",
			},
		},
		{
			name: "Add user mdm3. id: 3",
			args: args{
				name:   "mdm3",
				cnName: "马东明3",
			},
		},
		{
			name: "Add user mdm4. id: 4",
			args: args{
				name:   "mdm4",
				cnName: "马东明4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddUser(tt.args.name, tt.args.cnName); (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddUserToGroup(t *testing.T) {
	type args struct {
		userId  interface{}
		groupId interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Add mdm to test.",
			args: args{
				userId:  1,
				groupId: 1,
			},
		},
		{
			name: "Add mdm2 to test.",
			args: args{
				userId:  2,
				groupId: 1,
			},
		},
		{
			name: "Add mdm3 to test2.",
			args: args{
				userId:  3,
				groupId: 2,
			},
		},
		{
			name: "Add mdm4 to test3.",
			args: args{
				userId:  4,
				groupId: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddUserToGroup(tt.args.userId, tt.args.groupId); (err != nil) != tt.wantErr {
				t.Errorf("AddUserToGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateNode(t *testing.T) {
	type args struct {
		name        string
		description string
		userId      int
		parentId    int
		opts        []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Add root node: MyCompany.ID:1",
			args: args{
				name:        "MyCompany",
				description: "我的公司",
				userId:      1,
				parentId:    0,
			},
			wantErr: false,
		},
		{
			name: "Add child of root node: front.ID:2",
			args: args{
				name:        "front",
				description: "前端组",
				userId:      1,
				parentId:    1,
			},
			wantErr: false,
		},
		{
			name: "Add child of root node: backend.ID:3",
			args: args{
				name:        "backend",
				description: "后端组",
				userId:      1,
				parentId:    1,
				opts:        []interface{}{nil, nil, 1, nil, nil, "group=backend"},
			},
			wantErr: false,
		},
		{
			name: "Add child of backend node: app1.ID:4",
			args: args{
				name:        "app1",
				description: "App 1",
				userId:      1,
				parentId:    3,
			},
			wantErr: false,
		},
		{
			name: "Add child of backend node: app2.ID:5",
			args: args{
				name:        "app2",
				description: "App 2",
				userId:      1,
				parentId:    3,
			},
			wantErr: false,
		},
		{
			name: "Add child of backend node: app3.ID:6",
			args: args{
				name:        "app3",
				description: "App 3",
				userId:      1,
				parentId:    3,
			},
			wantErr: false,
		},
		{
			name: "Add child of backend node: app4.ID:7",
			args: args{
				name:        "app4",
				description: "App 4",
				userId:      1,
				parentId:    3,
			},
			wantErr: false,
		},
		{
			name: "Add child of root node: others.ID:8",
			args: args{
				name:        "others",
				description: "Others",
				userId:      2,
				parentId:    1,
			},
			wantErr: false,
		},
		{
			name: "Add child of others node: app5.ID:9",
			args: args{
				name:        "app5",
				description: "App 5",
				userId:      2,
				parentId:    8,
			},
			wantErr: false,
		},
		{
			name: "Add child of others node: app6.ID:10",
			args: args{
				name:        "app6",
				description: "App 6",
				userId:      2,
				parentId:    8,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNode(tt.args.name, tt.args.description, tt.args.userId, tt.args.parentId, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("CreateNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddNodeToNode(t *testing.T) {
	type args struct {
		srcNodeId interface{}
		tarNodeId interface{}
		userId    []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Add relation to App1 and App2.",
			args: args{
				srcNodeId: 4,
				tarNodeId: 5,
			},
			wantErr: false,
		},
		{
			name: "Add relation to App1 and App2.",
			args: args{
				srcNodeId: 4,
				tarNodeId: 6,
			},
			wantErr: false,
		},
		{
			name: "Add relation to App1 and App2.",
			args: args{
				srcNodeId: 5,
				tarNodeId: 7,
			},
			wantErr: false,
		},
		{
			name: "Add relation to App1 and App2.",
			args: args{
				srcNodeId: 6,
				tarNodeId: 7,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddNodeToNode(tt.args.srcNodeId, tt.args.tarNodeId, tt.args.userId...); (err != nil) != tt.wantErr {
				t.Errorf("AddNodeToNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddGroupToNode(t *testing.T) {
	type args struct {
		groupId     interface{}
		nodeId      int
		permissions []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Add group test2(ID:2) to node others(ID:8)",
			args: args{
				groupId: 2,
				nodeId:  8,
			},
			wantErr: false,
		},
		{
			name: "Add group test3(ID:3) to node front(ID:2)",
			args: args{
				groupId: 3,
				nodeId:  2,
			},
			wantErr: false,
		},
		{
			name: "Add group test3(ID:3) to node app4(ID:7)",
			args: args{
				groupId: 3,
				nodeId:  7,
			},
			wantErr: false,
		},
		{
			name: "Add group test3(ID:3) to node others(ID:8)",
			args: args{
				groupId: 3,
				nodeId:  8,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddGroupToNode(tt.args.groupId, tt.args.nodeId, tt.args.permissions...); (err != nil) != tt.wantErr {
				t.Errorf("AddGroupToNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
