package api

import (
	"cloud-batch/api/docs"
	"cloud-batch/api/middleware"
	"cloud-batch/api/v1"
	"cloud-batch/configs"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/paradeum-team/gin-prometheus-ext"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"strings"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	if configs.Server.IsOpenMetrics {
		p := ginprometheusext.NewPrometheus("gin")
		p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
			url := c.Request.URL.Path
			paramsKeys := map[string]string{"dgst": ":dgst", "id": ":id", "parentid": ":parentid"}
			for _, p := range c.Params {
				v, ok := paramsKeys[p.Key]
				if ok {
					url = strings.Replace(url, p.Value, v, 1)
					break
				}
			}
			return url
		}
		metricsRouter := gin.New()
		//metricsRouter.GET("/env", api.GetEnv)
		p.SetListenAddressWithRouter(configs.Server.MetricsAddr, metricsRouter)
		p.Use(r)
		pprof.Register(metricsRouter)
	}

	r.Use(gin.Recovery())
	r.Use(requestid.New())
	r.Use(middleware.LoggerToFile(configs.LogConfig.AppName), gin.Recovery())
	r.Use(middleware.Cors(), gin.Recovery())

	//r.StaticFS("/upload/files", http.Dir(upload.GetMediaFullPath()))

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Cloud Batch API"
	docs.SwaggerInfo.Description = "Cloud Batch"
	docs.SwaggerInfo.Version = configs.Server.Version
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/version", v1.GetVersion)
		auth := apiv1.Group("auth/")
		{
			// 登录生成token
			auth.POST("/login", v1.Login)
			auth.Use(middleware.JWT())
			auth.PUT("/passwd", v1.UpdateAuth)
		}

		apiv1.GET("servers", v1.BatchGetServers)

		batch := apiv1.Group("batch/")
		batch.Use(middleware.JWT())
		{
			batch.POST("servers", v1.BatchCreateServers)
			batch.DELETE("servers", v1.BatchDeleteServers)
		}
	}

	//api := r.Group("api/")
	//
	//// bos user
	//r.GET("/logout", api.Logout)
	//r.POST("/login", api.Login)
	//r.POST("/change_password", api.ChangePassword)
	//user := r.Group("/user")
	//{
	//	user.Use(bos.Auth())
	//	user.POST("/change_password", api.ChangePassword)
	//	user.POST("/upload", api.UserUploadMedia)
	//	user.GET("/query/:dgst", api.UserQueryMedia)
	//	user.GET("/filelist", api.UserFileListMedia)
	//
	//}
	return r
}
