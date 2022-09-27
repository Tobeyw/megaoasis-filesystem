package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

func OpenConfigFile() (Config, error) {
	absPath, _ := filepath.Abs("./config.yaml")
	f, err := os.Open(absPath)
	if err != nil {
		return Config{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Closing file error: %v", err)
		}
	}(f)
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

type Config struct {
	DatabaseMongoTest struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_mongo_test"`
	DatabaseMongoStaging struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_mongo_staging"`
	DatabaseMysqlTest struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
	} `yaml:"database_mysql_test"`
	DatabaseMysqlStaging struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
	} `yaml:"database_mysql_staging"`
}

func (cfg Config) InitializeMongoClient(ctx context.Context) (*mongo.Client, string, error) {
	rt := os.ExpandEnv("${RUNTIME}")
	var clientOptions *options.ClientOptions
	var dbOnline string
	switch rt {
	case "test":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.DatabaseMongoTest.User + ":" + cfg.DatabaseMongoTest.Pass + "@" + cfg.DatabaseMongoTest.Host + ":" + cfg.DatabaseMongoTest.Port + "/" + cfg.DatabaseMongoTest.Database)
		dbOnline = cfg.DatabaseMongoTest.Database
	case "staging":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.DatabaseMongoStaging.User + ":" + cfg.DatabaseMongoStaging.Pass + "@" + cfg.DatabaseMongoStaging.Host + ":" + cfg.DatabaseMongoStaging.Port + "/" + cfg.DatabaseMongoStaging.Database)
		dbOnline = cfg.DatabaseMongoStaging.Database
	default:
		log.Fatalf("runtime environment mismatch")
	}

	clientOptions.SetMaxPoolSize(50)
	co, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("mongo connect error:%s", err)
	}
	err = co.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("ping mongo error:%s", err)
	}
	return co, dbOnline, nil
}

func (cfg Config) InitializeMysqlClient() (*gorm.DB, error) {
	rt := os.ExpandEnv("${RUNTIME}")
	var dsn string

	switch rt {
	case "test":
		dsn = cfg.DatabaseMysqlTest.User + ":" + cfg.DatabaseMysqlTest.Pass + "@tcp(" + cfg.DatabaseMysqlTest.Host + ":" + cfg.DatabaseMysqlTest.Port + ")/" + cfg.DatabaseMysqlTest.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	case "staging":
		dsn = cfg.DatabaseMysqlStaging.User + ":" + cfg.DatabaseMysqlStaging.Pass + "@tcp(" + cfg.DatabaseMysqlStaging.Host + ":" + cfg.DatabaseMysqlStaging.Port + ")/" + cfg.DatabaseMysqlStaging.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	default:
		log.Fatalf("runtime environment mismatch")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
