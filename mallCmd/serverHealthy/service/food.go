package service

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/server/pkg/database/mongo"
	"MyTestMall/mallBase/server/pkg/database/orm"
	"MyTestMall/mallBase/server/pkg/pager"
	"go.mongodb.org/mongo-driver/bson"
	"serverHealthy/modelHealthy"
)

func FoodPage(name string, pageInfo orm.IndexPage) (list []modelHealthy.Food, total int64, err error) {
	var where []pager.E
	if name != "" {
		where = append(where, pager.E{Key: "name", Value: name})
	}
	total, err = pager.New(pager.NewMongoDriver(), pageInfo).SetIndex("body_param").Where(where).Find(&list)
	return
}

func FoodSave(m modelHealthy.Food) (res modelHealthy.Food, err error) {
	if m.Id > 0 {
		if count, _ := mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Count(); count == 0 {
			return m, app.Err(app.Fail, "未找到记录")
		} else {
			m.Pic, err = DownLoadPic(m.Name)
			if err != nil {
				return m, app.Err(app.Fail, "图像录入失败")
			}
			_, _ = mongo.Collection(&m).Where(bson.M{"_id": m.Id}).UpdateOne(&m)
		}
	} else {
		if count, _ := mongo.Collection(&res).Where(bson.M{"name": m.Name}).Count(); count == 0 {
			if m.Heat > 0 {
				m.Energy = 4185.85 * m.Heat
				_, _ = mongo.Collection(&res).InsertOne(m)
			} else {
				return m, app.Err(app.Fail, "食物热量必须存在")
			}
		} else {
			return m, app.Err(app.Fail, "存在重复记录")
		}
	}
	return m, err
}

func FoodList(id int64) (res modelHealthy.Food, err error) {
	err = mongo.Collection(&res).Where(bson.M{"_id": id}).FindOne(&res)
	return res, err
}

func FoodDelete(m modelHealthy.Food) (err error) {
	if count, _ := mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Count(); count == 0 {
		return app.Err(app.Fail, "未找到记录")
	} else {
		if _, err = mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Delete(); err != nil {
			return app.Err(app.Fail, "删除失败")
		}
	}
	return err
}
