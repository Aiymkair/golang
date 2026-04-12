package v1

import (
	"Assignment-7/internal/entity"
	"Assignment-7/internal/usecase"
	"Assignment-7/pkg/logger"
	"Assignment-7/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRoutes struct {
	t usecase.UserInterface
	l logger.Logger
}

func NewUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Logger) {
	r := &userRoutes{t: t, l: l}

	h := handler.Group("/users")
	{
		// публичные
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)

		// защищённые (JWT)
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		{
			protected.GET("/me", r.GetMe)
		}

		// только для админов
		admin := protected.Group("/")
		admin.Use(utils.RoleMiddleware("admin"))
		{
			admin.PATCH("/promote/:id", r.PromoteUser)
		}
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var dto entity.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		r.l.Error("RegisterUser bind error: " + err.Error())
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	user, err := r.t.RegisterUser(&dto)
	if err != nil {
		r.l.Error("RegisterUser usecase error: " + err.Error())
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var dto entity.LoginUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		r.l.Error("LoginUser bind error: " + err.Error())
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	token, err := r.t.LoginUser(&dto)
	if err != nil {
		r.l.Error("LoginUser usecase error: " + err.Error())
		c.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid user id"})
		return
	}

	user, err := r.t.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse{Error: "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid user id"})
		return
	}

	if err := r.t.PromoteUser(userID); err != nil {
		r.l.Error("PromoteUser error: " + err.Error())
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}
