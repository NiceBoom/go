package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

type WeightRecord struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Weight     int    `json:"weight"`
	CreateTime string `json:"create_date"`
}

func main() {
	//拼接数据库连接字符串
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.dbname")

	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// 连接数据库
	db, err := sql.Open("mysql", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 初始化Gin引擎
	router := gin.Default()

	// 定义路由处理程序
	router.GET("/weight", func(c *gin.Context) {
		// 获取查询参数 name
		name := c.Query("name")

		// 执行数据库查询
		query := "SELECT id, name, weight, create_date FROM users WHERE name = ?"
		rows, err := db.Query(query, name)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		// 解析查询结果
		var weightRecords []WeightRecord
		for rows.Next() {
			var record WeightRecord
			err := rows.Scan(&record.ID, &record.Name, &record.Weight, &record.CreateTime)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			weightRecords = append(weightRecords, record)
		}

		// 返回查询结果
		c.JSON(http.StatusOK, weightRecords)
	})

	// 启动Web服务器
	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
