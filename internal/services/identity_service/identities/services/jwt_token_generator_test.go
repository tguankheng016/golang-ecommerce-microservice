package services

import (
	"context"
	"database/sql"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	miniRedis "github.com/alicebob/miniredis/v2"
	"github.com/gofrs/uuid"
	jwtGo "github.com/golang-jwt/jwt/v5"
	redisCli "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/constants"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ctx = context.TODO()
)

func TestGenerateToken(t *testing.T) {
	securityStamp, _ := uuid.NewV4()
	refreshTokenKey, _ := uuid.NewV4()
	client := getMockRedisClient(t)
	defaultAuthOptions := &jwt.AuthOptions{
		Issuer:    "testing",
		Audience:  "testing",
		SecretKey: "G4jR8dL9pK6cE3nB2H7mQ8iJ9tY5aS4rF7gT6eW8zC2xV",
	}

	tests := []struct {
		name            string
		scenario        int
		user            *models.User
		refreshTokenKey string
		tokenType       jwt.TokenType
		wantErr         bool
		wantErrMessage  string
	}{
		{
			name:            "can generate access token",
			scenario:        0,
			user:            &models.User{Id: 1, SecurityStamp: securityStamp},
			refreshTokenKey: refreshTokenKey.String(),
			tokenType:       jwt.AccessToken,
			wantErr:         false,
			wantErrMessage:  "",
		},
		{
			name:            "can generate refresh token",
			scenario:        0,
			user:            &models.User{Id: 1, SecurityStamp: securityStamp},
			refreshTokenKey: "",
			tokenType:       jwt.RefreshToken,
			wantErr:         false,
			wantErrMessage:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDb, gormDb, mock := getMockDb(t)
			defer sqlDb.Close()

			mockByScenario(mock, tt.scenario)

			jwtTokenGenerator := NewJwtTokenGenerator(gormDb, client, defaultAuthOptions)
			var accessToken string
			var err error
			if tt.tokenType == jwt.AccessToken {
				accessToken, _, err = jwtTokenGenerator.GenerateAccessToken(ctx, tt.user, tt.refreshTokenKey)
			} else {
				accessToken, _, _, err = jwtTokenGenerator.GenerateRefreshToken(ctx, tt.user)
			}

			assert.Equal(t, err != nil, tt.wantErr)
			if tt.wantErr {
				assert.Equal(t, tt.wantErrMessage, err.Error())
			} else {
				assert.NotEmpty(t, accessToken)

				token, err := jwtGo.ParseWithClaims(accessToken, jwtGo.MapClaims{}, func(token *jwtGo.Token) (interface{}, error) {
					return []byte(defaultAuthOptions.SecretKey), nil
				})
				assert.False(t, err != nil)
				assert.True(t, token.Valid)
				assert.Equal(t, defaultAuthOptions.Issuer, token.Header["iss"])
				assert.Equal(t, defaultAuthOptions.Audience, token.Header["aud"])

				claims, ok := token.Claims.(jwtGo.MapClaims)
				assert.True(t, ok)

				tokenType, _ := strconv.Atoi(claims["token_type"].(string))
				assert.Equal(t, int(tt.tokenType), tokenType)

				sub, _ := claims["sub"].(string)
				assert.Equal(t, strconv.FormatInt(tt.user.Id, 10), sub)

				securityStamp, _ := claims[constants.SecurityStampKey].(string)
				assert.Equal(t, tt.user.SecurityStamp.String(), securityStamp)

				tokenKey, ok := claims[constants.TokenValidityKey].(string)
				assert.True(t, ok)
				assert.NotEmpty(t, tokenKey)

				if tt.refreshTokenKey != "" {
					refreshTokenKeyRetrieved, ok := claims[constants.RefreshTokenValidityKey].(string)
					assert.True(t, ok)
					assert.NotEmpty(t, refreshTokenKeyRetrieved)
					assert.Equal(t, tt.refreshTokenKey, refreshTokenKeyRetrieved)
				}
			}
		})
	}
}

func getMockDb(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	gormdb, err := gorm.Open(gorm_postgres.New(gorm_postgres.Config{
		Conn: sqldb,
	}))
	if err != nil {
		t.Fatal(err)
	}

	return sqldb, gormdb, mock
}

func mockByScenario(mock sqlmock.Sqlmock, scenario int) {
	expectedSQL := "INSERT INTO \"user_tokens\" (.+) VALUES (.+)"
	addRow := sqlmock.NewRows([]string{"id"}).AddRow("1")

	mock.ExpectBegin()

	switch scenario {
	default:
		mock.ExpectQuery(expectedSQL).WillReturnRows(addRow)
	}

	mock.ExpectCommit()
}

func getMockRedisClient(t *testing.T) *redisCli.Client {
	s := miniRedis.RunT(t)
	client := redisCli.NewClient(&redisCli.Options{
		Addr: s.Addr(),
	})

	return client
}
