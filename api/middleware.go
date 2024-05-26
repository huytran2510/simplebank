package api

import (
	"errors"
	"net/http"
	"simplebank/token"
	"strings"

	"github.com/gin-gonic/gin"
)


const (
	authorizationHeaderKey  = "authorization" 
	authorizationTypeBearer = "bearer"
	authorizationPayload_key = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
        authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		field := strings.Fields(authorizationHeader)
		if len(field) < 2 {
            err := errors.New("invalid authorization format")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
            return
        }

		authorizationType := strings.ToLower(field[0])
		if authorizationType!= authorizationTypeBearer {
            err := errors.New("invalid authorization type")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
            return
        }

		accessToken := field[1]
		payload,err := tokenMaker.VerifyToken(accessToken)
		if err!= nil {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
            return
        }

		ctx.Set(authorizationPayload_key,payload)
		ctx.Next()
    }
}