package modelHealthy

import "image"

type BodyParam struct {
	Id        int64   `bson:"_id" json:"id,omitempty"`        //身体参数对照表ID
	High      float32 `bson:"high" json:"high,omitempty"`     //身高cm
	Weight    float32 `bson:"weight" json:"weight,omitempty"` //体重kg
	BMI       float32 `bson:"bmi" json:"bmi,omitempty"`       //体重指数
	CreatedAt int64   `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt int64   `bson:"updated_at" json:"updated_at,omitempty"`
}

func (m *BodyParam) TableName() string {
	return "body_param"
}

type Food struct {
	Id        int64    `bson:"_id" json:"id,omitempty"`
	Name      string   `bson:"name" json:"name,omitempty"`     //食物名称
	Heat      float64  `bson:"heat" json:"heat,omitempty"`     //热量 kcal/100g
	Energy    float64  `bson:"energy" json:"energy,omitempty"` //能量j/100g
	Pic       []string `bson:"pic" json:"pic,omitempty"`       // 图片 JSON数组对应图片ID 最多上传15张
	CreatedAt int64    `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt int64    `bson:"updated_at" json:"updated_at,omitempty"`
}

func (m *Food) TableName() string {
	return "food"
}

type Sports struct {
	Id        int64    `bson:"_id" json:"id,omitempty"`
	Name      string   `bson:"name" json:"name,omitempty"`          //运动名称
	Equipment []string `bson:"equipment"json:"equipment,omitempty"` //所需器材
	Hour      int32    `bson:"hour" json:"hour,omitempty"`          //运动时长
	CreatedAt int64    `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt int64    `bson:"updated_at" json:"updated_at,omitempty"`
}

func (m *Sports) TableName() string {
	return "sports"
}

type Pic struct {
	Id        int64       `bson:"_id" json:"id,omitempty"`
	Name      string      `bson:"name" json:"name,omitempty"`           //所属种类
	FilePath  string      `bson:"file_path" json:"file_path,omitempty"` //路径
	Size      image.Point `bson:"size" json:"size,omitempty"`           //属性
	CreatedAt int64       `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt int64       `bson:"update_at" json:"update_at"`
}

func (m *Pic) TableName() string {
	return "pic"
}
