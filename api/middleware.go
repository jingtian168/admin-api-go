package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jingtian168/admin-api-go/api/e"
	"github.com/jingtian168/admin-api-go/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appg := Gin{C: ctx}
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			appg.AbortWithStatusJSON(http.StatusUnauthorized, e.ErrorAuth, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		authHeaderLen := 2
		if len(fields) < authHeaderLen {
			err := errors.New("invalid authorization header format")
			appg.AbortWithStatusJSON(http.StatusUnauthorized, e.ErrorAuth, err)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			appg.AbortWithStatusJSON(http.StatusUnauthorized, e.ErrorAuth, err)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			appg.AbortWithStatusJSON(http.StatusUnauthorized, e.ErrorAuth, err)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
	
		ctx.Next()
	}
}
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "*")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
