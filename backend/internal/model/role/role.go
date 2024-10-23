package role

type Role struct {
	ID          uint   `gorm:"column:id;primaryKey" json:"id"`
	Name        string `gorm:"-" json:"name"`
	Permissions []uint `gorm:"column:permissions" json:"permissions"`
}
