package service

type Store interface {
	SaveURL(fullURL string, shortURL string) error
	GetURL(shortURL string) (string, error)
}

type Service struct {
	s Store
}

const maxRetries = 3

func (s *Service) SaveURL(fullURL string) (string, error) {
	// generate short link (math/rand, crypto/rand) (max length = 8)
	// check that this short link does not yet exist
	// if it's already in the store, regenerate short link
	// max n retries
	// when the short link is unique, save it to store
}

func (s *Service) GetURL(shortURL string) (string, error) {

}
