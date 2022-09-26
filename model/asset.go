package model

import (
	"errors"
	"gorm.io/gorm"
)

type AssetList struct {
	ID            uint   `gorm:"column:id;type:int(11) unsigned NOT NULL AUTO_INCREMENT;PRIMARY_KEY"`
	Asset         string `gorm:"column:asset;type:VARCHAR(255);NOT NULL"`
	TokenId       string `gorm:"column:tokenid;type:VARCHAR(255);NOT NULL"`
	Image         string `gorm:"column:image;type:VARCHAR(255);NOT NULL"`
	Thumbnail     string `gorm:"column:thumbnail;type:VARCHAR(255)"`
	Timestamp     int64 `gorm:"column:timestamp;type:bigint(20);NOT NULL"`
}

func (assetList *AssetList) TableName() string {
	return "asset"
}

type  AssetListDao struct {
	db *gorm.DB
}

func NewAssetListDao(db *gorm.DB) *AssetListDao {
	return &AssetListDao{
		db: db,
	}
}

func (dao *AssetListDao) FindByAsset(asset string) ([]*AssetList, bool, error) {
	var blackList []*AssetList
	err := dao.db.Where("asset = ?", asset).Find(&blackList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return blackList, true, nil
}

func (dao *AssetListDao) BatchesCreate(blackList []*AssetList) error {
	return dao.db.CreateInBatches(blackList,len(blackList)).Error
}

func (dao *AssetListDao) Create(blackList *AssetList) error {
	return dao.db.Create(blackList).Error
}

func (dao *AssetListDao) FindByAssetTokenid(asset string, tokenid string) (*AssetList, bool, error) {
	var banned *AssetList

	err := dao.db.Where("asset = ? AND tokenid = ?", asset, tokenid).First(&banned).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return banned, true, nil
}
