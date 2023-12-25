package request

type ReqBodyParam struct {
	Id     int64   `json:"id"`                         //身体参数对照表ID
	High   float32 `validate:"required" json:"high"`   //身高cm
	Weight float32 `validate:"required" json:"weight"` //体重kg
	BMI    float32 ` json:"bmi"`                       //体重指数
}

type ReqFood struct {
	Id     int64    `json:"id,omitempty"`
	Name   string   `validate:"required" json:"name"`           //食物名称
	Heat   float64  `validate:"required" json:"heat,omitempty"` //热量 kcal/100g
	Energy float64  `json:"energy,omitempty"`                   //能量j/100g
	Pic    []string `json:"pic,omitempty"`                      // 图片 JSON数组对应图片ID 最多上传15张
}

type ReqSports struct {
	Id        int64    `json:"id,omitempty"`
	Name      string   `validate:"required" json:"name,omitempty"`      //运动名称
	Equipment []string `validate:"required" json:"equipment,omitempty"` //所需器材
	Hour      int32    `validate:"required" json:"hour,omitempty"`      //运动时长
}
