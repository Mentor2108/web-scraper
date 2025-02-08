package user

import (
	"net/http"
	"backend-service/defn"

	"github.com/julienschmidt/httprouter"
)

type UserRoutesHandler struct {
	service defn.Service
}

func NewUserRoutesHandler(service defn.Service) *UserRoutesHandler {
	return &UserRoutesHandler{
		service: service,
	}
}

func AddRoutes(router *httprouter.Router, userHandler *UserRoutesHandler) {
	router.Handle(http.MethodPost, "/signup", userHandler.Create)
	// router.Handle()
	// router.Handle(http.MethodPost, "/signup", Signup)
	// router.Handle(http.MethodGet, "/portfolio", RetrievePortfolio)
}
