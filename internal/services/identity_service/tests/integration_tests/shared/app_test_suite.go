package shared

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/caching"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/messaging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	pg "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security"
	jwt "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/configurations"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/data/seeds"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities"
	identityService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users"
	userConsts "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/constants"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/services"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type AppTestSuite struct {
	suite.Suite
	App               *fxtest.App
	Config            *config.Config
	Pool              *pgxpool.Pool
	Ctx               context.Context
	PgContainer       *postgres.PostgresContainer
	Client            *resty.Client
	CacheManager      *cache.Cache[string]
	JwtTokenGenerator identityService.IJwtTokenGenerator
	JwtTokenHandler   jwt.IJwtTokenHandler
}

// this function executes before the test suite begins execution
func (suite *AppTestSuite) SetupSuite() {
	fmt.Println(">>> From SetupSuite")

	if err := os.Setenv("APP_ENV", environment.Test.GetEnvironmentName()); err != nil {
		suite.T().Fatalf("failed to set test environment: %v", err)
	}

	suite.Ctx = context.Background()

	pgContainer, dbPool, err := CreatePostgresContainer(suite.Ctx)
	if err != nil {
		suite.T().Fatalf("failed to set postgres test container: %v", err)
	}

	if err := pg.RunGooseMigration(dbPool); err != nil {
		suite.T().Fatalf("failed to set migration for container: %v", err)
	}

	suite.Pool = dbPool
	suite.PgContainer = pgContainer

	app := fxtest.New(
		suite.T(),
		fx.Options(
			fx.Provide(
				func() *pgxpool.Pool {
					return dbPool
				},
				environment.ConfigureAppEnv,
				config.InitConfig,
			),
			logging.Module,
			caching.Module,
			security.Module,
			permissions.Module,
			identities.Module,
			users.Module,
			seeds.Module,
			configurations.Module,
			http.Module,
			grpc.Module,
			messaging.Module,
			fx.Invoke(func(
				config *config.Config,
				cacheManager *cache.Cache[string],
				jwtTokenGenerator identityService.IJwtTokenGenerator,
				jwtTokenHandler jwt.IJwtTokenHandler,
			) {
				suite.Config = config
				suite.Client = resty.New()
				suite.Client.SetBaseURL(suite.Config.ServerOptions.GetBasePath())
				suite.CacheManager = cacheManager
				suite.JwtTokenGenerator = jwtTokenGenerator
				suite.JwtTokenHandler = jwtTokenHandler
			}),
		),
	)

	suite.App = app

	if err := suite.App.Start(suite.Ctx); err != nil {
		suite.T().Fatalf("failed to start the Uber FX application: %v", err)
	}
}

// this function executes after all tests executed
func (suite *AppTestSuite) TearDownSuite() {
	fmt.Println(">>> From TearDownSuite")

	if err := suite.App.Stop(suite.Ctx); err != nil {
		log.Fatalf("error stopping app: %s", err)
	}

	if err := suite.PgContainer.Terminate(suite.Ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *AppTestSuite) LoginAs(userName string) (string, error) {
	userManager := userService.NewUserManager(suite.Pool)
	user, err := userManager.GetUserByUserName(suite.Ctx, userName)
	if err != nil {
		return "", nil
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	accessToken, _, err := suite.JwtTokenGenerator.GenerateAccessToken(suite.Ctx, *user, "")
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (suite *AppTestSuite) LoginAsAdmin() (string, error) {
	return suite.LoginAs(userConsts.DefaultAdminUserName)
}

func (suite *AppTestSuite) LoginAsUser() (string, error) {
	return suite.LoginAs(userConsts.DefaultUserUserName)
}
