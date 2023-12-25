package response

type ResFood struct {
	Id     int64    `json:"id,omitempty"`
	Name   string   `validate:"required" json:"name"`           //食物名称
	Heat   float64  `validate:"required" json:"heat,omitempty"` //热量 kcal/100g
	Energy float64  `json:"energy,omitempty"`                   //能量j/100g
	Pic    []string `json:"pic,omitempty"`                      // 图片 JSON数组对应图片ID 最多上传15张
}
