package model

type ExamList struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	GID  string `json:"-"`
}
