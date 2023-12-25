package request

type PageInfo struct {
	Page     int32  `validate:"required" json:"page" form:"page"`           // 当前页
	PageSize int32  `validate:"required" json:"page_size" form:"page_size"` // 每页显示数
	OrderKey string `json:"order_key" form:"orderKey"`                      // 默认排序字段 -filed1,+field2,field3 (-Desc 降序)
}

type ReqId struct {
	Id int64 `validate:"required" json:"id" form:"id"`
}
