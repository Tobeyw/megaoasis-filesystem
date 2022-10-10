package api

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"log"
	"metaoasis-filesystem/consts"
	"metaoasis-filesystem/model"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type T struct {
	Client      *model.T
	MysqlClient *model.AssetListDao
}

func (me *T) GetCollection(args struct {
	Collection string
}) (*mongo.Collection, error) {
	collection := me.Client.C_online.Database(me.Client.Db_online).Collection(args.Collection)
	return collection, nil
}

func (me *T) ScanNep11Data(assetArr []string) error {
	//

	fmt.Println("into scan data...")
	message := make(json.RawMessage, 0)
	ret := &message
	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Nep11Properties",
		Index:      "scandata",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": bson.M{"$in": assetArr}, "properties": bson.M{"$ne": "{}"}},
		Query:      []string{},
	}, ret)

	if err != nil {
		return err
	}
	result := make([]*model.AssetList, 0)
	currentPath, err := os.Getwd()
	for _, item := range r1 {
		//获取nft 属性
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		nftproperties := item["properties"]
		image := ""
		thumbnail := ""
		if nftproperties != nil && nftproperties != "" {
			extendData := nftproperties.(string)
			if extendData != "" {
				var data map[string]interface{}
				if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
					img, ok := data["image"]
					if ok {
						image = img.(string)
					}
					tokenuri, ok := data["tokenURI"]
					if ok {
						if image == "" {
							checkpath := currentPath + "/image/" + asset + "/" + tokenid
							isExit, _ := PathExists(checkpath)
							if isExit {
								continue
							}
							image, err = GetImgFromTokenURL(tokenurl(tokenuri.(string)))
							if err != nil {
								return err
							}
						}
					}
					thb, ok6 := data["thumbnail"]
					if ok6 {
						tb, err22 := base64.URLEncoding.DecodeString(thb.(string))
						if err22 != nil {
							return err22
						}
						thumbnail = string(tb[:])
					}
				} else {
					return err
				}

			}
		}

		if image != "" {
			nft := model.AssetList{
				Asset:     asset,
				TokenId:   tokenid,
				Image:     image,
				Thumbnail: thumbnail,
				Timestamp: time.Now().Unix(),
			}
			result = append(result, &nft)
		}
	}
	// insert mysql database
	//me.MysqlClient.BatchesCreate(result)
	for i, item := range result {
		go LoadAndSave(me, item)
		if i == len(result)-1 {
		}
	}

	return nil
}

func GetImgFromTokenURL(tokenurl string) (string, error) {
	var image string
	response, err := http.Get(tokenurl)
	if err != nil {
		log.Println("http get error: ", err)
		return "", err
	}
	defer response.Body.Close()
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		log.Println("ioutil read error: ", err)
	}
	jsonData := make(map[string]interface{})
	//fmt.Println(string(body))
	err = json.Unmarshal([]byte(string(body)), &jsonData)
	if err != nil {
		log.Println("imag from json error :", err, tokenurl)
		return "", nil
	}
	//处理ipfs
	ipfs := jsonData["image"].(string)
	str := strings.Replace(ipfs, ":", "", 1)
	image = "https://" + consts.IPFS_GATEWAY + "/" + str
	return image, nil
}
func LoadAndSave(me *T, list *model.AssetList) error {
	image := list.Image
	thumbnail := list.Thumbnail
	asset := list.Asset
	tokenid := list.TokenId
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}
	//查看本地是否存在数据
	checkpath := currentPath + "/image/" + asset + "/" + tokenid
	isExit, _ := PathExists(checkpath)
	if !isExit {
		if err != nil {
			return err
		}
		// download image
		if image != "" {
			imagePath := CreateDateDir(currentPath+"/image/", asset+"/image/")
			path := imagePath + "/" + tokenid
			err := LoadImage(image, path)
			if err != nil {
				return err
			}
			list.Image = "/image/" + asset + "/image/" + tokenid
		}
		if thumbnail != "" {
			thumbnailPath := CreateDateDir(currentPath+"/image/", asset+"/thumbnail/")
			path := thumbnailPath + "/" + tokenid
			err := LoadImage(thumbnail, path)
			if err != nil {
				return err
			}
			list.Thumbnail = "/image/" + asset + "/thumbnail/" + tokenid
		}

		//err = me.MysqlClient.Create(list)
		//if err != nil {
		//	return err
		//}
		fmt.Println("update one record successfully")
	}

	return nil
}
func LoadImage(imagurl string, path string) error {

	_, err := url.ParseRequestURI(imagurl)
	if err != nil {
		panic(err)
	}
	client := http.DefaultClient
	//client.Timeout = time.Second * 120 //设置超时时间
	resp, err := client.Get(imagurl)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch URL %q:%q,%v", imagurl, path, err))
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(fmt.Errorf("unexpected http GET status %q:%q,%s", imagurl, path, resp.Status))
	}
	if resp.ContentLength <= 0 {
		log.Println("[*] Destination server does not support breakpoint download.")
	}

	raw := resp.Body
	defer raw.Close()

	//reader := bufio.NewReaderSize(raw, 1024*64)

	out, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	wt := bufio.NewWriter(out)

	defer out.Close()

	n, err := io.Copy(wt, resp.Body)
	fmt.Println("write", n)
	if err != nil {
		panic(err)
	}
	wt.Flush()

	return nil
}

func CreateDateDir(basepath string, folderName string) string {

	folderPath := filepath.Join(basepath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0777)
		if err != nil {
			fmt.Println("Create dir error: %v", err)
		}
		err = os.Chmod(folderPath, 0777)
		if err != nil {
			fmt.Println("Chmod error: %v", err)
		}
	}
	return folderPath
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true, nil
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	//其它类型，不确定是否存在
	return false, err
}

func tokenurl(url string) string {

	str := url[:4]
	if str == "ipfs" {
		gateway := "https://cloudflare-ipfs.com"
		httpurl := strings.Replace(url, "ipfs.io", gateway, 1)
		return httpurl
	}

	return url
}
