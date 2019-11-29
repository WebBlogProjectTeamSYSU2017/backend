package main

import (
	"github.com/WebBlogProjectTeamSYSU2017/backend/router"
	"os"

	flag "github.com/spf13/pflag"
)

const (
	PORT string = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}
	server := router.NewServer()
	server.Run(":" + port)
	// user := model.User{
	// 	Password: "123456",
	// 	Username: "lauerncelzx",
	// 	Email:    "382673304@qq.com",
	// }
	// db.AddUser(user)
	// // time.Sleep(2 * time.Second)
	// _, userResult := db.GetUser("382673304@qq.com")
	// fmt.Println("-----")
	// fmt.Println(userResult)
}
