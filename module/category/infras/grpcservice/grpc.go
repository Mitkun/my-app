package grpcservice

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"google.golang.org/grpc"
	"log"
	"my-app/common"
	"my-app/module/category/query"
	"my-app/proto/category"
	"net"
)

type service struct {
	port int
	sctx sctx.ServiceContext
}

func NewCateGRPCService(port int, sctx sctx.ServiceContext) *service {
	return &service{port: port, sctx: sctx}
}

type categoryServer struct {
	category.UnimplementedCategoryServer
	sctx sctx.ServiceContext
}

func NewCategoryServer(sctx sctx.ServiceContext) *categoryServer {
	return &categoryServer{sctx: sctx}
}

func (cs *categoryServer) GetCategoriesByIds(ctx context.Context, request *category.GetCateIdsRequest) (*category.CateIdsResponse, error) {
	var cates []query.CategoryDTO

	dbCtx := cs.sctx.MustGet(common.KeyGorm).(common.DbContext)

	ids := make([]uuid.UUID, len(request.Ids))
	for i := range ids {
		ids[i] = uuid.MustParse(request.Ids[i])
	}

	if err := dbCtx.GetDB().Table(query.CategoryDTO{}.TableName()).
		Where("id in (?)", ids).
		Find(&cates).Error; err != nil {
		return nil, core.ErrBadRequest.WithError("cannot list categories").WithDebug(err.Error())
	}

	results := make([]*category.CategoryDTO, len(cates))

	for i := range results {
		results[i] = &category.CategoryDTO{
			Id:    cates[i].Id.String(),
			Title: cates[i].Title,
		}
	}

	return &category.CateIdsResponse{Data: results}, nil
}

func (s *service) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	// Create a gRPC server object
	grpcServ := grpc.NewServer()
	// Attach the Greeter service to the server

	category.RegisterCategoryServer(grpcServ, NewCategoryServer(s.sctx))

	log.Println(fmt.Sprintf("Serving gRPC on 0.0.0.0:%d", s.port))
	log.Fatal(grpcServ.Serve(lis))

	return nil
}
