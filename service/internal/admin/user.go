package admin

import (
	"crypto/md5"
	"fmt"
	"gocms/model"
	"gocms/router"
	"gocms/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 用户
// 路由
type user struct {
	*AdminRouter
}

func NewUser(admin *AdminRouter) *user {
	return &user{AdminRouter: admin}
}

type userRes struct {
	*model.UserGroup
	CurrentRole model.RuleRoleGroup `json:"currentRole"`
}
type UserReq struct {
	*model.User
	Captcha string `json:"captcha"`
}

// 登录
func (this *user) Login(ctx *gin.Context) {
	userReq := UserReq{}
	if err := ctx.ShouldBindBodyWith(&userReq, binding.JSON); err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	var userRes userRes
	user, err := this.model.UserModel.GetUserByUsername(userReq.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "账号不存在",
		})
		return
	}
	this.log.Debug("captcha: ", zap.String("X-Captcha-ID", ctx.GetHeader("X-Captcha-ID")), zap.Any("TT", strings.ToUpper(userReq.Captcha)))
	res := this.cache.Get(ctx, "captcha-"+ctx.GetHeader("X-Captcha-ID"))
	if res.Err() != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "验证码错误",
		})
		return
	}
	if res.Val() != strings.ToUpper(userReq.Captcha) {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "验证码错误",
		})
		return
	}
	userRes.UserGroup = user
	if len(user.Role) > 0 {
		userRes.CurrentRole = user.Role[0]
	}
	if userRes.Password != fmt.Sprintf("%x", md5.Sum([]byte(userReq.Password))) {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "密码错误",
		})
		return
	}
	super := 0
	if userRes.User.ID == 1 {
		super = 1
	}
	token, _ := utils.GenerateToken(map[string]any{
		"user_id":  user.User.ID,
		"username": user.User.Username,
		"role_id":  userRes.CurrentRole.ID,
		"is_super": super,
	}, this.config)
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: gin.H{
			"accessToken": token,
			"user":        userRes,
		},
	})
}

// 获取用户信息
func (this *user) Detail(ctx *gin.Context) {
	username := ctx.GetString("username")
	var userRes userRes
	user, err := this.model.UserModel.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "账号不存在",
		})
		return
	}
	userRes.UserGroup = user
	if len(user.Role) > 0 {
		userRes.CurrentRole = user.Role[0]
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: userRes,
	})
}

// 修改用户信息
func (this *user) Update(ctx *gin.Context) {
	userReq := model.User{}
	if err := ctx.ShouldBindBodyWith(&userReq, binding.JSON); err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	userReq.Updatetime = time.Now()
	err := this.model.UserModel.UpdateUser(int64(userReq.ID), &userReq)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: userReq,
	})
}

// 增加用户
func (this *user) Add(ctx *gin.Context) {
	type UserReq struct {
		*model.UserGroup
		RoleIds []int64 `json:"roleIds"`
	}
	userReq := &UserReq{}
	if err := ctx.ShouldBindBodyWith(userReq, binding.JSON); err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	userReq.Createtime = time.Now()
	userReq.Updatetime = time.Now()
	userReq.Password = fmt.Sprintf("%x", md5.Sum([]byte(userReq.Password)))
	err := this.model.UserModel.CreateUser(userReq.UserGroup)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	for _, v := range userReq.RoleIds {
		this.e.AddGroupingPolicy(utils.ToString(userReq.User.ID), utils.ToString(v))
	}
	this.e.LoadPolicy()
	this.e.SavePolicy()
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: userReq,
	})
}

// 删除用户
func (this *user) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "id不能为空",
		})
	}
	iid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	err = this.model.UserModel.DeleteUser(int64(iid))
	this.model.RuleModel.Delete(&model.Rule{
		Ptype: "g",
		V0:    utils.ToString(iid),
	})
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
	})
}

// 获取用户列表
func (this *user) List(ctx *gin.Context) {
	users, err := this.model.UserModel.GetUserList()
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: users,
	})
}

// 用户权限
func (this *user) Permission(ctx *gin.Context) {
	id := ctx.GetFloat64("role_id")
	isSuper := ctx.GetInt("is_super")
	if isSuper == 1 {
		id = -1
	}
	plist, err := this.model.PermissionModel.GetPermissionsByRoleID(int64(id))
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: plist,
	})
}

// 获取验证码
func (this *user) Captcha(ctx *gin.Context) {
	svg, code := utils.GenerateSVG(80, 40)
	tu, _ := uuid.NewUUID()
	// 将验证码存储到 Redis
	err := this.cache.Set(ctx, "captcha-"+tu.String(), code, time.Minute*5)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 设置 Content-Type 为 "image/svg+xml"
	ctx.Header("Content-Type", "image/svg+xml; charset=utf-8")
	ctx.Header("X-Captcha-ID", tu.String())
	// 返回验证码
	ctx.Data(http.StatusOK, "image/svg+xml", svg)
}

// 刷新token
func (this *user) RefreshToken(ctx *gin.Context) {
	token, _ := utils.GenerateToken(map[string]any{
		"user_id":  ctx.GetFloat64("user_id"),
		"username": ctx.GetString("username"),
		"role_id":  ctx.GetFloat64("role_id"),
	}, this.config)
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: gin.H{
			"accessToken": token,
		},
	})
}

// ResetPassword
func (this *user) ResetPassword(ctx *gin.Context) {
	id := ctx.Param("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	userReq := model.UserProfileGroup{}
	if err := ctx.ShouldBindBodyWith(&userReq, binding.JSON); err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	user := model.User{
		ID:       int64(iid),
		Password: fmt.Sprintf("%x", md5.Sum([]byte(userReq.Password))),
	}
	err = this.model.UserModel.UpdateUser(int64(iid), &user)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
	})
}
