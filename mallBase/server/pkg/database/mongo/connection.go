package mongo

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/log"
	baseMongo "MyTestMall/mallBase/basics/pkg/mongo"
	"MyTestMall/mallBase/basics/pkg/unique"
	"MyTestMall/mallBase/basics/tools/dinterface"
	"MyTestMall/mallBase/server/pkg/database"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type (
	// CollectionInfo 集合包含的连接信息和查询等操作信息
	CollectionInfo struct {
		Db       *mongo.Client
		Database *mongo.Database
		Table    *mongo.Collection
		filter   interface{}
		limit    int64
		skip     int64
		sort     interface{}
		fields   bson.M
	}
)

func Connection(c configs) (db *mongo.Client, err error) {
	mongoOptions := options.Client()
	if c.MaxConnIdleTime > 0 {
		mongoOptions.SetMaxConnIdleTime(time.Duration(c.MaxConnIdleTime) * time.Second)
	}
	if c.MaxPoolSize > 0 {
		mongoOptions.SetMaxPoolSize(uint64(c.MaxPoolSize))
	}
	if c.Username != "" && c.Password != "" {
		mongoOptions.SetAuth(options.Credential{AuthSource: c.AuthSource, Username: c.Username, Password: c.Password})
	}
	if c.ReplSetName != "" {
		mongoOptions.ReplicaSet = &c.ReplSetName
	}
	if c.Direct {
		mongoOptions.Direct = &c.Direct
	}
	db, err = mongo.NewClient(mongoOptions.ApplyURI("mongodb://" + c.URL))
	if err != nil {
		log.Error("mongo connect error, you can't use orm support", err.Error())
		return db, err
	} else {
		log.Debug(fmt.Sprintf("mongo connect: %s %s", c.URL, c.Database), nil)
	}
	if err = db.Connect(context.Background()); err != nil {
		return db, err
	}
	if err = db.Ping(context.Background(), nil); err != nil {
		log.Error("mongo Ping异常", err.Error())
		return db, err
	}
	return db, err
}

// GetIndexModel
// @Description: 转换为索引格式
// @param data
// @return []mongo.IndexModel
func GetIndexModel(data []baseMongo.IndexData) []mongo.IndexModel {
	var index []mongo.IndexModel
	for _, val := range data {
		var opt = options.IndexOptions{}
		opt.SetBackground(true) // 后台构建索引
		if val.Weights > 0 {
			opt.SetWeights(val.Weights)
		}
		if val.Unique {
			opt.SetUnique(val.Unique)
		}
		if val.Exp > 0 {
			opt.SetExpireAfterSeconds(val.Exp) //设置过期时间
		}
		Keys := bson.D{}
		for _, k := range val.Keys {
			if k.Key == "1" || k.Key == "-1" {
				if Key, err := strconv.Atoi(k.Key); err == nil {
					Keys = append(Keys, bson.E{Key: k.Name, Value: Key})
				}
			} else {
				Keys = append(Keys, bson.E{Key: k.Name, Value: k.Key})
			}
		}
		index = append(index, mongo.IndexModel{Keys: Keys, Options: &opt})
	}
	return index
}

// GetDataBase
// @Description: 获取数据库连接
// @param db
// @param database
// @return *CollectionInfo
func GetDataBase(db *mongo.Client, database string) *CollectionInfo {
	collection := &CollectionInfo{
		Db:       db,
		Database: db.Database(database),
		filter:   bson.M{},
	}
	return collection
}

// GetCollection
// @Description: 得到一个mongo操作对象
// @param db
// @param database
// @param table
// @return *CollectionInfo
func GetCollection(db *mongo.Client, database string, table database.Table) *CollectionInfo {
	dbTmp := db.Database(database)
	return &CollectionInfo{
		Db:       db,
		Database: dbTmp,
		Table:    dbTmp.Collection(table.TableName()),
		filter:   bson.M{},
	}
}

// IsRecordNotFoundError
// @Description: 检查获取当前错误是否为未找到数据
// @param err
// @return bool
func IsRecordNotFoundError(err error) bool {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true
		}
	}
	return false
}

// 转换

// Collection 得到一个mongo操作对象
func (collection *CollectionInfo) Collection(table database.Table) *CollectionInfo {
	collection.Table = collection.Database.Collection(table.TableName())
	collection.filter = nil
	collection.limit = 0
	collection.skip = 0
	collection.sort = nil
	collection.fields = bson.M{}
	return collection
}

// SetTable 设置集合名称
func (collection *CollectionInfo) SetTable(name string) *CollectionInfo {
	collection.Table = collection.Database.Collection(name)
	return collection
}

