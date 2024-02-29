package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	jwtCecret := os.Getenv("JWT_SECRET")

	tokenProvider := component.NewJWTProvider(jwtCecret,
		60*60*24*7, 60*60*24*24)

	authClient := usecase.NewIntrospectUC(repository.NewUserRepo(db), repository.NewSessionMySQLRepo(db), tokenProvider)

	r.GET("/ping", middleware.RequireAuth(authClient), func(c *gin.Context) {

		requester := c.MustGet(common.KeyRequester).(common.Requester)

		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"requester": requester.LastName(),
		})
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

	//userUC := usecase.NewUseCase(repository.NewUserRepo(db), repository.NewSessionMySQLRepo(db), &common.Hasher{}, tokenProvider)

	userUseCase := usecase.UseCaseWithBuilder(builder.NewSimpleBuilder(db, tokenProvider))

	httpservice.NewUserService(userUseCase).Routes(v1)

	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type fakeAuthClient struct{}

func (fakeAuthClient) IntrospectToken(ctx context.Context, assessToken string) (common.Requester, error) {
	return common.NewRequester(
		uuid.MustParse("018dd9c5-3da2-710a-aa53-9ba16fb8d451"),
		uuid.MustParse("018ddf53-d3aa-793f-bde7-dc9611970497"),
		"Thiện",
		"Trương Cong",
		"user",
		"acttivated",
	), nil
}

//type mockSessionRepo struct {
//}
//
//func (m mockSessionRepo) Create(ctx context.Context, data *domain.Session) error {
//	return nil
//}
