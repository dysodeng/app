package model

import "github.com/dysodeng/app/internal/infrastructure/persistence/model/common"

type MailConfig struct {
	ID        uint64 `json:"id"`
	User      string `json:"user"`
	FromName  string `json:"from_name"`
	Transport string `json:"transport"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
}

func (c *MailConfig) ToModel() *common.MailConfig {
	dataModel := &common.MailConfig{
		User:      c.User,
		FromName:  c.FromName,
		Transport: c.Transport,
		Username:  c.Username,
		Password:  c.Password,
		Host:      c.Host,
		Port:      c.Port,
	}
	if c.ID > 0 {
		dataModel.ID = c.ID
	}
	return dataModel
}

func MailConfigFormModel(config *common.MailConfig) *MailConfig {
	return &MailConfig{
		ID:        config.ID,
		User:      config.User,
		FromName:  config.FromName,
		Transport: config.Transport,
		Username:  config.Username,
		Password:  config.Password,
		Host:      config.Host,
		Port:      config.Port,
	}
}
