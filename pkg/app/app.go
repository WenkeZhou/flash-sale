package app

import (
	"fmt"
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	fmt.Printf("StatusCode:%v, Data: %v \n", http.StatusOK, data)
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	response := gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	}
	fmt.Printf("StatusCode:%v, Data: %v \n", http.StatusOK, response)
	r.Ctx.JSON(http.StatusOK, response)
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	fmt.Printf("StatusCode:%v, Data: %v \n", err.StatusCode(), response)
	r.Ctx.JSON(err.StatusCode(), response)
}
