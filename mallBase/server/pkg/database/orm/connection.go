package orm

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"gorm.io/gorm"
)

// 分页条件
type PageWhere struct {
	Where string
	Value []string
}

type CollectionInfo struct {
	client *gorm.DB
}

// 分页参数返回
type IndexPage struct {
	Total    int64  `json:"total"`     // 总数
	Page     int64  `json:"page"`      // 页数
	PageSize int64  `json:"page_size"` // 每页显示数
	OrderKey string `json:"order_key"` // 默认排序字段 -filed1,+field2,field3 (-Desc 降序)
}

// Set 设置钩子值
func (c *CollectionInfo) Set(key string, value interface{}) *CollectionInfo {
	c.client = c.client.Set(key, value)
	return c
}

// Create
// @Description: 添加记录
// @receiver c
// @param value
// @param gen 从分布式发号器获取ID
// @return error
func (c *CollectionInfo) Create(value interface{}, gen bool) error {
	return c.client.Set("gen", gen).Create(value).Error
}

// BatchCreate 批量创建
func (c *CollectionInfo) BatchCreate(values interface{}, batchSize int, gen bool) error {
	return c.client.Set("gen", gen).CreateInBatches(values, batchSize).Error
}

// Save 更新全部字段
func (c *CollectionInfo) Save(value interface{}) error {
	return c.client.Save(value).Error
}

// Updates 条件更新
func (c *CollectionInfo) Updates(where interface{}, value interface{}) error {
	return c.client.Model(where).Omit("id").Updates(value).Error
}

// Delete 模型删除
func (c *CollectionInfo) DeleteByModel(model interface{}) (count int64, err error) {
	db := c.client.Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// Delete 条件删除
func (c *CollectionInfo) DeleteByWhere(model, where interface{}) (int64, error) {
	db := c.client.Where(where).Delete(model)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// Delete 根据ID
func (c *CollectionInfo) DeleteByID(model interface{}, id int64) (count int64, err error) {
	db := c.client.Where("id=?", id).Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// Delete 根据ID批量删除
func (c *CollectionInfo) DeleteByIDS(model interface{}, ids []int64) (count int64, err error) {
	db := c.client.Where("id in (?)", ids).Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// First 根据ID查找
func (c *CollectionInfo) FirstByID(elect string, out interface{}, id int64) (err error) {
	if len(elect) > 0 {
		err = c.client.Select(elect).First(out, id).Error
	} else {
		err = c.client.First(out, id).Error
	}
	return
}

// 获取数量
func (c *CollectionInfo) Count(model, where interface{}) (count int64, err error) {
	err = c.client.Model(model).Where(where).Select("count(*)").Limit(1).Count(&count).Error
	return count, err
}

// 获取某个字段值
func (c *CollectionInfo) GetField(model, where interface{}, elect string) (out string, err error) {
	var info []string
	err = c.client.Model(model).Where(where).Limit(1).Pluck(elect, &info).Error
	if err == nil && len(info) == 1 {
		return info[0], err
	} else {
		return "", app.Err(app.Fail, "record not found")
	}
}

// First 查找单条记录
func (c *CollectionInfo) First(elect string, where interface{}, out interface{}) (err error) {
	if len(elect) > 0 {
		err = c.client.Select(elect).Where(where).First(out).Error
	} else {
		err = c.client.Where(where).First(out).Error
	}
	return
}

// Find 查找批量记录
func (c *CollectionInfo) Find(elect string, where interface{}, out interface{}, orders ...string) error {
	db := c.client.Where(where)
	if len(elect) > 0 {
		db = db.Select(elect)
	}
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// Scan
func (c *CollectionInfo) Scan(elect string, model, where interface{}, out interface{}) (err error) {
	db := c.client.Model(model)
	if len(elect) > 0 {
		db = db.Select(elect)
	}
	err = db.Where(where).Scan(out).Error
	return
}

// ScanList
func (c *CollectionInfo) ScanList(elect string, model, where interface{}, out interface{}, orders ...string) error {
	db := c.client.Model(model).Where(where)
	if len(elect) > 0 {
		db = db.Select(elect)
	}
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Scan(out).Error
}

// PluckList
func (c *CollectionInfo) PluckList(elect string, model, where interface{}, fieldName string, out interface{}) error {
	db := c.client.Model(model).Where(where)
	if len(elect) > 0 {
		db = db.Select(elect)
	}
	return db.Pluck(fieldName, out).Error
}
