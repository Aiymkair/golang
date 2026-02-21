package handler

import (
	"assignment-2/internal/usecase"
)

type Handler struct {
	UserUsecase usecase.UserUsecase
}

func NewHandler(uu usecase.UserUsecase) *Handler {
	return &Handler{UserUsecase: uu}
}
