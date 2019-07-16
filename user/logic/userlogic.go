package logic

import (
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/yakaa/log4g"
	"integral-mall/common/baseerror"
	"integral-mall/user/model"
	"strconv"
)

type (
	UserLogic struct {
		userModel        *model.UserModel
		redisCache       *redis.Client
	}
	RegisterRequest struct {
		Mobile   string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	RegisterResponse struct {
	}

	LoginRequest struct {
		Mobile   string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	LoginResponse struct {
		Authorization string `json:"authorization"`
	}
)

var (
	ErrRecordExist        = baseerror.NewBaseError("此手机号已经存在")
	ErrUserNameOrPassword = baseerror.NewBaseError("用户名或密码错误")
)

func NewUserLogic(userModel *model.UserModel,redisCache *redis.Client) *UserLogic {
	return &UserLogic{userModel:userModel,redisCache:redisCache}
}

/**
   注册操作
 */
func (l *UserLogic) Register(r *RegisterRequest) (*RegisterResponse,error) {
	response :=new(RegisterResponse)
	//1、所添加的手机号内容是否存在
	b,err:=l.userModel.ExistByMobile(r.Mobile)
	if err != nil {
		return nil, err
	}
	if b {
		return nil, ErrRecordExist
	}


	//添加操作
	if _,err :=l.userModel.Insert(&model.User{
		Mobile:r.Mobile,
		Password:fmt.Sprintf("%x", md5.Sum([]byte(r.Password))),
	});err !=nil{
		return nil,err
	}

	return response,nil
}

/**
   用户登录
 */
func (l *UserLogic) Login(r *LoginRequest) (*LoginResponse,error) {
	response := new(LoginResponse)
	user, err := l.userModel.FindByMobile(r.Mobile)
	if err != nil {
		return nil, err
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(r.Password))) {
		return nil, ErrUserNameOrPassword
	}
	response.Authorization = fmt.Sprintf("%x", md5.Sum([]byte(user.Mobile+strconv.Itoa(int(user.Id)))))
	sts:=l.redisCache.Set(response.Authorization, user.Id,0)
	 log4g.Info(sts.Result())
	return response, nil
}