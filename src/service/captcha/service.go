package captcha

import (
	"github.com/julienschmidt/httprouter"
	"github.com/k42-software/go-altcha"
	"math/rand"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/repository/t"
	"sync"
	"time"
)

type Service struct {
	signatures map[string]time.Time
	mutex      *sync.RWMutex
}

func NewService() *Service {
	return &Service{
		signatures: map[string]time.Time{},
		mutex:      &sync.RWMutex{},
	}
}

func (s *Service) GetChallenge(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	api.Success(w, r, []byte(s.Challenge()))
}

func (s *Service) Challenge() string {
	challenge := altcha.NewChallengeWithParams(altcha.Parameters{})
	s.mutex.Lock()
	s.signatures[challenge.Signature] = time.Now().Add(2 * time.Minute)
	if rand.Intn(1000) == 0 {
		now := time.Now()
		for key, t := range s.signatures {
			if t.Before(now) {
				delete(s.signatures, key)
			}
		}
	}
	s.mutex.Unlock()
	return challenge.Encode()
}

func (s *Service) Solve(response string) error {
	msg, err := altcha.DecodeResponse(response)
	if err != nil {
		return t.Errorf("can't decode response")
	}
	s.mutex.RLock()
	created, ok := s.signatures[msg.Signature]
	s.mutex.RUnlock()
	if !ok {
		return t.Errorf("challenge not found")
	}
	s.mutex.Lock()
	delete(s.signatures, msg.Signature)
	s.mutex.Unlock()
	if created.Before(time.Now()) {
		return t.Errorf("challenge too old")
	}
	ok = msg.IsValidResponse()
	if !ok {
		return t.Errorf("response invalid")
	}
	return nil
}
