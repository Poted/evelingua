package service

import "evelinqua/internal/repository"

type WordService struct {
	repo *repository.WordRepository
}

func NewWordService(repo *repository.WordRepository) *WordService {
	return &WordService{repo: repo}
}

func (s *WordService) AddWord(word repository.Word) error {
	return s.repo.AddWord(word)
}

func (s *WordService) SearchWords(query string, fuzzy bool) ([]repository.Word, error) {
	return s.repo.SearchWords(query, fuzzy)
}
