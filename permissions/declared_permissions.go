package permissions

type DeclaredPermission struct {
	PermissionCode string
	PermissionName string
	PermissionType string
	ResourceCode   string
	GroupCode      string
	Description    string
}

func DeclaredPermissions() []DeclaredPermission {
	return []DeclaredPermission{
		{
			PermissionCode: "role:read",
			PermissionName: "查看角色",
			PermissionType: "READ",
			ResourceCode:   "role",
			GroupCode:      "SYSTEM",
			Description:    "角色列表与详情查看",
		},
		{
			PermissionCode: "role:write",
			PermissionName: "管理角色",
			PermissionType: "WRITE",
			ResourceCode:   "role",
			GroupCode:      "SYSTEM",
			Description:    "角色创建、编辑、删除、分配权限",
		},
		{
			PermissionCode: "permission:read",
			PermissionName: "查看权限字典",
			PermissionType: "READ",
			ResourceCode:   "permission",
			GroupCode:      "SYSTEM",
			Description:    "查看权限码与权限分组信息",
		},
		{
			PermissionCode: "permission:write",
			PermissionName: "管理权限字典",
			PermissionType: "WRITE",
			ResourceCode:   "permission",
			GroupCode:      "SYSTEM",
			Description:    "维护权限字典与权限元数据",
		},
	}
}
