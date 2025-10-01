package main

import (
	"fmt"
	"net"
	"os"
	"time"

	_ "github.com/faizisyellow/indocoffee/docs"
	"github.com/faizisyellow/indocoffee/internal/auth"
	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/logger"
	"github.com/faizisyellow/indocoffee/internal/repository/beans"
	"github.com/faizisyellow/indocoffee/internal/repository/carts"
	"github.com/faizisyellow/indocoffee/internal/repository/forms"
	"github.com/faizisyellow/indocoffee/internal/repository/invitations"
	"github.com/faizisyellow/indocoffee/internal/repository/products"
	"github.com/faizisyellow/indocoffee/internal/repository/roles"
	"github.com/faizisyellow/indocoffee/internal/repository/users"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/uploader/uploadthing"
	"github.com/faizisyellow/indocoffee/internal/utils"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
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
// @host		localhost:8080
// @BasePath	/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatalw("error loading .env file", zap.Error(err))
	}

	dbConfig := DBConf{
		Addr:            os.Getenv("DB_ADDR"),
		MaxOpenConn:     30,
		MaxIdleConn:     30,
		MaxLifeTime:     "4m",
		MaxIdleLifeTime: "4m",
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

	services := service.New(
		&users.UsersRepository{Db: dbs},
		&invitations.InvitationRepository{Db: dbs},
		&beans.BeansRepository{Db: dbs},
		&forms.FormsRepository{Db: dbs},
		&roles.RolesRepository{Db: dbs},
		&products.ProductRepository{Db: dbs},
		upt,
		&db.TransactionDB{Db: dbs},
		ud,
		&carts.CartsRepository{Db: dbs},
	)

	jwtTokenConfig := JwtConfig{
		SecretKey: os.Getenv("SECRET_KEY"),
		Iss:       "authentication",
		Sub:       "user",
		Exp:       time.Now().Add(time.Hour * 24 * 3).Unix(),
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
