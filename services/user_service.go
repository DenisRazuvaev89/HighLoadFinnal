package services

import (
	"errors"
	"sync"

	"go-microservice/models"
)

type UserService struct {
	users  map[int]*models.User
	mu     sync.RWMutex
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]*models.User),
		nextID: 1,
	}
}

func (s *UserService) GetAll() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		result = append(result, user)
	}
	return result
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) Create(user models.User) *models.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = s.nextID
	s.nextID++
	s.users[user.ID] = &user
	return &user
}

func (s *UserService) Update(id int, user models.User) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return nil, errors.New("user not found")
	}

	user.ID = id
	s.users[id] = &user
	return &user, nil
}

func (s *UserService) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(s.users, id)
	return nil
}