// Where 条件查询, bson.M{"field": "value"} 、 bson.D{{"field", "value"}}
func (collection *CollectionInfo) Where(m interface{}) *CollectionInfo {
	collection.filter = m
	return collection
}

// Limit 限制条数
func (collection *CollectionInfo) Limit(n int64) *CollectionInfo {
	if n > 0 {
		collection.limit = n
	}
	return collection
}

// Skip 跳过条数
func (collection *CollectionInfo) Skip(n int64) *CollectionInfo {
	collection.skip = n
	return collection
}

// Sort 排序 bson.D{{"created_at",-1}}
func (collection *CollectionInfo) Sort(sorts bson.D) *CollectionInfo {
	collection.sort = sorts
	return collection
}

// Fields 指定查询字段  field1,field2
func (collection *CollectionInfo) Fields(fields string) *CollectionInfo {
	if len(fields) > 0 {
		sorts := make(bson.M)
		for _, v := range strings.Split(fields, ",") {
			sorts[strings.TrimSpace(v)] = 1
		}
		collection.fields = sorts
	}
	return collection
}

func (collection *CollectionInfo) SetFields(fields bson.M) *CollectionInfo {
	collection.fields = fields
	return collection
}

// StartSession
// @Description: 打开 MONGO 事务
// @param db
// @return sess
// @return ctx
// @return err
// mongodb.StartSession(mongodb.Client())
// err = session.AbortTransaction(sc)
// err = session.CommitTransaction(sc)
// session.EndSession(ctx)
func StartSession(db *mongo.Client) (sess mongo.Session, ctx context.Context, err error) {
	if sess, err = db.StartSession(); err != nil {
		return sess, ctx, err
	}
	if err = sess.StartTransaction(); err != nil {
		return sess, ctx, err
	}
	return sess, context.TODO(), err
}

// CreateMany
// @Description: 批量创建索引
// @receiver collection
// @param req
// @return res
// @return err
//
// mongo.Collection(&req).CreateMany([]mongo.IndexData{{Name: "uid", Key: "1", Unique: true}, {Name: "created_at", Key: "-1"}})
func (collection *CollectionInfo) CreateMany(req []baseMongo.IndexData) (res []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err = collection.Table.Indexes().CreateMany(ctx, GetIndexModel(req), options.CreateIndexes().SetMaxTime(10*time.Second))
	return res, err
}

// InsertOne
// @Description: 写入单条数据
// @receiver collection
// @param document
// @return res
// @return err
func (collection *CollectionInfo) InsertOne(document interface{}) (res *mongo.InsertOneResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err = collection.Table.InsertOne(ctx, BeforeCreate(document))
	return res, err
}

// InsertMany
// @Description: 写入多条数据
// @receiver collection
// @param documents
// @return res
// @return err
func (collection *CollectionInfo) InsertMany(documents interface{}) (res *mongo.InsertManyResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var data []interface{}
	data = BeforeCreate(documents).([]interface{})
	res, err = collection.Table.InsertMany(ctx, data)
	return res, err
}

// UpdateOrInsert
// @Description: 存在更新,不存在写入
// @receiver collection
// @param document
// @return res
// @return err
func (collection *CollectionInfo) UpdateOrInsert(document interface{}) (res *mongo.UpdateResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err = collection.Table.UpdateOne(ctx, collection.filter, bson.D{
		{"$set", BeforeUpdate(document)},
		{"$setOnInsert", bson.M{"created_at": time.Now().Unix()}},
	}, options.Update().SetUpsert(true))
	return res, err
}

// UpdateOne
// @Description: 更新一条
// @receiver collection
// @param document
// @return res
// @return err
func (collection *CollectionInfo) UpdateOne(document interface{}) (res *mongo.UpdateResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err = collection.Table.UpdateOne(ctx, collection.filter, bson.D{
		{"$set", BeforeUpdate(document)},
	})
	return res, err
}

// UpdateMany
// @Description: 更新多条
// @receiver collection
// @param document
// @return res
// @return err
func (collection *CollectionInfo) UpdateMany(document interface{}) (res *mongo.UpdateResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err = collection.Table.UpdateMany(ctx, collection.filter, bson.D{{"$set", BeforeUpdate(document)}})
	return res, err
}

// FindOne
// @Description: 查询一条数据
// @receiver collection
// @param document
// @return error
func (collection *CollectionInfo) FindOne(document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.Table.FindOne(ctx, collection.filter, &options.FindOneOptions{
		Skip:       &collection.skip,
		Sort:       collection.sort,
		Projection: collection.fields,
	}).Decode(document)
	return err
}

