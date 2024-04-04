package contacts

import (
	"context"
	"time"

	"contact-management-service/internal/contacts/model"
	"gorm.io/gorm"
)

type (
	Service interface {
		HealthCheck(ctx context.Context) error
		FetchExternalUserData(ctx context.Context, url string) (*UserExternalAPIResponse, error)
		AddContact(ctx context.Context, user *model.Contacts) error
		GetContact(ctx context.Context, userId string) (*[]contactResponse, error)
		UpdateContact(ctx context.Context, userId string, user *model.Contacts) error
		DeleteContact(ctx context.Context, userId string) error
	}

	Repo interface {
		StoreDataInDB(ctx context.Context, tx *gorm.DB, user *model.Contacts) error
		GetContactByID(ctx context.Context, userId string) ([]model.Contacts, error)
		UpdateDataInDB(ctx context.Context, tx *gorm.DB, userId string, user *model.Contacts) error
		DeleteContact(ctx context.Context, tx *gorm.DB, userId string) error
		RunTransaction(ctx context.Context, fn func(tx *gorm.DB) error, onCommitFail func() error) error
		HealthCheck(ctx context.Context) error
	}
)

type (
	UserExternalAPIResponse struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Address  struct {
			Street  string `json:"street"`
			Suite   string `json:"suite"`
			City    string `json:"city"`
			Zipcode string `json:"zipcode"`
			Geo     struct {
				Lat string `json:"lat"`
				Lng string `json:"lng"`
			} `json:"geo"`
		} `json:"address"`
		Phone   string `json:"phone"`
		Website string `json:"website"`
		Company struct {
			Name        string `json:"name"`
			CatchPhrase string `json:"catchPhrase"`
			Bs          string `json:"bs"`
		} `json:"company"`
	}

	contactResponse struct {
		ID         int       `json:"ID"`
		Timestamp  time.Time `json:"Timestamp"`
		Name       string    `json:"Name"`
		Email      string    `json:"Email"`
		Phone      string    `json:"Phone"`
		ExternalID int       `json:"ExternalID"`
	}
)
