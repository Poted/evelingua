package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Word struct {
	ID          string    `json:"id,omitempty"`
	Word        string    `json:"word"`
	Language    string    `json:"language"`
	Translation string    `json:"translation"`
	AddedAt     time.Time `json:"added_at"`
}

type WordRepository struct {
	client *elasticsearch.Client
	index  string
}

func NewWordRepository(client *elasticsearch.Client, index string) *WordRepository {
	return &WordRepository{client: client, index: index}
}

func (r *WordRepository) AddWord(word Word) error {
	word.AddedAt = time.Now()
	data, err := json.Marshal(word)
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

func (r *WordRepository) SearchWords(query string, fuzzy bool) ([]Word, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"word": map[string]interface{}{
					"query": query,
				},
			},
		},
	}

	if fuzzy {
		searchQuery["query"].(map[string]interface{})["match"].(map[string]interface{})["word"].(map[string]interface{})["fuzziness"] = "AUTO"
	}

	data, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source Word `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	words := []Word{}
	for _, hit := range result.Hits.Hits {
		words = append(words, hit.Source)
	}

	return words, nil
}
