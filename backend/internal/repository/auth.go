package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Auth struct {
	ID          string    `json:"id,omitempty"`
	Auth        string    `json:"auth"`
	Language    string    `json:"language"`
	Translation string    `json:"translation"`
	AddedAt     time.Time `json:"added_at"`
}

type AuthRepository struct {
	client *elasticsearch.Client
	index  string
}

func NewAuthRepository(client *elasticsearch.Client, index string) *AuthRepository {
	return &AuthRepository{client: client, index: index}
}

func (r *AuthRepository) Login(auth Auth) error {

	auth.AddedAt = time.Now()
	data, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	req := bytes.NewReader(data)
	res, err := r.client.Index(r.index, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	return nil
}
