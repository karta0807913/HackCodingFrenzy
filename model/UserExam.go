package model

type UserExam struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	UserID    uint     `gorm:"uniqueIndex:UserExamIndex" json:"-"`
	UserData  UserData `gorm:"foreignkey:ID;references:UserID" json:"user_data"`
	ExamID    uint     `gorm:"uniqueIndex:UserExamIndex" json:"-"`
	ExamData  ExamList `gorm:"foreignkey:ID;references:ExamID" json:"exam_data"`
	SecretKey string   `gorm:"index;type:VARCHAR(20)" json:"-"`
	State     uint     `json:"state"`
}
