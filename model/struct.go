package model

//User 定义在数据库中保存的user数据结构
type User struct {
	Username  string
	Password  string
	Email     string
	BlockList []string //保存BlogName
}

type Blog struct {
	ID           string `json:"id"`
	CreatorEmail string `json:"creatoremail"`
	CreateTime   string `json:"createtime"`
	Title        string `json:"title"`
	IsPublic     string `json:"ispublic"`
	Content      string `json:"content"`
	Tag          string `json:"tag"`
}

type Comment struct {
}