// FindOneAndUpdate
// @Description: 查询单条数据后修改该数据
// @receiver collection
// @param document
// @return error
func (collection *CollectionInfo) FindOneAndUpdate(document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.Table.FindOneAndUpdate(ctx, collection.filter, bson.D{{"$set", BeforeUpdate(document)}}, &options.FindOneAndUpdateOptions{
		Sort:       collection.sort,
		Projection: collection.fields,
	}).Decode(document)
	return err
}

// FindOneAndDelete
// @Description: 查询单条数据后删除该数据
// @receiver collection
// @param document
// @return error
func (collection *CollectionInfo) FindOneAndDelete(document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.Table.FindOneAndDelete(ctx, collection.filter, &options.FindOneAndDeleteOptions{
		Sort:       collection.sort,
		Projection: collection.fields,
	}).Decode(document)
	return err
}

// FindMany
// @Description: 查询多条数据
// @receiver collection
// @param documents
// @return error
func (collection *CollectionInfo) FindMany(documents interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.Table.Find(ctx, collection.filter, &options.FindOptions{
		Skip:       &collection.skip,
		Limit:      &collection.limit,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	if err != nil {
		return err
	}
	defer result.Close(ctx)
	err = result.All(ctx, documents)
	return nil
}

// AggregateMany
// @Description: 聚合查询
// @receiver collection
// @param documents
// @param pipeline
// @return error
//
// conditions := []bson.M{{"$unwind": "$test"}, {"$match": bson.M{"test.id": id}}, {"$project": bson.M{"test": 1}}}
// err = mongo.Collection(&mongoSql.CoUserStaffGroup{}).AggregateMany(&tmp, conditions)
func (collection *CollectionInfo) AggregateMany(documents interface{}, pipeline interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.Table.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer result.Close(ctx)
	err = result.All(ctx, documents)
	return err
}

// Delete
// @Description: 删除数据,并返回删除成功的数量
// @receiver collection
// @return int64
// @return error
func (collection *CollectionInfo) Delete() (int64, error) {
	if collection.filter == nil || reflect.ValueOf(collection.filter).Len() == 0 {
		return 0, app.Err(app.Fail, "you can't delete all documents, it's very dangerous")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.Table.DeleteMany(ctx, collection.filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, err
}

// SoftDeleteOne
// @Description: 软删除数据,
// @receiver collection
// @return int64
// @return error
func (collection *CollectionInfo) SoftDeleteOne() (int64, error) {
	if collection.filter == nil || reflect.ValueOf(collection.filter).Len() == 0 {
		return 0, app.Err(app.Fail, "you can't delete all documents, it's very dangerous")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var document bson.D
	err := collection.Table.FindOne(ctx, collection.filter).Decode(&document)
	if err == nil {
		document = append(document, bson.E{Key: "deleted_at", Value: time.Now().UTC().Unix()})
		_, err = collection.Database.Collection(collection.Table.Name()+"_del").InsertOne(ctx, &document)
	}
	if err != nil {
		return 0, err
	}
	result, err := collection.Table.DeleteOne(ctx, collection.filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, err
}

// Count
// @Description: 根据指定条件获取总条数
// @receiver collection
// @return int64
// @return error
func (collection *CollectionInfo) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.Table.CountDocuments(ctx, collection.filter)
	if err != nil {
		return 0, err
	}
	return count, err
}

// BeforeCreate
// @Description: 创建数据前置操作
// @param document
// @param omitempty
// @return interface{}
func BeforeCreate(document interface{}, omitempty ...bool) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)
	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeCreate(val.Elem().Interface(), omitempty...)
	case reflect.Array, reflect.Slice:
		if val.Type() == reflect.TypeOf(bson.D{}) {
			return document
		} else {
			var sliceData = make([]interface{}, val.Len(), val.Cap())
			for i := 0; i < val.Len(); i++ {
				sliceData[i] = BeforeCreate(val.Index(i).Interface(), omitempty...)
			}
			return sliceData
		}
	case reflect.Struct:
		data := bson.D{}
		for i := 0; i < typ.NumField(); i++ {
			tag := strings.Split(typ.Field(i).Tag.Get("bson"), ",")[0]
			if tag != "-" && tag != "created_at" && tag != "updated_at" {
				if tag == "_id" {
					// ID 主键
					id := val.Field(i).Interface()
					if val.Field(i).Type() == reflect.TypeOf(primitive.ObjectID{}) && id.(primitive.ObjectID).IsZero() {
						id = primitive.NewObjectID()
					} else if IsIntn(val.Field(i).Kind()) && dinterface.IsNil(id) {
						id = unique.ID()
					} else if val.Field(i).Kind() == reflect.String && dinterface.IsNil(id) {
						id = primitive.NewObjectID().Hex()
					}
					data = append(data, bson.E{Key: tag, Value: id})
				} else {
					if len(omitempty) == 0 || !isZero(val.Field(i)) {
						data = append(data, bson.E{Key: tag, Value: val.Field(i).Interface()})
					}
				}
			}
		}
		now := time.Now().Unix()
		data = append(data, bson.E{Key: "created_at", Value: now})
		data = append(data, bson.E{Key: "updated_at", Value: now})
		return data
	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			if !val.MapIndex(reflect.ValueOf("_id")).IsValid() {
				val.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(unique.ID()))
			}
			now := time.Now().Unix()
			val.SetMapIndex(reflect.ValueOf("created_at"), reflect.ValueOf(now))
			val.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(now))
		}
		return val.Interface()
	}
}

