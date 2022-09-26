package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type Config struct {
	RemotedbTest    MongoConfig `yaml:"database_test"`
	RemotedbStaging MongoConfig `yaml:"database_staging"`
	LocaldbTest     MysqlConfig `yaml:"database_local_test"`
	LocaldbStaging  MysqlConfig `yaml:"database_local_staging"`
}


type MongoConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Database string `yaml:"database"`
	DBName   string `yaml:"dbname"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Database string `yaml:"database"`
}

func OpenConfigFile() (Config, error) {
	absPath, _ := filepath.Abs("./config.yaml")
	f, err := os.Open(absPath)
	if err != nil {
		return Config{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("open config file error:",err)
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


func  (c *Config) InitializeMongoClient(ctx context.Context) (*mongo.Client, string,error) {
	rt := os.ExpandEnv("${RUNTIME}")
	var clientOptions *options.ClientOptions
	var dbOnline string
	switch rt {
	case "test":
		clientOptions = options.Client().ApplyURI("mongodb://" + c.RemotedbTest.User+ ":" + c.RemotedbTest.Pass + "@" + c.RemotedbTest.Host + ":" + c.RemotedbTest.Port + "/" + c.RemotedbTest.Database)
		dbOnline = c.RemotedbTest.Database
	case "staging":
		clientOptions = options.Client().ApplyURI("mongodb://" + c.RemotedbStaging.User + ":" + c.RemotedbStaging.Pass + "@" + c.RemotedbStaging.Host + ":" + c.RemotedbStaging.Port + "/" + c.RemotedbStaging.Database)
		dbOnline = c.RemotedbStaging.Database
	default:
		fmt.Sprintf("runtime environment mismatch")
	}

	clientOptions.SetMaxPoolSize(50)
	co, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("mongo connect error:%s\n", err)
		return nil, "", err
	}
	err = co.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("ping mongo error:%s\n", err)
		return nil, "", err
	}
	return co, dbOnline ,nil
}



func (c *Config) InitializeMysqlClient()  (*gorm.DB,error ){

	rt := os.ExpandEnv("${RUNTIME}")

	var dsn string
	switch rt {
	case "test":
		dsn = c.LocaldbTest.User + ":" + c.LocaldbTest.Pass + "@tcp(" + c.LocaldbTest.Host + ")/" + c.LocaldbTest.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	case "staging":
		dsn = c.LocaldbStaging.User + ":" + c.LocaldbStaging.Pass + "@tcp(" + c.LocaldbStaging.Host + ")/" + c.LocaldbStaging.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	default:
		fmt.Sprintf("runtime environment mismatch")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db ,nil
}
