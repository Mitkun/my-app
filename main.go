package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"my-app/builder"
	"my-app/common"
	"my-app/component"
	"my-app/middleware"
	"my-app/module/product/controller"
	productusecase "my-app/module/product/domain/usecase"
	productmysql "my-app/module/product/repository/mysql"
	"my-app/module/user/infras/httpservice"
	"my-app/module/user/infras/repository"
	"my-app/module/user/usecase"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	//r.Use(middleware.RequireAuth())

	jwtSecret := os.Getenv("JWT_SECRET")

	tokenProvider := component.NewJWTProvider(jwtSecret,
		60*60*24*7, 60*60*24*14)

	authClient := usecase.NewIntrospectUC(repository.NewUserRepo(db), repository.NewSessionMySQLRepo(db), tokenProvider)

	r.GET("/ping", middleware.RequireAuth(authClient), func(c *gin.Context) {

		requester := c.MustGet(common.KeyRequester).(common.Requester)

		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"requester": requester.LastName(),
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

	httpservice.NewUserService(userUseCase).Routes(v1)

	err = r.Run(":3000")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
