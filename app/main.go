package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"metaoasis-filesystem/api"
	"metaoasis-filesystem/config"
	"metaoasis-filesystem/consts"
	"metaoasis-filesystem/model"
	"os"
)

func main() {
	// loading config file
	cfg, err := config.OpenConfigFile()
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.TODO()
	mongoColl, mongoOnline, err := cfg.InitializeMongoClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	mysqlColl, err := cfg.InitializeMysqlClient()
	if err != nil {
		fmt.Println(err)
	}
	mongoClient := model.T{
		Db_online: mongoOnline,
		C_online:  mongoColl,
		Ctx:       ctx,
	}

	assetDAO := model.NewAssetListDao(mysqlColl)

	fmt.Sprintf(mongoClient.Db_online, assetDAO)

	apiClent := api.T{
		Client:      &mongoClient,
		MysqlClient: assetDAO,
	}
	//listening.....
	//======================
	router := gin.Default()

	/// setting safe proxies

	// using a CDN service
	//router.TrustedPlatform = gin.PlatformGoogleAppEngine
	//router.TrustedPlatform = "X-CDN-IP"

	router.SetTrustedProxies([]string{"127.0.0.1"})
	//router.StaticFile("/favicon.ico", "./image/favicon.ico")
	router.GET("/upload", func(c *gin.Context) {

		copyContext := c.Copy()
		// 异步处理
		go func() {
			asset := copyContext.Query("asset")
			tokenid := copyContext.Query("tokenid")
			imagepath := "./image/" + asset + "/image/" + tokenid
			copyContext.File(imagepath)
		}()

	})
	//watching.....
	go func() {
		rt := os.ExpandEnv("${RUNTIME}")
		market := consts.Market_Main
		switch rt {
		case "test":
			market = consts.Market_Test
		case "staging":
			market = consts.Market_Main
		default:
			fmt.Sprintf("runtime environment mismatch")
		}

		fmt.Println(market)
		conn, err := apiClent.GetCollection(struct{ Collection string }{Collection: "MarketNotification"})
		if err != nil {
			fmt.Println("conn :", err)
		}
		cs, err := conn.Watch(context.TODO(), mongo.Pipeline{})
		//cs, err := conn.Watch(context.TODO(),mongo.Pipeline{bson.D{{"$match", bson.D{{"market", market},{"eventname", "AddAsset"}}}}})
		if err != nil {
			//return nil,err
			fmt.Println("watch error:", err)
		}
		fmt.Println("watching.....")

		for cs.Next(context.TODO()) {
			fmt.Println("watching addAsset")
			var changeEvent map[string]interface{}
			err := cs.Decode(&changeEvent)
			if err != nil {
				log.Fatal(err)
			}
			eventItem := changeEvent["fullDocument"].(map[string]interface{})

			asset := eventItem["asset"].(string)
			event := eventItem["eventname"].(string)

			fmt.Println(event == "AddAsset", asset)
			if event == "AddAsset" {
				assetArr := []string{asset}
				err = apiClent.ScanNep11Data(assetArr)
				if err != nil {
					fmt.Println("watching Error :: scan data err: ", err)
				}
			}
		}
	}()

	//// scan data
	fmt.Println("scaning.....")
	assetArr, err := apiClent.GetMarketWhiteList()
	if err != nil {
		log.Fatal("getwhitelist error: ", err)
	}
	err = apiClent.ScanNep11Data(assetArr)

	if err != nil {
		fmt.Println("Error :: scan data err: ", err)
	}

	router.Run(":8080")

}
