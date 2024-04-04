package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"contact-management-service/config"
	"contact-management-service/internal/contacts"
	"contact-management-service/internal/contacts/model"
	"contact-management-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type contactHandler struct {
	s    contacts.Service
	repo contacts.Repo
	cnf  *config.Config
}

func SetupContactHandler(router gin.IRouter, s contacts.Service, repo contacts.Repo, cnf *config.Config) {
	hh := &contactHandler{
		s:    s,
		repo: repo,
		cnf:  cnf,
	}

	router.GET("/health/readiness", hh.ReadinessProbe)
	router.GET("/health/liveness", hh.ReadinessProbe)

	router.GET("/home", func(c *gin.Context) {
		c.File("./templates/index.html")
	})

	endpoints := router.Group("/api/v1")
	endpoints.GET("/contacts", hh.getContacts)
	endpoints.POST("/contacts", hh.addContact)
	endpoints.GET("/contacts/:id", hh.getContact)
	endpoints.PUT("/contacts/:id", hh.updateContact)
	endpoints.DELETE("/contacts/:id", hh.deleteContact)
}

func (hh *contactHandler) ReadinessProbe(c *gin.Context) {
	fmt.Print("Entering health check...")
	err := hh.s.HealthCheck(c)

	if err != nil {
		response.EncodeJSONResp(c, gin.H{
			"status":  "failed",
			"message": err.Error,
		}, http.StatusServiceUnavailable)
		return
	}

	response.EncodeJSONResp(c, gin.H{
		"status": "success",
	}, http.StatusOK)
}

type Contact struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	ExternalID int    `json:"external_id"`
}

func (hh *contactHandler) getContacts(c *gin.Context) {
	data, err := hh.s.GetContact(c.Request.Context(), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	responseData := gin.H{
		"status":  c.Writer.Status(),
		"message": "Success",
		"data":    data,
	}

	response.EncodeJSONResp(c, responseData, http.StatusOK)
}

func (hh *contactHandler) addContact(c *gin.Context) {
	var contact Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseURL := hh.cnf.ExternalAPIConfig.UserApiUrl
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error parsing URL: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing JSON."})
		return
	}

	path := strconv.Itoa(contact.ExternalID)
	u.Path = u.Path + "/" + path

	log.Info().Msg(fmt.Sprintf("User URL: %s", u.String()))

	userData, err := hh.s.FetchExternalUserData(c, u.String())
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to fetch weather data: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather data."})
		return
	}

	timestamp := time.Now()
	data := model.Contacts{
		Timestamp:  timestamp,
		Name:       userData.Name,
		Email:      userData.Email,
		Phone:      userData.Phone,
		ExternalID: contact.ExternalID,
	}

	err = hh.s.AddContact(c, &data)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to store data in postgresql: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to store data in postgresql. - %v", err)})
		return
	}

	responseData := gin.H{
		"status":  http.StatusCreated,
		"message": "Success",
		"data":    data,
	}

	response.EncodeJSONResp(c, responseData, http.StatusCreated)
}

func (hh *contactHandler) getContact(c *gin.Context) {
	userId := c.Params.ByName("id")
	data, err := hh.s.GetContact(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	responseData := gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    data,
	}

	response.EncodeJSONResp(c, responseData, http.StatusOK)
}

func (hh *contactHandler) updateContact(c *gin.Context) {
	id := c.Param("id")
	var contact Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newData := model.Contacts{
		Name:  contact.Name,
		Email: contact.Email,
		Phone: contact.Phone,
	}

	err := hh.s.UpdateContact(c, id, &newData)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to update data in postgresql: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update data in postgresql. - %v", err)})
		return
	}

	responseData := gin.H{
		"status":  http.StatusOK,
		"message": "Contact updated successfully",
	}

	response.EncodeJSONResp(c, responseData, http.StatusOK)
}

func (hh *contactHandler) deleteContact(c *gin.Context) {
	id := c.Param("id")

	err := hh.s.DeleteContact(c.Request.Context(), id)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to deelete data in postgresql: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete data in postgresql. - %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}
