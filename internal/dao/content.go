package dao

import (
	"cms/v2/internal/model"

	"gorm.io/gorm"
)

type ContentDao struct {
	db *gorm.DB
}

type FindParams struct {
	ID       int
	Author   string
	Title    string
	Page     int
	PageSize int
}

func NewContentDao(db *gorm.DB) *ContentDao {
	return &ContentDao{db: db}
}

func (c *ContentDao) Get(params *FindParams) ([]*model.ContentDetail, int64, error) {
	query := c.db.Model(&model.ContentDetail{})

	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}

	if params.Author != "" {
		query = query.Where("author = ?", params.Author)
	}

	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}

	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var page, pageSize = 1, 10

	if params.Page > 0 {
		page = params.Page
	}

	if params.PageSize > 0 {
		pageSize = params.PageSize
	}

	offset := (page - 1) * pageSize
	var data []*model.ContentDetail

	if err := query.Offset(offset).Limit(pageSize).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (c *ContentDao) Create(detail model.ContentDetail) error {
	if err := c.db.Create(&detail).Error; err != nil {
		return err
	}

	return nil
}

func (c *ContentDao) Update(id int, detail model.ContentDetail) error {
	if err := c.db.Where("id = ?", id).Updates(&detail).Error; err != nil {
		return err
	}

	return nil
}

func (c *ContentDao) Delete(id int) error {
	if err := c.db.Where("id = ?", id).Delete(&model.ContentDetail{}).Error; err != nil {
		return err
	}

	return nil
}

func (c *ContentDao) IsExist(ID int) (bool, error) {
	err := c.db.Where("id=?", ID).First(&model.ContentDetail{}).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
