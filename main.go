package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"hack-frenzy.com/model"
)

func main() {
	client, err := FrenzyLogin(Account, Password)
	if err != nil {
		log.Fatal(err)
	}
	db, err := model.SqliteDB("db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	err = model.InitDB(db)
	if err != nil {
		log.Fatal(err)
	}

	var itemCount int64
	err = db.Model(&model.ExamList{}).Count(&itemCount).Error
	if err != nil {
		log.Fatal(err)
	}

	if itemCount == 0 {
		exam, err := client.Load49()
		if err != nil {
			log.Fatal(err)
		}
		for _, exam := range exam {
			db.Create(&model.ExamList{
				Name: exam.Name,
				GID:  exam.GID,
			})
		}
	}

	server := gin.Default()
	route(server, db)

	server.Run(":8870")

}
