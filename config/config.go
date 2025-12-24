package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hunderaweke/sma-tui/utils"
)

type Room struct {
	PrivateKey   string `json:"private_key,omitempty"`
	PublicKey    string `json:"public_key,omitempty"`
	UniqueString string `json:"unique_string,omitempty"`
	IsPublic     bool   `json:"is_public,omitempty"`
}

type Config struct {
	Rooms           map[string]Room `json:"rooms,omitempty"`
	DefaultIdentity string          `json:"default_room,omitempty"`
}

func New(handler utils.PGPHandler) (*Config, error) {
	key, err := handler.GenerateKey()
	if err != nil {
		return nil, err
	}
	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		return nil, err
	}
	privateKey, err := key.Armor()
	if err != nil {
		return nil, err
	}
	r := Room{
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		UniqueString: key.GetFingerprint()[:12],
		IsPublic:     false,
	}
	var c Config
	c.Rooms = make(map[string]Room)
	c.DefaultIdentity = r.UniqueString
	c.AddRoom(r)
	return &c, nil
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
func (c *Config) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
func (c *Config) AddRoom(r Room) {
	c.Rooms[r.UniqueString] = r
}
func (c *Config) GetRoom(uniqueString string) (*Room, error) {
	r, ok := c.Rooms[uniqueString]
	if !ok {
		return nil, fmt.Errorf("error getting room with unique string: %s", uniqueString)
	}
	return &r, nil
}
