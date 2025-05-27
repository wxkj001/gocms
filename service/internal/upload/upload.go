package upload

import (
	"gocms/cache"
	"gocms/extend"
	"gocms/middleware"
	"gocms/model"
	"gocms/plugin"
	"gocms/router"
	"gocms/utils"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewUploadRouter(log *zap.Logger, config *viper.Viper, plugin *plugin.Plugins, middle middleware.MiddlewareParams, models model.ModelParams, extend extend.ExtendParams, e *casbin.Enforcer, cache cache.Cache) *UploadRouter {
	return &UploadRouter{log: log, config: config, middle: middle, model: models, e: e, cache: cache, plugins: plugin, extend: extend}
}

type UploadRouter struct {
	middle  middleware.MiddlewareParams
	log     *zap.Logger
	config  *viper.Viper
	model   model.ModelParams
	g       *gin.Engine
	e       *casbin.Enforcer
	cache   cache.Cache
	plugins *plugin.Plugins
	extend  extend.ExtendParams
}

// 注册路由
func (c *UploadRouter) RouteRegister(g *gin.Engine, r *gin.RouterGroup) {
	r.POST("/upload", c.middle.Jwt.AdminJWT(), c.Upload)
	r.GET("/upload/*filepath", c.GetFile)
}
func (c *UploadRouter) Upload(ctx *gin.Context) {
	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		c.log.Error("获取上传文件失败", zap.Error(err))
		ctx.JSON(http.StatusOK, router.Response{
			Code:    400,
			Message: "获取上传文件失败: " + err.Error(),
		})
		return
	}
	defer file.Close()
	sys, err := c.model.SysConfigModel.GetAllConfigMap()
	if err != nil {
		c.log.Error("获取系统配置失败", zap.Error(err))
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "获取系统配置失败",
		})
		return
	}
	// 获取文件大小
	fileSize := header.Size
	// 获取文件类型
	contentType := header.Header.Get("Content-Type")

	// 验证文件大小
	maxSize := int64(utils.ToInt(sys["max_size"]) * 1024 * 1024) // 默认单位为MB
	if maxSize > 0 && fileSize > maxSize {
		c.log.Error("文件大小超过限制", zap.Int64("size", fileSize), zap.Int64("maxSize", maxSize))
		ctx.JSON(http.StatusOK, router.Response{
			Code:    400,
			Message: "文件大小超过限制",
		})
		return
	}

	// 验证文件类型
	allowTypes := strings.Split(sys["allow_types"], ",") //c.config.GetStringSlice("upload.allow_types")
	if len(allowTypes) > 0 {
		valid := false
		for _, t := range allowTypes {
			if t == contentType {
				valid = true
				break
			}
		}
		if !valid {
			c.log.Error("文件类型不允许", zap.String("contentType", contentType))
			ctx.JSON(http.StatusOK, router.Response{
				Code:    400,
				Message: "文件类型不允许",
			})
			return
		}
	}

	// 生成唯一文件名
	// fileExt := filepath.Ext(header.Filename)
	uniqueName := header.Filename
	objectName := time.Now().Format("2006/01/02/") + uniqueName
	if utils.ToInt(sys["is_oss"]) == 1 {
		s3Config := c.extend.S3.GetConfig()
		s3Config.AccessKeyID = sys["AccessKeyID"]
		s3Config.SecretAccessKey = sys["SecretAccessKey"]
		s3Config.Endpoint = sys["Endpoint"]
		s3Config.BucketName = sys["BucketName"]
		s3Config.UseSSL = utils.ToInt(sys["UseSSL"]) == 1
		s3Config.Token = sys["Token"]
		s3 := c.extend.S3.Client(s3Config)
		if s3.Err() != nil {
			c.log.Error("初始化S3客户端失败", zap.Error(s3.Err()))
			ctx.JSON(http.StatusOK, router.Response{
				Code:    500,
				Message: "初始化S3客户端失败",
			})
			return
		}
		// 使用S3服务上传文件
		_, err := s3.UploadWithReader(s3Config.BucketName, objectName, file, fileSize, contentType)
		if err != nil {
			c.log.Error("上传文件到对象存储失败", zap.Error(err))
			ctx.JSON(http.StatusOK, router.Response{
				Code:    500,
				Message: "上传文件失败: " + err.Error(),
			})
			return
		}

	} else {
		prefix := utils.ToString(sys["upload_path"])
		if prefix != "" && !strings.HasSuffix(prefix, "/") {
			prefix = prefix + "/"
		}
		prefix += objectName
		err = ctx.SaveUploadedFile(header, objectName)
		if err != nil {
			c.log.Error("保存文件失败", zap.Error(err))
			ctx.JSON(http.StatusOK, router.Response{
				Code:    500,
				Message: "上传文件失败: " + err.Error(),
			})
			return
		}
	}
	host := ctx.Request.Host
	fileURL := host
	if !strings.HasSuffix(host, "/") && !strings.HasPrefix(objectName, "/") {
		fileURL += "/"
	}
	fileURL += "file/" + objectName
	c.model.MediasModel.Add(&model.Medias{
		MediaUrl:    objectName,
		MediaName:   uniqueName,
		ContentType: contentType,
		Size:        fileSize,
	})
	// 返回上传结果
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: gin.H{
			"url":  fileURL,
			"size": fileSize,
			"name": header.Filename,
			"type": contentType,
		},
	})
}
func (c *UploadRouter) GetFile(ctx *gin.Context) {
	filename := ctx.Param("filepath")
	if filename == "" {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    400,
			Message: "文件名不能为空",
		})
		return
	}
	c.log.Debug("获取文件", zap.String("filename", filename))
	sys, err := c.model.SysConfigModel.GetAllConfigMap()
	if err != nil {
		c.log.Error("获取系统配置失败", zap.Error(err))
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "获取系统配置失败",
		})
		return
	}
	if utils.ToInt(sys["is_oss"]) == 1 {
		s3Config := c.extend.S3.GetConfig()
		s3Config.AccessKeyID = sys["AccessKeyID"]
		s3Config.SecretAccessKey = sys["SecretAccessKey"]
		s3Config.Endpoint = sys["Endpoint"]
		s3Config.BucketName = sys["BucketName"]
		s3Config.UseSSL = utils.ToInt(sys["UseSSL"]) == 1
		s3Config.Token = sys["Token"]
		s3 := c.extend.S3.Client(s3Config)
		if s3.Err() != nil {
			c.log.Error("初始化S3客户端失败", zap.Error(s3.Err()))
			ctx.JSON(http.StatusOK, router.Response{
				Code:    500,
				Message: "初始化S3客户端失败",
			})
			return
		}
		// 移除filename里第一个/
		filename = strings.TrimLeft(filename, "/")
		// 使用S3服务获取文件
		object, err := s3.GetObjectURL(s3Config.BucketName, filename, 3600)
		if err != nil {
			c.log.Error("获取文件失败", zap.Error(err))
			ctx.JSON(http.StatusOK, router.Response{
				Code:    500,
				Message: "获取文件失败: " + err.Error(),
			})
			return
		}
		ctx.Redirect(http.StatusMovedPermanently, object)
	} else {
		prefix := utils.ToString(sys["upload_path"])
		if prefix != "" && !strings.HasSuffix(prefix, "/") {
			prefix = prefix + "/"
		}
		prefix += filename
		ctx.File(prefix)
	}
}
