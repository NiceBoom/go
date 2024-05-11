package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

type WeightRecord struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Weight     int    `json:"weight"`
	CreateTime string `json:"create_date"`
}

func main() {
	//加载并解析配置文件
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file:", err)
	}
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
	router.GET("/getWeightByName", func(c *gin.Context) {
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

	// 添加插入数据的路由
	router.POST("/newWeight", func(c *gin.Context) {
		// 解析请求体中的JSON数据
		var record WeightRecord
		if err := c.ShouldBindJSON(&record); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		// 执行数据库插入操作
		result, err := db.Exec("INSERT INTO users (name, weight, create_date) VALUES (?, ?, ?)", record.Name, record.Weight, currentTime)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to insert record"})
			return
		}

		// 获取插入数据的ID
		id, _ := result.LastInsertId()

		// 返回插入成功的响应
		c.JSON(200, gin.H{"message": "Record inserted", "id": id})
	})

	// 启动Web服务器
	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
