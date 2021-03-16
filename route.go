package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hack-frenzy.com/model"
)

func RandStringRunes(n int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func renderIndex(code int, c *gin.Context, data interface{}) {
	tpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error :<"})
		return
	}
	c.Status(code)
	c.Header("Content-Type", "text/html;charset=utf-8")
	err = tpl.Execute(c.Writer, data)
	if err != nil {
		log.Println(err)
	}
}

func route(server *gin.Engine, db *gorm.DB) {
	server.GET("/index", func(c *gin.Context) {
		renderIndex(http.StatusOK, c, gin.H{})
	})

	server.POST("/index", func(c *gin.Context) {
		type Body struct {
			Account  string `form:"account" binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		var body Body
		err := c.ShouldBind(&body)
		if err != nil {
			log.Println(err)
			renderIndex(http.StatusBadRequest, c, gin.H{
				"Error":   "請輸入帳號密碼",
				"Account": body.Account,
			})
			return
		}
		client, err := FrenzyLogin(body.Account, body.Password)
		if err != nil {
			renderIndex(http.StatusBadRequest, c, gin.H{
				"Error":   err.Error(),
				"Account": body.Account,
			})
			return
		}

		var userData model.UserData
		if db.First(&userData, "student_id=?", body.Account).Error != nil {
			userData.StudentID = body.Account
			err = db.Create(&userData).Error
			if err != nil {
				log.Println(err)
				renderIndex(http.StatusBadRequest, c, gin.H{
					"Error":   "建立使用者資料時發生錯誤 X_X",
					"Account": body.Account,
				})
				return
			}
		}

		var list []model.ExamList
		err = db.Find(&list).Error
		if err != nil {
			log.Println(err)
			renderIndex(http.StatusInternalServerError, c, gin.H{
				"Error":   "內部錯誤，請稍後再試",
				"Account": body.Account,
			})
			return
		}
		for _, exam := range list {
			secret := RandStringRunes(20)
			err := client.UpdateExamSave(exam.GID, body.Account+"@"+exam.Name+":"+secret)
			if err != nil {
				log.Println(err)
				renderIndex(http.StatusInternalServerError, c, gin.H{
					"Error":   err.Error() + "，請稍後再試",
					"Account": body.Account,
				})
				return
			}

			err = db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "exam_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"secret_key": secret}),
			}).Create(&model.UserExam{
				UserID:    userData.ID,
				ExamID:    exam.ID,
				SecretKey: secret,
			}).Error
		}
		renderIndex(http.StatusOK, c, gin.H{
			"Error": "成功拉，現在可以打開你的程式了",
		})
	})

	server.GET("/admin", func(c *gin.Context) {
		val, ok := c.GetQuery("key")
		if !ok {
			c.Redirect(http.StatusTemporaryRedirect, "/index")
			return
		}
		if val != "ad@mIn123" {
			c.Redirect(http.StatusTemporaryRedirect, "/index")
			return
		}
		tpl, err := template.ParseFiles("templates/layout.html", "templates/admin.html")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error :<"})
			return
		}
		var students []model.UserData
		var exam []model.ExamList
		err = db.Find(&students).Error
		if err != nil {
			log.Println(err)
			tpl.Execute(c.Writer, gin.H{
				"Result": err.Error(),
			})
			return
		}

		err = db.Find(&exam).Error
		if err != nil {
			log.Println(err)
			tpl.Execute(c.Writer, gin.H{
				"Result": err.Error(),
			})
			return
		}

		result := gin.H{
			"Students": students,
			"Exam":     exam,
		}
		queryStudent, sok := c.GetQuery("student")
		queryExam, eok := c.GetQuery("exam")
		if sok && eok {
			var res model.UserExam
			err = db.Preload("UserData").Preload("ExamData").First(&res, "user_id = ? and exam_id = ?", queryStudent, queryExam).Error
			if err != nil {
				log.Println(err)
				result["Result"] = "紀錄未找到"
			} else {
				result["Result"] = "使用者 " + res.UserData.StudentID + " 題目 [" + res.ExamData.Name + "] 的secret是 " + res.SecretKey
			}
		}

		err = tpl.Execute(c.Writer, result)
		if err != nil {
			log.Println(err)
		}

	})
	server.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/index")
	})
}
