package app

import (
	"context"
	"database/sql"
	"flag"
	"go-nuxt-blogs/db"
	"go-nuxt-blogs/models"
	"go-nuxt-blogs/pkg/config"
	"go-nuxt-blogs/pkg/luckid"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fzzp/gotk"
	"github.com/fzzp/gotk/token"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	Wg     sync.WaitGroup
	JWT    token.Maker
	Conf   *config.Config
	Repo   db.Repository
	FileUp models.BlobUploader
}

func Start() {
	var err error
	var conf config.Config

	// 优先启动验证器，解析json配置需要使用其验证
	gotk.InitValidation("zh")

	var configFile string
	flag.StringVar(&configFile, "path", "app.json", "config 配置文件地址")
	flag.Parse()
	if err = config.NewConfig(configFile, &conf); err != nil {
		log.Fatalln("解析JSON配置失败: ", err)
	}

	// 打印查看一下
	// conf.Println()

	gotk.NewLogger(conf.LogLevel, true, os.Stdout)

	// 连接SQLite，设置基础参数
	conn := openDB()
	defer conn.Close()
	conn.SetMaxOpenConns(conf.MaxOpenConns)
	conn.SetMaxIdleConns(conf.MaxIdleConns)
	conn.SetConnMaxLifetime(conf.MaxLifetime.Duration)
	conn.SetConnMaxIdleTime(conf.MaxIdleTime.Duration)

	repo := db.NewRepository(conn)

	err = luckid.New(conf.LuckFile, 3, 33)
	if err != nil {
		log.Fatal(err)
	}
	defer luckid.Save()

	jwt, err := token.NewJWTMaker(conf.TokenSecret, conf.TokenIssuer)
	if err != nil {
		log.Fatal(err)
	}

	maxFileSize := 4 << 20
	fileUp := models.NewBlobUploader(int64(maxFileSize), conf.AllowFiles...)

	app := application{
		JWT:    jwt,
		Repo:   repo,
		Conf:   &conf,
		FileUp: fileUp,
	}

	if err = app.serve(); err != nil {
		log.Fatal(err)
	}
}

func openDB() *sql.DB {
	conn, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalln("openSQLite() failed: ", err)
	}

	// 连接到SQLite3,最多尝试连接5次
	for i := 5; i > 0; i-- {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := conn.PingContext(ctx)
		if i <= 1 && err != nil {
			log.Fatalln("db.PingContext() failed: ", err)
		}
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return conn
}
