package router

import (
	"contact-management-service/config"
	svc "contact-management-service/internal/contacts"
	"contact-management-service/internal/contacts/handler"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine, s svc.Service, repo svc.Repo, cnf *config.Config) {
	handler.SetupContactHandler(r, s, repo, cnf)
}