func StructToBson(document interface{}) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)
	switch typ.Kind() {
	case reflect.Ptr:
		return StructToBson(val.Elem().Interface())
	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			tag := strings.Split(typ.Field(i).Tag.Get("bson"), ",")[0]
			if tag != "_id" && tag != "-" && tag != "created_at" && tag != "updated_at" {
				data[tag] = val.Field(i).Interface()
			}
		}
		return data
	}
	return nil
}

func WhereToBson(document interface{}) (result bson.D) {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)
	for i := 0; i < val.NumField(); i++ {
		filed := val.Field(i)
		tag := strings.Split(typ.Field(i).Tag.Get("bson"), ",")[0]
		if tag != "-" {
			switch filed.Kind() {
			case reflect.String:
				if v := filed.String(); v != "" {
					result = append(result, bson.E{Key: tag, Value: v})
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v := filed.Int(); v != 0 {
					result = append(result, bson.E{Key: tag, Value: v})
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				if v := filed.Uint(); v != 0 {
					result = append(result, bson.E{Key: tag, Value: v})
				}
			case reflect.Float32, reflect.Float64:
				if v := filed.Float(); v != 0 {
					result = append(result, bson.E{Key: tag, Value: v})
				}
			case reflect.Bool:
				if v := filed.Bool(); v == true {
					result = append(result, bson.E{Key: tag, Value: v})
				}
			case reflect.Slice:
				if filed.Len() > 0 {
					result = append(result, bson.E{Key: tag, Value: bson.D{{"$in", filed.Interface()}}})
				}
			}
		}
	}
	return result
}

// BeforeUpdate 更新数据前置操作
// omitempty 空值不更新
func BeforeUpdate(document interface{}, omitempty ...bool) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)
	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeUpdate(val.Elem().Interface(), omitempty...)
	case reflect.Slice, reflect.Array:
		if val.Type() == reflect.TypeOf(bson.D{}) {
			data := document.(bson.D)
			data = append(data, bson.E{Key: "updated_at", Value: time.Now().Unix()})
			return data
		} else {
			var sliceData = make([]interface{}, val.Len(), val.Cap())
			for i := 0; i < val.Len(); i++ {
				sliceData[i] = BeforeUpdate(val.Index(i).Interface(), omitempty...).(bson.M)
			}
			return sliceData
		}
	case reflect.Struct:
		data := bson.D{}
		for i := 0; i < typ.NumField(); i++ {
			if len(omitempty) > 0 || !isZero(val.Field(i)) {
				tag := strings.Split(typ.Field(i).Tag.Get("bson"), ",")[0]
				if tag != "_id" && tag != "-" && tag != "created_at" && tag != "updated_at" {
					data = append(data, bson.E{Key: tag, Value: val.Field(i).Interface()})
				}
			}
		}
		data = append(data, bson.E{Key: "updated_at", Value: time.Now().Unix()})
		return data
	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			val.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(time.Now().Unix()))
		}
		return val.Interface()
	}
}

// IsIntn
// @Description: 是否为整数
// @param p
// @return bool
func IsIntn(p reflect.Kind) bool {
	return p == reflect.Int || p == reflect.Int64 || p == reflect.Uint64 || p == reflect.Int32 || p == reflect.Uint32
}

func isZero(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// Struct2Json
// @Description:
// @param form
// @return []byte
func Struct2Json(form interface{}) []byte {
	bsonx, _ := bson.Marshal(form)
	return bsonx
}

// Unmarshal
// @Description: 读取值转换
// @param req
// @param v
// @return error
func Unmarshal(req interface{}, v interface{}) error {
	bsonx, err := bson.Marshal(req)
	if err == nil {
		err = bson.Unmarshal(bsonx, v)
	}
	return err
}
