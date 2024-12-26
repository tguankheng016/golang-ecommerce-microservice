package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	miniRedis "github.com/alicebob/miniredis/v2"
	redisCli "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ctx = context.TODO()
)

func TestGetUserRolePermissions(t *testing.T) {
	allPermissions := make(map[string]struct{})
	for key := range permissions.GetAppPermissions().Items {
		allPermissions[key] = struct{}{}
	}

	tests := []struct {
		name            string
		userId          int64
		wantPermissions map[string]struct{}
	}{
		{
			name:            "get admin user granted permissions",
			userId:          1,
			wantPermissions: allPermissions,
		},
		{
			name:   "get normal user granted permissions",
			userId: 2,
			wantPermissions: map[string]struct{}{
				"Pages.Administration.Users": {},
				"Pages.Administration.Roles": {},
			},
		},
		{
			name:   "get normal user granted permissions with prohibited permissions",
			userId: 3,
			wantPermissions: map[string]struct{}{
				"Pages.Administration.Users": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDb, gormDb, mock := getMockDb(t)
			defer sqlDb.Close()

			mockByScenarioByUserId(mock, tt.userId)

			userRolePermissionManager := NewUserRolePermissionManager(gormDb, getMockRedisClient(t))
			gotPermissions, err := userRolePermissionManager.SetUserPermissions(ctx, tt.userId)
			assert.Equal(t, err, nil)
			assert.Equal(t, tt.wantPermissions, gotPermissions)
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

func mockByScenarioByUserId(mock sqlmock.Sqlmock, userId int64) {
	expectUserSQL := "^SELECT (.+) FROM \"users\" WHERE (.+)$"
	expectUserRoleSQL := "^SELECT (.+) FROM \"user_roles\" WHERE (.+)$"
	expectRoleSQL := "^SELECT (.+) FROM \"roles\" WHERE (.+)$"
	expectedUserPermissionSQL := "^SELECT (.+) FROM \"user_role_permissions\" WHERE user_id =(.+)$"
	expectedRolePermissionSQL := "^SELECT (.+) FROM \"user_role_permissions\" WHERE role_id =(.+)$"

	userRows := sqlmock.NewRows([]string{"id", "user_name"})
	userRoleRows := sqlmock.NewRows([]string{"user_id", "role_id"})
	roleRows := sqlmock.NewRows([]string{"id", "name"})
	roleRows2 := sqlmock.NewRows([]string{"id", "name"})
	userPermissionRows := sqlmock.NewRows([]string{"name", "user_id", "role_id", "is_granted"})
	rolePermissionRows := sqlmock.NewRows([]string{"name", "user_id", "role_id", "is_granted"})

	switch userId {
	case 1:
		mock.ExpectQuery(expectUserSQL).WillReturnRows(userRows.AddRow(userId, "admin"))
		mock.ExpectQuery(expectUserRoleSQL).WillReturnRows(userRoleRows.AddRow(userId, 1))
		mock.ExpectQuery(expectRoleSQL).WillReturnRows(roleRows.AddRow(1, "Admin"))
		mock.ExpectQuery(expectedUserPermissionSQL).WillReturnRows(userPermissionRows)
		mock.ExpectQuery(expectRoleSQL).WillReturnRows(roleRows2.AddRow(1, "Admin"))
		mock.ExpectQuery(expectedRolePermissionSQL).WillReturnRows(rolePermissionRows)
	case 2:
		mock.ExpectQuery(expectUserSQL).WillReturnRows(userRows.AddRow(userId, "user"))
		mock.ExpectQuery(expectUserRoleSQL).WillReturnRows(userRoleRows.AddRow(userId, 1))
		mock.ExpectQuery(expectRoleSQL).WillReturnRows(roleRows.AddRow(1, "User"))
		mock.ExpectQuery(expectedUserPermissionSQL).WillReturnRows(userPermissionRows)
		mock.ExpectQuery(expectRoleSQL).WillReturnRows(roleRows2.AddRow(1, "User"))

		rolePermissionRows.AddRow("Pages.Administration.Users", nil, 1, true)
		rolePermissionRows.AddRow("Pages.Administration.Roles", nil, 1, true)
		mock.ExpectQuery(expectedRolePermissionSQL).WillReturnRows(rolePermissionRows)
	case 3:
		mock.ExpectQuery(expectUserSQL).WillReturnRows(userRows.AddRow(userId, "user"))
		mock.ExpectQuery(expectUserRoleSQL).WillReturnRows(userRoleRows.AddRow(userId, 1))
		mock.ExpectQuery(expectRoleSQL).WillReturnRows(roleRows.AddRow(1, "User"))

		userPermissionRows.AddRow("Pages.Administration.Roles", userId, nil, false)

		mock.ExpectQuery(expectedUserPermissionSQL).WillReturnRows(userPermissionRows)
		mock.ExpectQuery(expectRoleSQL).WillReturnRows(roleRows2.AddRow(1, "User"))

		rolePermissionRows.AddRow("Pages.Administration.Users", nil, 1, true)
		rolePermissionRows.AddRow("Pages.Administration.Roles", nil, 1, true)
		mock.ExpectQuery(expectedRolePermissionSQL).WillReturnRows(rolePermissionRows)
	}
}

func getMockRedisClient(t *testing.T) *redisCli.Client {
	s := miniRedis.RunT(t)
	client := redisCli.NewClient(&redisCli.Options{
		Addr: s.Addr(),
	})

	return client
}
