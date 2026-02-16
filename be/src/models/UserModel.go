package models

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type User struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *User) PrepareUser(step string) error {
	if err := u.validateUser(step); err != nil {
		return err
	}

	err := u.formatUser(step)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) validateUser(step string) error {
	if u.Name == "" || u.Nickname == "" || u.Email == "" || (u.Password == "" && step == "Create") {
		return errors.New("Todas as informações são obrigatórias e não podem estar em branco")
	}

	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return err
	}

	return nil
}

func (u *User) formatUser(step string) error {
	u.Name = strings.TrimSpace(u.Name)
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.Email = strings.TrimSpace(u.Email)

	if step == "Create" {
		hashingPassword, err := security.Hash(u.Password)

		if err != nil {
			return err
		}

		u.Password = string(hashingPassword)
	}

	return nil

}
