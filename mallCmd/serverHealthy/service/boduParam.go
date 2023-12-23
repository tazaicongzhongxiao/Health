package service

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/server/pkg/database/mongo"
	"MyTestMall/mallBase/server/pkg/database/orm"
	"MyTestMall/mallBase/server/pkg/pager"
	"go.mongodb.org/mongo-driver/bson"
	"serverHealthy/modelHealthy"
)

func BodyParamPage(bmi float32, pageInfo orm.IndexPage) (list []modelHealthy.BodyParam, total int64, err error) {
	var where []pager.E
	if bmi != 0 {
		where = append(where, pager.E{Key: "bmi", Value: bmi})
	}
	total, err = pager.New(pager.NewMongoDriver(), pageInfo).SetIndex("body_param").Where(where).Find(&list)
	return
}

func BodyParamSave(m modelHealthy.BodyParam) (res modelHealthy.BodyParam, err error) {
	if m.Id > 0 {
		if count, _ := mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Count(); count == 0 {
			return m, app.Err(app.Fail, "未找到记录")
		} else {
			_, err = mongo.Collection(&m).Where(bson.M{"_id": m.Id}).UpdateOne(&m)
			return m, err
		}
	} else {
		if count, _ := mongo.Collection(&res).Where(bson.M{"bmi": m.BMI}).Count(); count == 0 {
			_, err = mongo.Collection(&res).InsertOne(m)
			return m, err
		} else {
			return m, app.Err(app.Fail, "存在重复记录")
		}
	}
}

func BodyParamList(id int64) (res modelHealthy.BodyParam, err error) {
	err = mongo.Collection(&res).Where(bson.M{"_id": id}).FindOne(&res)
	return res, err
}

func BodyParamDelete(m modelHealthy.BodyParam) (err error) {
	if count, _ := mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Count(); count == 0 {
		return app.Err(app.Fail, "未找到记录")
	} else {
		if _, err = mongo.Collection(&m).Where(bson.M{"_id": m.Id}).Delete(); err != nil {
			return app.Err(app.Fail, "删除失败")
		}
	}
	return err
}
