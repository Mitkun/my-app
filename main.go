package main

import (
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/gormc"
	"log"
	"my-app/builder"
	"my-app/common"
	"my-app/component"
	"my-app/middleware"
	"my-app/module/image"
	"my-app/module/product/controller"
	productusecase "my-app/module/product/domain/usecase"
	"my-app/module/product/infras/producthttp"
	productmysql "my-app/module/product/repository/mysql"
	"my-app/module/user/infras/httpservice"
	"my-app/module/user/infras/repository"
	"my-app/module/user/usecase"
	"net/http"
)

func newService() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName("G11"),
		sctx.WithComponent(gormc.NewGormDB(common.KeyGorm, "")),
		sctx.WithComponent(component.NewJWT(common.KeyJWT)),
		sctx.WithComponent(component.NewAWSS3Provider(common.KeyAWSS3)),
	)
}

func main() {

	service := newService()

	service.OutEnv()

	if err := service.Load(); err != nil {
		log.Fatalln(err)
	}

	db := service.MustGet(common.KeyGorm).(common.DbContext).GetDB()

	r := gin.Default()

	r.Use(middleware.Recovery())

	tokenProvider := service.MustGet(common.KeyJWT).(component.TokenProvider)

	authClient := usecase.NewIntrospectUC(repository.NewUserRepo(db), repository.NewSessionMySQLRepo(db), tokenProvider)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.DELETE("/v1/revoke-token", middleware.RequireAuth(authClient), func(c *gin.Context) {
		requester := c.MustGet(common.KeyRequester).(common.Requester)

		repo := repository.NewSessionMySQLRepo(db)
		if err := repo.Delete(c.Request.Context(), requester.TokenId()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": true})
	})

	// Setup dependencies
	repo := productmysql.NewMysqlRepository(db)
	useCase := productusecase.NewCreateProductUseCase(repo)
	api := controller.NewAPIController(useCase)

	v1 := r.Group("/v1")
	{
		products := v1.Group("/products")
		{
			products.POST("", api.CreateProductAPI(db))
		}

	}

	userUseCase := usecase.UseCaseWithBuilder(builder.NewSimpleBuilder(db, tokenProvider))

	httpservice.NewUserService(userUseCase, service).SetAuthClient(authClient).Routes(v1)
	image.NewHTTPService(service).Routes(v1)
	producthttp.NewHttpService(service).Routes(v1)

	err := r.Run(":3000")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
