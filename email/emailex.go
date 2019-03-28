package email

import (
	"encoding/json"
	"io/ioutil"
	"net/mail"
	"net/smtp"
)

func SendMessage(config *Config, m *Message) error {
	m.From = mail.Address{Name: config.FromName, Address: config.FromAddress}
	m.To = config.To
	return Send(config.Smtp+":"+config.Port, getAuth(config), m)
}

func SendMessageTo(config *Config, m *Message, toEmails []string) error {
	m.From = mail.Address{Name: config.FromName, Address: config.FromAddress}
	m.To = toEmails
	return Send(config.Smtp+":"+config.Port, getAuth(config), m)
}

func getAuth(config *Config) smtp.Auth {
	return smtp.PlainAuth("", config.FromAddress, config.Password, config.Smtp)
}

func LoadConfig(filename string) (*Config, error) {
	cfg := &Config{}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	return cfg, json.Unmarshal(b, cfg)
}

type Config struct {
	FromName    string   `json:"name"`
	FromAddress string   `json:"address"`
	Password    string   `json:"password"`
	Smtp        string   `json:"smtp"`
	Port        string   `json:"port"`
	To          []string `json:"to"`
}
