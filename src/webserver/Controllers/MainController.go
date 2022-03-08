package controllers

import "github.com/antonPalmFolkmann/DevOps2022/models"

type Response struct {
	Data []models.User `json:"data"`
	Message string `json:"message"`
}