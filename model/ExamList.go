package model

type ExamList struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	GID  string
}
