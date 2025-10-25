package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/faizisyellow/indocoffee/docs"
	"github.com/faizisyellow/indocoffee/internal/auth"
	"github.com/faizisyellow/indocoffee/internal/db"
	loginLimiter "github.com/faizisyellow/indocoffee/internal/limiter/login"
	"github.com/faizisyellow/indocoffee/internal/logger"
	"github.com/faizisyellow/indocoffee/internal/repository/beans"
	"github.com/faizisyellow/indocoffee/internal/repository/carts"
	"github.com/faizisyellow/indocoffee/internal/repository/forms"
	"github.com/faizisyellow/indocoffee/internal/repository/invitations"
	"github.com/faizisyellow/indocoffee/internal/repository/orders"
	"github.com/faizisyellow/indocoffee/internal/repository/products"
	"github.com/faizisyellow/indocoffee/internal/repository/roles"
	"github.com/faizisyellow/indocoffee/internal/repository/users"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/uploader/uploadthing"
	"github.com/faizisyellow/indocoffee/internal/utils"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

//	@title			Indocoffee REST APIs service
//	@version		1.0
//	@description	Rest API Documentation.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						Authorization

// @schemes	http https
// @BasePath	/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatalw("error loading .env file", zap.Error(err))
	}

	docs.SwaggerInfo.Host = net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT"))

	dbConfig := DBConf{
		Addr:            os.Getenv("DB_ADDR"),
		MaxOpenConn:     5,
		MaxIdleConn:     5,
		MaxLifeTime:     "3m",
		MaxIdleLifeTime: "3m",
	}

	dbs, err := db.New(
		dbConfig.Addr,
		dbConfig.MaxOpenConn,
		dbConfig.MaxIdleConn,
		dbConfig.MaxIdleLifeTime,
		dbConfig.MaxLifeTime,
	)
	if err != nil {
		logger.Logger.Fatalw("error connecting to database", zap.Error(err))
	}

	defer dbs.Close()
	logger.Logger.Infow("database connection pool has established")

	ud := utils.UUID{Plaintoken: uuid.New().String()}

	upt := uploadthing.New(
		os.Getenv("UPLOADTHING_API_KEY"),
		os.Getenv("UPLOADTHING_PRESIGNED_URL"),
		os.Getenv("UPLOADTHING_POOL_UPLOAD_URL"),
		"public-read",
		"imageUploader",
		os.Getenv("UPLOADTHING_UPLOAD_BY"),
		os.Getenv("UPLOADTHING_META_URL"),
		os.Getenv("UPLOADTHING_CALLBACK_URL"),
		os.Getenv("UPLOADTHING_DELETE_URL"),
		os.Getenv("UPLOADTHING_APP_ID"),
	)

	rdbname, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logger.Logger.Fatalw("error loading .env file", zap.Error(err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PW"),
		DB:       rdbname,
	})

	defer rdb.Close()

	loginRateLimiter := loginLimiter.RedisLoginLimiter{
		Rdb:      rdb,
		Limit:    2, // start from 0
		Duration: 12 * time.Hour,
	}

	services := service.New(
		&loginRateLimiter,
		&users.UsersRepository{Db: dbs},
		&invitations.InvitationRepository{Db: dbs},
		&beans.BeansRepository{Db: dbs},
		&forms.FormsRepository{Db: dbs},
		&roles.RolesRepository{Db: dbs},
		&products.ProductRepository{Db: dbs},
		upt,
		&db.TransactionDB{Db: dbs},
		ud,
		utils.Ulid(ulid.Make().String),
		&carts.CartsRepository{Db: dbs},
		&orders.OrdersRepository{Db: dbs},
	)

	jwtTokenConfig := JwtConfig{
		SecretKey: os.Getenv("SECRET_KEY"),
		Iss:       "authentication",
		Sub:       "user",
		Exp:       24 * time.Hour,
	}

	jwtAuthentication := auth.New(jwtTokenConfig.SecretKey, jwtTokenConfig.Iss, jwtTokenConfig.Sub)

	application := Application{
		Port:           os.Getenv("PORT"),
		Host:           os.Getenv("HOST"),
		Env:            os.Getenv("ENV"),
		DbConfig:       dbConfig,
		Services:       *services,
		JwtAuth:        jwtTokenConfig,
		Authentication: jwtAuthentication,
		Logger:         logger.Logger,

		//http:domain:port/version/swagger/*
		SwaggerUrl: fmt.Sprintf("http://%v/v%v/swagger/doc.json", net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT")), 1),
	}

	err = application.Run(application.Mux())
	if err != nil {
		logger.Logger.Fatalw("error running application", zap.Error(err))
	}

}
