package database

import "github.com/brightside-dev/dev-blog/database/client"

func New() client.DatabaseService {
	return client.NewMySQL()
}
