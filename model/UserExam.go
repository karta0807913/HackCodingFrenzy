package model

type UserExam struct {
	ID        uint     `gorm:"primaryKey"`
	UserID    uint     `gorm:"uniqueIndex:UserExamIndex"`
	UserData  UserData `gorm:"foreignkey:ID;references:UserID"`
	ExamID    uint     `gorm:"uniqueIndex:UserExamIndex"`
	ExamData  ExamList `gorm:"foreignkey:ID;references:ExamID"`
	SecretKey string   `gorm:"index;type:VARCHAR(20)"`
	State     uint
}
