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
	"net/http"
	"os"
	"time"
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
	fmt.Println(apiClent)
	//listening.....
	//======================
	router := gin.Default()

	/// setting safe proxies

	// using a CDN service
	//router.TrustedPlatform = gin.PlatformGoogleAppEngine
	//router.TrustedPlatform = "X-CDN-IP"
	//
	//router.SetTrustedProxies([]string{"127.0.0.1"})

	//router.StaticFile("/","./image/")
	router.GET("/images/:asset/:tokenid", func(c *gin.Context) {
		//image := "image"
		pwd, _ := os.Getwd()
		//copyContext := c.Copy()
		time.Sleep(1 * time.Second)

		asset := c.Param("asset")
		tokenid := c.Param("tokenid")

		//imagepath := pwd + "\\image\\" + asset + "\\image\\" + tokenid
		imagepath := pwd + "/image/" + asset + "/image/" + tokenid
		fmt.Println(imagepath)
		c.File(imagepath)

	})

	router.GET("/thumbnail/:asset/:tokenid", func(c *gin.Context) {
		//image := "image"
		pwd, _ := os.Getwd()
		//copyContext := c.Copy()
		time.Sleep(1 * time.Second)
		p := c.Params
		u := c.Request.RequestURI
		fmt.Println("param: ", p, u)
		asset := c.Param("asset")
		tokenid := c.Param("tokenid")
		//imagepath := pwd + "\\image\\" + asset + "\\image\\" + tokenid
		imagepath := pwd + "/image/" + asset + "/thumbnail/" + tokenid
		fmt.Println(imagepath)
		c.File(imagepath)
	})

	router.POST("/rename", func(c *gin.Context) {
		srcdir := c.PostForm("srcdir")
		dstdir := c.PostForm("dstdir")
		err := api.ImagRename(srcdir, dstdir)

		if err != nil {
			c.String(http.StatusBadRequest, "copyRename err", err)
		} else {
			c.String(http.StatusOK, "copyRename success")
		}

	})

	router.POST("/upload", func(c *gin.Context) {
		asset := c.PostForm("asset")
		tokenid := c.PostForm("tokenid")
		isImage := c.PostForm("isImage")
		image := ""
		thumbnail := ""

		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		currentPath, err := os.Getwd()
		if err != nil {
			log.Fatal("get current path error :", err)
		}
		imagePath := api.CreateDateDir(currentPath+"/image/", asset+"/image/")
		image = imagePath
		if isImage == "false" {
			imagePath = api.CreateDateDir(currentPath+"/image/", asset+"/thumbnail/")
			thumbnail = imagePath
			image = ""
		}

		if err := c.SaveUploadedFile(file, imagePath+"/"+tokenid); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}

		nft := model.AssetList{
			Asset:     asset,
			TokenId:   tokenid,
			Image:     image,
			Thumbnail: thumbnail,
			Timestamp: time.Now().Unix(),
		}
		err = apiClent.MysqlClient.Create(&nft)
		if err != nil {
			c.String(http.StatusOK, "insert err")
		} else {
			c.String(http.StatusOK, "Uploaded successfully %d files with fields name=%s and email=%s.")
		}

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
	//assetArr := []string{}
	if err != nil {
		log.Fatal("getwhitelist error: ", err)
	}
	err = apiClent.ScanNep11Data(assetArr)

	if err != nil {
		fmt.Println("Error :: scan data err: ", err)
	}

	router.Run(":8080")

}
