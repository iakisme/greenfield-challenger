package dao

import (
	"github.com/bnb-chain/greenfield-challenger/db/model"
	"gorm.io/gorm"
)

type BlockDao struct {
	DB *gorm.DB
}

func NewGreenfieldBlockDao(db *gorm.DB) *BlockDao {
	return &BlockDao{
		DB: db,
	}
}

func (d *BlockDao) GetLatestBlock() (*model.Block, error) {
	block := model.Block{}
	err := d.DB.Model(model.Block{}).Order("Height DESC").Take(&block).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &block, nil
}
