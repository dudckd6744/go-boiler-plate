package user

import (
	"database/sql"
	"fmt"
)

type UserService struct {
	repository *UserRepository
}

var Service UserService

func (s *UserService) InitService(db *sql.DB) {

	s.repository = &Repository
	s.repository.AssignDB(db)
}

func (s *UserService) CreateUser(payload UserDto) (sql.Result, error) {
	raw, err := s.repository.CreateUser(payload)

	return raw, err
}

func (s *UserService) FindAllUser() (*[]UserRaw, error) {

	result, err := s.repository.FindAllUser()

	return result, err
}

func (s *UserService) FindDetailUser(id int) (*UserRaw, error) {

	result, err := s.repository.FindDetailUser(id)

	return result, err
}

func (s *UserService) PatchUserName(id int, name *string) (sql.Result, error) {
	fmt.Println(&name)

	result, err := s.repository.PatchUserName(id, name)

	return result, err
}
