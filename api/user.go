package api

import (
	"database/sql"
	"net/http"
	"time"
    "log"
	"github.com/jingtian168/admin-api-go/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jingtian168/admin-api-go/api/e"
	db "github.com/jingtian168/admin-api-go/db/sqlc"
	"github.com/jingtian168/admin-api-go/util"
	"github.com/lib/pq"

	
)


type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum" lable:"用户名"`
	Password string `json:"password" binding:"required,min=6" lable:"密码"`
	FullName string `json:"full_name" binding:"required" lable:"姓名"`
	Email    string `json:"email" binding:"required_without=Phone,omitempty,email" lable:"邮箱"`
	Phone    string `json:"phone" binding:"required_without=Email,omitempty,phone" lable:"手机号码"`
}

type userResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		Phone:            user.Phone,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	appG := Gin{C: ctx}
	var req createUserRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidPrarms, err)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.Error, err)
		return
	}

	// TODO 检查 邮箱和手机号唯一
	
	id, err := uuid.NewRandom()
	arg := db.CreateUserParams{
		ID:             id,
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				appG.Response(http.StatusForbidden, e.ErrorUsernameExit, nil)
				return
			}
		}
		appG.Response(http.StatusInternalServerError, e.Error, err)
		return
	}

	resp := newUserResponse(user)
	appG.Response(http.StatusOK, e.Success, resp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`

}

func (server *Server) loginUser(ctx *gin.Context) {
	appg := Gin{C: ctx}
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		appg.Response(http.StatusBadRequest, e.InvalidPrarms, err)
		return
	}

	user, err := server.store.GetUserByName(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			appg.Response(http.StatusNotFound, e.NotFound, gin.H{})
			return
		}
		appg.Response(http.StatusInternalServerError, e.Error, err)
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		appg.Response(http.StatusUnauthorized, e.Unauthorized, err)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		appg.Response(http.StatusInternalServerError, e.Error, err)
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		appg.Response(http.StatusInternalServerError, e.Error, err)
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		appg.Response(http.StatusInternalServerError, e.Error, err)
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}
	appg.Response(http.StatusOK, e.Success, rsp)
}
type userInfoResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
}

func (server *Server) getUserInfo(ctx *gin.Context) {	
	appG := Gin{C: ctx}
    payload  := ctx.MustGet("authorization_payload")
   
   
    id := payload.(*token.Payload).UserID
	log.Println(id)
	user, err := server.store.GetUser(ctx, id)
	if err != nil {	
		appG.Response(http.StatusInternalServerError, e.Error, err)
		return
	}

	rep := userInfoResponse{
		Username:        user.Username,
		FullName:         user.FullName,
		Email:            	user.Email,
		Phone:            	user.Phone,	}	
	appG.Response(http.StatusOK, e.Success, rep)
    
}
