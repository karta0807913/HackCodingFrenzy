package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SqliteDB(name string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(name), &gorm.Config{})
}

func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&UserData{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&ExamList{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&UserExam{})
	if err != nil {
		return err
	}
	return nil
}
