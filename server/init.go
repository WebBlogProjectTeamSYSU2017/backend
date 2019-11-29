package server

import (
	"github.com/WebBlogProjectTeamSYSU2017/backend/database"
)

var dbServer database.DBServiceInterface

func init() {
	dbServer = &database.DBService{}
}
