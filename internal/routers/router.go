package routers

import (
	"blog-service/global"
	"blog-service/internal/middleware"
	"blog-service/internal/routers/api"
	v1 "blog-service/internal/routers/api/v1"
	"blog-service/pkg/limiter"
	"github.com/gin-gonic/gin"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	//if global.ServerSetting.RunMode == "debug" {
	//	r.Use(gin.Logger())
	//	r.Use(gin.Recovery())
	//} else {
	r.Use(middleware.CorsMiddleware()) //跨域设置
	r.Use(middleware.AccessLog())      //请求日志
	r.Use(middleware.Recovery())
	//}
	r.Use(middleware.RateLimiter(methodLimiters))                             //限流
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout)) //超时时间

	r.POST("/auth", api.GetAuth)

	//upload := v1.NewUpload()
	//r.POST("/upload/file", upload.UploadFile)
	//r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))

	tag := v1.NewTag()
	article := v1.NewArticle()
	chat := v1.NewChat()

	apiv1 := r.Group("/api/v1").Use(middleware.JWT())
	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags/:id", tag.Get)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}
	r.POST("/session", chat.Session)
	r.POST("/chat-process", chat.ChatProcess)

	return r
}
