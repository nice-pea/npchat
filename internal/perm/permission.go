package perm

type Kind string

const (
	Administrator           Kind = "Administrator"
	ManageMemberPermissions Kind = "ManageMemberPermissions"
	AddMembers              Kind = "AddMembers"
	DeleteMembers           Kind = "DeleteMembers"
	DeleteMemberMessages    Kind = "DeleteMemberMessages"
	EditChatInfo            Kind = "EditChatInfo"
	SendMessages            Kind = "SendMessages"
)
