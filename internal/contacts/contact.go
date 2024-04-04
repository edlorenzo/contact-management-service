package contacts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"contact-management-service/internal/contacts/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type (
	contactService struct {
		repo Repo
		db   *gorm.DB
	}
)

func NewContactService(
	repo Repo,
	db *gorm.DB,
) Service {
	return &contactService{
		repo: repo,
		db:   db,
	}
}

func (c *contactService) HealthCheck(ctx context.Context) error {
	return c.repo.HealthCheck(ctx)
}

func (c *contactService) FetchExternalUserData(ctx context.Context, url string) (*UserExternalAPIResponse, error) {
	var rs UserExternalAPIResponse
	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Debug().Msg(fmt.Sprintf("error: %s", resp.Status))
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
		return nil, err
	}

	return &rs, nil
}

func (c *contactService) AddContact(ctx context.Context, user *model.Contacts) error {
	transactionFn := func(tx *gorm.DB) error {
		err := c.repo.StoreDataInDB(ctx, tx, user)
		if err != nil {
			log.Debug().Msg(fmt.Sprintf("insert data failed: %v", err))
			return err
		}
		return nil
	}

	err := c.repo.RunTransaction(ctx, transactionFn, func() error {
		log.Debug().Msg("failed inserting of user data")
		return errors.New("failed inserting of user data")
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *contactService) GetContact(ctx context.Context, userId string) (*[]contactResponse, error) {
	var response []contactResponse
	rs, err := c.repo.GetContactByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	for i, r := range rs {
		if i >= len(response) {
			response = append(response, contactResponse{})
		}
		response[i].ID = int(r.ID)
		response[i].Name = r.Name
		response[i].Phone = r.Phone
		response[i].Email = r.Email
		response[i].ExternalID = r.ExternalID
		response[i].Timestamp = r.Timestamp
		response = append(response)
	}

	return &response, err
}

func (c *contactService) UpdateContact(ctx context.Context, userId string, user *model.Contacts) error {
	transactionFn := func(tx *gorm.DB) error {
		err := c.repo.UpdateDataInDB(ctx, tx, userId, user)
		if err != nil {
			log.Debug().Msg(fmt.Sprintf("update data failed: %v", err))
			return err
		}
		return nil
	}

	err := c.repo.RunTransaction(ctx, transactionFn, func() error {
		log.Debug().Msg("failed updating of user data")
		return errors.New("failed updating of user data")
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *contactService) DeleteContact(ctx context.Context, userId string) error {
	transactionFn := func(tx *gorm.DB) error {
		err := c.repo.DeleteContact(ctx, tx, userId)
		if err != nil {
			log.Debug().Msg(fmt.Sprintf("delete data failed: %v", err))
			return err
		}
		return nil
	}

	err := c.repo.RunTransaction(ctx, transactionFn, func() error {
		log.Debug().Msg("failed deleting of user data")
		return errors.New("failed deleting of user data")
	})

	if err != nil {
		return err
	}

	return nil
}
