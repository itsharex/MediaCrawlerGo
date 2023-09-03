package conf

import (
	"MediaCrawlerGo/models"
	"log"
	"os"

	"MediaCrawlerGo/util"
	"github.com/joho/godotenv"
)

func Init() {
	// load env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file and error:", err)
	}

	// build custom logger
	util.BuildLogger("MediaCrawlerGo", os.Getenv("LOG_LEVEL"))

	// connected mysql
	// 连接数据库
	models.Database(os.Getenv("MYSQL_DSN"))

}
