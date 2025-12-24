package app

import (
	"errors"
	"os"

	"github.com/hunderaweke/sma-tui/config"
	"github.com/hunderaweke/sma-tui/utils"
)

type App struct {
	Config     *config.Config
	pgpHandler *utils.PGPHandler
}

func NewApp() (*App, error) {
	var a App
	a.pgpHandler = utils.NewPGPHandler()
	c, err := config.Load(os.Getenv("CONFIG_PATH"))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		c, err = config.New(*a.pgpHandler)
		if err != nil {
			return nil, err
		}
		err = c.Save(os.Getenv("CONFIG_PATH"))
		if err != nil {
			return nil, err
		}
	}
	a.Config = c
	return &a, nil
}
