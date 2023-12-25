package controllers

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/client/pkg/comm"
	"MyTestMall/mallBase/client/pkg/validator"
	protoHealthy "MyTestMall/protoImp/healthy"
	"clientHealthy/client"
	"clientHealthy/request"
	"clientHealthy/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Food(ctx *gin.Context) {
	var form request.PageInfo
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	var Form protoHealthy.ReqFoodPage
	if err := app.UnmarshalJson(form, &Form.Page); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.FoodPage(Form)
	comm.ApiResponse(ctx, &res, err)
}

func FoodList(ctx *gin.Context) {
	var form request.ReqId
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	res, _ := client.GRPC.ServerHealthy.FoodList(protoHealthy.ReqId{
		Id: form.Id,
	})
	var result response.ResFood
	_ = app.Unmarshal(res, &result)
	for _, v := range result.Pic {
		v = fmt.Sprintf("http://192.168.2.139:31000/Pic/%s/%s.png", result.Name, v)
		result.Pic = append(result.Pic, v)
	}
	app.Println(3, result)
	ctx.JSON(http.StatusOK, &result)
}

func CreatPic(name string) (err error) {
	err = os.MkdirAll(fmt.Sprintf("./pic/%s", name), os.ModePerm)
	if err != nil {
		app.Println(3, err)
		return app.Err(app.Fail, "Error opening or creating file")
	}
	return nil
}

func FoodSave(ctx *gin.Context) {
	var form request.ReqFood
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	_ = CreatPic(form.Name)
	var Form protoHealthy.Food
	if err := app.UnmarshalJson(form, &Form); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.FoodSave(Form)
	comm.ApiResponse(ctx, &res, err)
}

func FoodDelete(ctx *gin.Context) {
	var form request.ReqFood
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	var Form protoHealthy.Food
	if err := app.UnmarshalJson(form, &Form); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.FoodDelete(Form)
	comm.ApiResponse(ctx, &res, err)
}
