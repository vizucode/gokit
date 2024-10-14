package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/vizucode/gokit/adapter/dbc"
	"github.com/vizucode/gokit/config"
	"github.com/vizucode/gokit/factory/server"
	gokitrpc "github.com/vizucode/gokit/factory/server/grpc"
	"github.com/vizucode/gokit/factory/server/rest"
	"github.com/vizucode/gokit/logger"
	pb "github.com/vizucode/gokit/protoc"
	"github.com/vizucode/gokit/utils/constant"
	"github.com/vizucode/gokit/utils/errorkit"
	"github.com/vizucode/gokit/utils/request"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type restRoute struct {
	pb.UnimplementedHelloServiceServer

	GORM  *gorm.DB
	SQLDB *sql.DB
	REDIS *redis.Client
	API   request.Client
}

func NewHandler(gormDB *gorm.DB, sqlDB *sql.DB, redis *redis.Client, apiClient request.Client) *restRoute {
	return &restRoute{
		GORM:  gormDB,
		SQLDB: sqlDB,
		REDIS: redis,
		API:   apiClient,
	}
}

func (r *restRoute) Router(router fiber.Router) {
	v1 := router.Group("v1")
	v1.Get("/with-gorm", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		var jabatan []Jabatan

		err := r.GORM.Model(&Jabatan{}).Find(&jabatan).Error
		if err != nil {
			logger.Log.Error(ctx, err)
			return err
		}

		return c.Status(200).JSON(map[string]interface{}{
			"data": jabatan,
		})
	})

	v1.Get("/with-sql", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		var jabatans []Jabatan

		rows, err := r.SQLDB.Query("SELECT * FROM jabatans")
		if err != nil {
			logger.Log.Error(ctx, err)
			return err
		}

		for rows.Next() {
			var jabatan Jabatan
			err = rows.Scan(&jabatan.Id, &jabatan.CreatedAt, &jabatan.UpdatedAt, &jabatan.CreatedBy, &jabatan.UpdatedBy, &jabatan.Name)
			if err != nil {
				logger.Log.Error(ctx, err)
				return err
			}
			jabatans = append(jabatans, jabatan)
		}

		return c.Status(200).JSON(map[string]interface{}{
			"data": jabatans,
		})
	})

	v1.Get("/with-redis", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		redisCMD := r.REDIS.Get(ctx, "example")
		if redisCMD.Err() != nil {
			logger.Log.Error(ctx, redisCMD.Err())
			return redisCMD.Err()
		}

		result, err := redisCMD.Result()
		if err != nil {
			logger.Log.Error(ctx, err)
			return err
		}

		return c.Status(200).JSON(map[string]interface{}{
			"data": result,
		})
	})

	v1.Get("/with-request-rest", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		req := r.API.Request(nil, "https://jsonplaceholder.typicode.com/todos/1", "Get:JsonPlaceholder")

		res, sc, err := req.Get(ctx)
		if err != nil {
			logger.Log.Error(ctx, err)
			return err
		}

		var resp map[string]interface{}
		err = json.Unmarshal(res, &resp)
		if err != nil {
			logger.Log.Errorf(ctx, "actual response %s. error %s", string(res), err.Error())

			return err
		}

		return c.Status(sc).JSON(map[string]interface{}{
			"data": resp,
		})

	})

	v1.Get("/with-request-grpc", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		host := fmt.Sprintf("%s:%s", "localhost", "3005")
		intercept := logger.NewInterceptor(host)

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(intercept.ChainUnaryClient(
				intercept.UnaryClientTracerInterceptor,
			)),
		}

		conn, err := grpc.NewClient(host, opts...)
		if err != nil {
			logger.Log.Error(ctx, err)
			return err
		}

		rpcHello := pb.NewHelloServiceClient(conn)
		respHello, err := rpcHello.SayHello(ctx, &pb.HelloRequest{
			Name: "Example Name",
		})
		if err != nil {
			logger.Log.Error(ctx, err)
			return errorkit.NewErrorStd(http.StatusInternalServerError, errorkit.InternalServer, err.Error())
		}

		return c.Status(200).JSON(map[string]interface{}{
			"data": respHello.GetMessage(),
		})
	})
}

func (r *restRoute) Register(srv *grpc.Server) {
	pb.RegisterHelloServiceServer(srv, r)
}

func (r *restRoute) SayHello(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloResponse, err error) {

	resp = &pb.HelloResponse{
		Message: "Hello " + req.Name,
	}

	return resp, nil
}

type Jabatan struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt time.Time
	CreatedBy sql.NullString
	UpdatedBy sql.NullString
	Name      string
}

func main() {

	serviceName := "TestingGoKit"

	config.Load(serviceName, ".")

	apiClient := request.NewRequest(&http.Client{
		Timeout: 5 * time.Second,
	})

	gormDB := dbc.NewGormConnection(
		dbc.SetGormURIConnection("host=localhost user=example password=examplepassword dbname=db_kpri_sehat port=5432 sslmode=disable TimeZone=Asia/Jakarta"),
		dbc.SetGormDatabaseName("db_example"),
		dbc.SetGormDriver(constant.Postgres),
		dbc.SetGormMaxIdleConnection(5),
		dbc.SetGormMaxPoolConnection(50),
	)

	sqlDB := dbc.NewSqlConnection(
		dbc.SetSqlDatabaseName("db_example"),
		dbc.SetSqlURIConnection("host=localhost user=example password=examplepassword dbname=db_kpri_sehat port=5432 sslmode=disable TimeZone=Asia/Jakarta"),
		dbc.SetSqlDriver(constant.Postgres),
	)

	redisRead := dbc.NewRedisConnection(
		dbc.SetRedisMaxIdleConnectionDuration(5*time.Minute),
		dbc.SetRedisURIConnection("redis://:@localhost:6379/0"),
		dbc.SetRedisServiceName(serviceName),
		dbc.SetRedisMinPoolConnection(1),
		dbc.SetRedisMaxPoolConnection(5),
	)

	app := server.NewService(
		server.SetServiceName(serviceName),
		server.SetRestHandler(NewHandler(gormDB.DB, sqlDB.DB, redisRead.DB, apiClient)),
		server.SetRestHandlerOptions(
			rest.SetHTTPHost("localhost"),
			rest.SetHTTPPort(3000),
			rest.SetErrorHandler(fiber.DefaultErrorHandler),
		),
		server.SetGrpcHandler(NewHandler(gormDB.DB, sqlDB.DB, redisRead.DB, apiClient)),
		server.SetGrpcHandlerOptions(
			gokitrpc.SetTCPHost("localhost"),
			gokitrpc.SetTCPPort(3001),
		),
	)

	appServer := server.New(app)
	appServer.Run()
}
