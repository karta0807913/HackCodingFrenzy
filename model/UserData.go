package model

type UserData struct {
	ID        uint   `gorm:"primaryKey" json:"-"`
	StudentID string `gorm:"index;type:VARCHAR(10);uniqueIndex" json:"student_id"`
}
