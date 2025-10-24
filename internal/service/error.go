package service

import "errors"

var (
	ErrNewsNotFound            = errors.New("news not found")
	ErrCategoriesAlreadyExists = errors.New("categories already exists")

	ErrNotFound      = errors.Join(ErrNewsNotFound)
	ErrAlreadyExists = errors.Join(ErrCategoriesAlreadyExists)
)
