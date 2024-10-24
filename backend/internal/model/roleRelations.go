package model

type RoleRelations struct {
	RoleID   uint `gorm:"primaryKey" json:"role_id"`
	MemberID uint `gorm:"primaryKey" json:"member_id"`
}
