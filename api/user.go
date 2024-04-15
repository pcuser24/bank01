package api

import (
	"database/sql"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/user2410/simplebank/db/sqlc"
	"github.com/user2410/simplebank/util"
)

type createUserRequest struct {
	Username string                `form:"username" binding:"required,alphanum"`
	Password string                `form:"password" binding:"required,min=6"`
	FullName string                `form:"fullname" binding:"required"`
	Avatar   *multipart.FileHeader `form:"avatar"`
	Email    string                `form:"email" binding:"required,email"`
}

type userResponse struct {
	Username  string    `json:"username"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	Avatar    *string   `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	ur := userResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	if user.Avatar.Valid {
		ur.Avatar = &user.Avatar.String
	}
	return ur
}

// 8mb in bytes
const MAX_AVATAR_SIZE = 1024 * 1024 * 8

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
	}

	// handle upload to S3
	fileHeaders := form.File["avatar"]
	if len(fileHeaders) >= 1 {
		fileHeader := fileHeaders[0]
		file, err := fileHeader.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ftype := fileHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(ftype, "image/") {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "avatar must be an image"})
			return
		}
		ftype = ftype[len("image/"):]
		fsize := fileHeader.Size
		if fsize > MAX_AVATAR_SIZE {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "avatar is too large"})
			return
		}
		avatarUrl, err := server.fileStorage.PutFile(file, fileHeader.Filename, ftype, fsize)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg.Avatar = sql.NullString{
			String: avatarUrl,
			Valid:  avatarUrl != "",
		}
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg.Password = hashedPassword

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			switch dbErr.Code {
			case "23505", "23514":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := newUserResponse(user)

	ctx.JSON(http.StatusOK, res)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	err = util.VerifyPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	})
}
