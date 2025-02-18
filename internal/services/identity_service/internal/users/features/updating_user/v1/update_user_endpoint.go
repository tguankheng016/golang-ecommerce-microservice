package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	v "github.com/RussellLuo/validating/v3"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	roleService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/roles/services"
	userConsts "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/constants"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/services"
)

// Request
type HumaUpdateUserRequest struct {
	Body struct {
		dtos.EditUserDto
	}
}

// Result
type HumaUpdateUserResult struct {
	Body struct {
		User dtos.UserDto
	}
}

// Validator
func (e HumaUpdateUserRequest) Schema() v.Schema {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return v.Schema{
		v.F("id", e.Body.Id): v.All(
			v.Nonzero[*int64]().Msg("Invalid user id"),
			v.Nested(func(ptr *int64) v.Validator { return v.Value(*ptr, v.Gt(int64(0)).Msg("Invalid user id")) }),
		),
		v.F("username", e.Body.UserName):    v.Nonzero[string]().Msg("Please enter the username"),
		v.F("first_name", e.Body.FirstName): v.Nonzero[string]().Msg("Please enter the first name"),
		v.F("last_name", e.Body.LastName):   v.Nonzero[string]().Msg("Please enter the last name"),
		v.F("email", e.Body.Email): v.All(
			v.Nonzero[string]().Msg("Please enter the email"),
			v.Match(pattern).Msg("Please enter a valid email"),
		),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	userRolePermissionManager services.IUserRolePermissionManager,
	publisher message.Publisher,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "UpdateUser",
			Method:        http.MethodPut,
			Path:          "/identities/user",
			Summary:       "Update User",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesAdministrationUsersEdit),
				postgres.SetupTransaction(api, pool),
			},
		},
		updateUser(userRolePermissionManager, publisher),
	)
}

func updateUser(userRolePermissionManager services.IUserRolePermissionManager, publisher message.Publisher) func(context.Context, *HumaUpdateUserRequest) (*HumaUpdateUserResult, error) {
	return func(ctx context.Context, request *HumaUpdateUserRequest) (*HumaUpdateUserResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		userManager := services.NewUserManager(tx)

		user, err := userManager.GetUserById(ctx, *request.Body.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if user == nil {
			return nil, huma.Error404NotFound("user not found")
		}

		if user.UserName == userConsts.DefaultAdminUserName && request.Body.UserName != user.UserName {
			return nil, huma.Error400BadRequest("you cannot change admin username")
		}

		if user.UserName == userConsts.DefaultAdminUserName && request.Body.Password != "" {
			return nil, huma.Error400BadRequest("you cannot change admin password")
		}

		if err := copier.Copy(&user, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := userManager.UpdateUser(ctx, user, request.Body.Password); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		if len(request.Body.RoleIds) > 0 {
			roleManager := roleService.NewRoleManager(tx)

			for _, roleId := range request.Body.RoleIds {
				if roleId == 0 {
					continue
				}

				role, err := roleManager.GetRoleById(ctx, roleId)
				if err != nil {
					return nil, huma.Error500InternalServerError(err.Error())
				}
				if role == nil {
					return nil, huma.Error400BadRequest("Invalid role id")
				}
			}
		}

		userRoleUpdated, err := userManager.UpdateUserRoles(ctx, user, request.Body.RoleIds)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if userRoleUpdated {
			userRolePermissionManager.RemoveUserRoleCaches(ctx, user.Id)
		}

		var userDto dtos.UserDto
		if err := copier.Copy(&userDto, &user); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var userUpdatedEvent events.UserUpdatedEvent
		if err := copier.Copy(&userUpdatedEvent, &user); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		payload, err := json.Marshal(userUpdatedEvent)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)
		publisher.Publish(events.UserUpdatedTopicV1, msg)

		result := HumaUpdateUserResult{}
		result.Body.User = userDto

		return &result, nil
	}
}
