package controllers

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/client/pkg/comm"
	"MyTestMall/mallBase/client/pkg/validator"
	protoHealthy "MyTestMall/protoImp/healthy"
	"clientHealthy/client"
	"clientHealthy/request"
	"github.com/gin-gonic/gin"
)

func Sports(ctx *gin.Context) {
	var form request.PageInfo
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	var Form protoHealthy.ReqSportsPage
	if err := app.UnmarshalJson(form, &Form.Page); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.SportsPage(Form)
	comm.ApiResponse(ctx, &res, err)
}

func SportsList(ctx *gin.Context) {
	var form request.ReqId
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	res, err := client.GRPC.ServerHealthy.SportsList(protoHealthy.ReqId{
		Id: form.Id,
	})
	comm.ApiResponse(ctx, &res, err)
}

func SportsSave(ctx *gin.Context) {
	var form request.ReqSports
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	var Form protoHealthy.Sports
	if err := app.UnmarshalJson(form, &Form); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.SportsSave(Form)
	comm.ApiResponse(ctx, &res, err)
}

func SportsDelete(ctx *gin.Context) {
	var form request.ReqSports
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	var Form protoHealthy.Sports
	if err := app.UnmarshalJson(form, &Form); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.SportsDelete(Form)
	comm.ApiResponse(ctx, &res, err)
}
