package pkg

import "errors"

var ErrClientIsInstalled = errors.New("client is already installed")
var ErrConfigAlreadyExists = errors.New("Config already exists")
var ErrValidatorAlreadyExists = errors.New("validator address already exists")
