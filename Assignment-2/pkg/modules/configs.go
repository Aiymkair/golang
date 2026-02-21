package modules

import "time"

type PostgreSQL struct {
	Host        string
	Port        string
	Username    string
	Password    string
	DBName      string
	SSLMode     string
	ExecTimeout time.Duration
}
