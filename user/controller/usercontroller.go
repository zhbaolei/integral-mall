package controller

import (
	"github.com/gin-gonic/gin"
	"integral-mall/common/baseresponse"
	"integral-mall/user/logic"
)

type (
	UserController struct {
		userLogic *logic.UserLogic
	}
)

func NewUserController(userLogic *logic.UserLogic) *UserController {
	return &UserController{userLogic:userLogic}
}

/**
  注册方法
 */
func (c *UserController) Register(ctx *gin.Context) {
	r := new(logic.RegisterRequest)
	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}

	res, err := c.userLogic.Register(r)
	baseresponse.HttpResponse(ctx, res, err)
	return
}

/**
  登录方法
 */
func (l *UserController) Login(ctx *gin.Context) {
	r := new(logic.LoginRequest)

	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}

	//调用登录方法
	res, err := l.userLogic.Login(r)
	baseresponse.HttpResponse(ctx, res, err)
	return
}