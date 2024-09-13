package common

import (
	"fmt"
	"strconv"

	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/enums"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/helpers"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/repositories"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReqContextKey struct{}
type ReqFileKey struct{}

type RequestContext struct {
	token           string
	id              string
	objId           primitive.ObjectID
	ip              string
	logger          *logger.Logger
	method          string
	path            string
	role            string
	is_auth         bool
	payload         *map[string]interface{}
	timezone_offset int
}

func NewRequestContext() *RequestContext {
	return &RequestContext{}
}

func (ctx *RequestContext) ParseHeader(headers map[string][]string) *RequestContext {
	if len(headers["X-Token"]) > 0 {
		ctx.token = headers["X-Token"][0]
	}
	if len(headers["Tz-Offset"]) > 0 {
		timezone_offset, err := strconv.Atoi(headers["Tz-Offset"][0])
		if err == nil {
			ctx.timezone_offset = timezone_offset
		}
	}
	return ctx
}

func (ctx *RequestContext) Auth() error {
	payload, err := helpers.VerifyAuthToken(ctx.token)
	if err != nil {
		return err
	}
	ctx.is_auth = true
	ctx.payload = &payload
	objId, err := primitive.ObjectIDFromHex(payload["_id"].(string))
	if err != nil {
		return exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
	}
	ctx.objId = objId
	return nil
}

func (ctx *RequestContext) GetIPForwardedFor(headers map[string][]string) *RequestContext {
	if len(headers["X-Forwarded-For"]) > 0 {
		ctx.ip = headers["X-Forwarded-For"][0]
	}
	return ctx
}

func (ctx *RequestContext) SetIP(ip string) *RequestContext {
	if ctx.ip == "" {
		ctx.ip = ip
	}
	return ctx
}

func (ctx *RequestContext) AuthUser() error {
	payload := ctx.payload
	if !ctx.is_auth {
		err := ctx.Auth()
		if err != nil {
			return err
		}
	}
	role := (*payload)["role"].(string)
	if role != enums.USER_ROLE_USER {
		return exceptions.Err_FORBIDEN().SetMessage("FORBIDEN")
	}

	user := entities.NewUserEntity()
	err, is_not_found := repositories.GetUserRepository().FindOneById(ctx.objId, user)
	if err != nil {
		return err
	}
	if is_not_found {
		return exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
	}
	ctx.id = user.Id.Hex()
	ctx.objId = user.Id
	ctx.role = enums.USER_ROLE_USER
	return nil
}

func (ctx *RequestContext) AuthAdmin() error {
	payload := ctx.payload
	if !ctx.is_auth {
		err := ctx.Auth()
		if err != nil {
			return err
		}
	}
	role := (*payload)["role"].(string)
	if role != enums.USER_ROLE_ADMIN {
		return exceptions.Err_FORBIDEN().SetMessage("FORBIDEN")
	}
	admin := entities.NewAdminEntity()
	err, is_not_found := repositories.GetAdminRepository().FindOneById(ctx.objId, admin)
	if err != nil {
		return err
	}
	if is_not_found {
		return exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
	}
	ctx.id = admin.Id.Hex()
	ctx.objId = admin.Id
	ctx.role = enums.USER_ROLE_ADMIN
	return nil
}

func (ctx *RequestContext) GetId() (primitive.ObjectID, error) {
	if !ctx.is_auth {
		err := ctx.Auth()
		if err != nil {
			return primitive.NilObjectID, err
		}
	}
	return primitive.ObjectIDFromHex(ctx.id)
}

func (ctx *RequestContext) GetObjId() primitive.ObjectID {
	if !ctx.is_auth {
		err := ctx.Auth()
		if err != nil {
			return primitive.NilObjectID
		}
	}
	return ctx.objId
}
func (ctx *RequestContext) SetObjId(id primitive.ObjectID) *RequestContext {
	ctx.objId = id
	ctx.is_auth = true
	return ctx
}

func (ctx *RequestContext) GetIp() string {
	return ctx.ip
}

func (ctx *RequestContext) GetRole() string {
	return ctx.role
}

func (ctx *RequestContext) GetTZOffset() int {
	return ctx.timezone_offset
}

func (ctx *RequestContext) SetMethod(method string) *RequestContext {
	ctx.method = method
	return ctx
}
func (ctx *RequestContext) SetPath(path string) *RequestContext {
	ctx.path = path
	return ctx
}

func (ctx *RequestContext) GetLogger() *logger.Logger {
	if ctx.logger == nil {
		ctx.logger = logger.NewLogger(fmt.Sprintf("%s %s %s", ctx.ip, ctx.method, ctx.path))
	}
	return ctx.logger
}

func GetRequestContext(c *fiber.Ctx) *RequestContext {
	ctx := c.Locals(ReqContextKey{})
	if ctx != nil {
		return ctx.(*RequestContext)
	}
	return nil
}
