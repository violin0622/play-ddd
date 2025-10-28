package novel

import (
	"sync"
)

type service struct {
	mu sync.Mutex

	repo Repo
}

func NewService() *service {
	return &service{}
}
