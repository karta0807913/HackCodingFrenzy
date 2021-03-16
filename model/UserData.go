package model

type UserData struct {
	ID        uint   `gorm:"primaryKey"`
	StudentID string `gorm:"index;type:VARCHAR(10);uniqueIndex"`
}
