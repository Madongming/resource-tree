package dao

import (
	"testing"
)

func TestUpdateNodeName(t *testing.T) {
	type args struct {
		name   string
		cnName string
		nodeId interface{}
		userId []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "None user,Modify App1(ID:4)",
			args: args{
				name:   "NewApp1",
				cnName: "新App1",
				nodeId: 4,
			},
			wantErr: false,
		},
		{
			name: "Use user mdm4(ID:4). Modify App4(ID:7) permission resource.",
			args: args{
				name:   "NewApp4",
				nodeId: 7,
				userId: []interface{}{4},
			},
			wantErr: false,
		},
		{
			name: "Use user mdm4(ID:4), modify App2(ID:5) none permission resource.",
			args: args{
				name:   "NewApp2",
				cnName: "新App2",
				nodeId: 5,
				userId: []interface{}{4},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateNodeName(tt.args.name, tt.args.cnName, tt.args.nodeId, tt.args.userId...); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNodeName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	tree, err := GetTree(1, true)
	if err != nil {
		t.Errorf("Finished modify, Error: %v", err)
	}
	t.Logf("Finish modify, tree is %s", tree)
}

func TestUpdataUserNodePermissions(t *testing.T) {
	type args struct {
		userId      interface{}
		nodeId      interface{}
		permissions int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Use user mdm1(ID:1). Modify App4(ID:7) permission resource.",
			args: args{
				userId:      1,
				nodeId:      7,
				permissions: 3,
			},
			wantErr: false,
		},
		{
			name: "Use user mdm4(ID:4), Modify App2(ID:5) none permission resource.",
			args: args{
				userId:      4,
				nodeId:      5,
				permissions: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdataUserNodePermissions(tt.args.userId, tt.args.nodeId, tt.args.permissions); (err != nil) != tt.wantErr {
				t.Errorf("UpdataUserNodePermissions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdataGroupNodePermissions(t *testing.T) {
	type args struct {
		groupId     interface{}
		nodeId      interface{}
		permissions int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Use group test1(ID:2). Modify Others(ID:8) permission resource.",
			args: args{
				groupId:     2,
				nodeId:      8,
				permissions: 3,
			},
			wantErr: false,
		},
		{
			name: "Use group test2(ID:3), Modify App2(ID:5) none permission resource.",
			args: args{
				groupId:     4,
				nodeId:      5,
				permissions: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdataGroupNodePermissions(tt.args.groupId, tt.args.nodeId, tt.args.permissions); (err != nil) != tt.wantErr {
				t.Errorf("UpdataGroupNodePermissions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
