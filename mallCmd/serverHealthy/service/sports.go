package service

import (
	"gitlab.mall.com/mallBase/basics/pkg/app"
	"gitlab.mall.com/mallBase/server/pkg/database/mongo"
	"gitlab.mall.com/mallBase/server/pkg/database/orm"
	"gitlab.mall.com/mallBase/server/pkg/pager"
	"go.mongodb.org/mongo-driver/bson"
	"serverHealthy/modelHealthy"
)

func SportsPage(name string, pageInfo orm.IndexPage) (list []modelHealthy.Sports, total int64, err error) {
	var where []pager.E
	if name != "" {
		where = append(where, pager.E{Key: "name", Value: name})
	}
	total, err = pager.New(pager.NewMongoDriver(), pageInfo).SetIndex("body_param").Where(where).Find(&list)
	return
}

func SportsSave(m modelHealthy.Sports) (res modelHealthy.Sports, err error) {
	if m.Id > 0 {
		if count, _ := mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Count(); count == 0 {
			return m, app.Err(app.Fail, "未找到记录")
		} else {
			_, err = mongo.Collection(&m).Where(bson.M{"_id": m.Id}).UpdateOne(&m)
			return m, err
		}
	} else {
		if count, _ := mongo.Collection(&res).Where(bson.M{"name": m.Name}).Count(); count == 0 {
			_, err = mongo.Collection(&res).InsertOne(m)
			return m, err
		} else {
			return m, app.Err(app.Fail, "存在重复记录")
		}
	}
}

func SportsList(id int64) (res modelHealthy.Sports, err error) {
	err = mongo.Collection(&res).Where(bson.M{"_id": id}).FindOne(&res)
	return res, err
}

func SportsDelete(m modelHealthy.Sports) (err error) {
	if count, _ := mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Count(); count == 0 {
		return app.Err(app.Fail, "未找到记录")
	} else {
		if _, err = mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Delete(); err != nil {
			return app.Err(app.Fail, "删除失败")
		}
	}
	return err
}
