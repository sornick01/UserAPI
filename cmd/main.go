package main

import (
	"context"
	"fmt"
	"github.com/sornick01/UserAPI/internal/user/repository"
	"github.com/sornick01/UserAPI/models"
	"log"
)

func main() {
	repo := repository.NewJsonRepo("users.json")
	s, err := repo.CreateUser(context.Background(), &models.User{
		DisplayName: "mpeanuts",
		Email:       "abc@mail.ru",
	})

	if err != nil {
		log.Fatal(err)
	}

	s, err = repo.CreateUser(context.Background(), &models.User{
		DisplayName: "qwerty",
		Email:       "hello@gmail.com",
	})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}
