package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/WebBlogProjectTeamSYSU2017/backend/lib"
	"github.com/WebBlogProjectTeamSYSU2017/backend/model"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

//CreateUserRequest 用于创建网站用户请求
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

//DeleteBlogRequest 用于用户删除自己的博客
type DeleteBlogRequest struct {
	ID string `json:"id"`
}

//UserLogInRequest 用于用户登录请求
type UserLogInRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

//CreateBlogRequest 用户创建博客请求
type CreateBlogRequest struct {
	Title    string `json:"title"`
	IsPublic string `json:"ispublic"`
	Content  string `json:"content"`
	Tag      string `json:"tag"`
}

type PublicBlogsResponse struct {
	ID          string `json:"id"`
	CreatorName string `json:"creatorname"`
	CreateTime  string `json:"createtime"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Tag         string `json:"tag"`
}

//AddHandler 用于测试
func AddHandler(w http.ResponseWriter, req *http.Request) {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	vars := mux.Vars(req)
	id1 := vars["id1"]
	id2 := vars["id2"]
	id_1, _ := strconv.Atoi(id1)
	id_2, _ := strconv.Atoi(id2)
	id := id_1 + id_2
	id_str := strconv.Itoa(id)
	formatter.JSON(w, http.StatusOK, struct{ Test string }{id1 + " + " + id2 + " = " + id_str})
}

//DeleteBlogHandler 提供删除博客服务
func DeleteBlogHandler(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	body, err := ioutil.ReadAll(req.Body)
	vars := mux.Vars(req)
	useremail := vars["email"]
	if ok, _ := dbServer.GetUserFromEmail(useremail); ok != true {
		return false, "用户不存在"
	}
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	deleteBlogRequest := DeleteBlogRequest{}
	if err := json.Unmarshal(body, &deleteBlogRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}

	if deleteBlogRequest.ID == "" {
		log.Print("博客ID内容不能为空")
		return false, "博客名、博客内容不能为空"
	}

	if err := dbServer.DeleteBlogByID(deleteBlogRequest.ID, useremail); err != nil {
		return false, fmt.Sprintf("删除博客失败,错误信息为：%s", err)
	}
	return true, ""
}

//GetAllBlogPublic 获取当前所有public博客
func GetAllBlogPublic(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	vars := mux.Vars(req)
	useremail := vars["email"]
	if ok, _ := dbServer.GetUserFromEmail(useremail); ok != true {
		return false, "用户不存在"
	}
	ok, allBlogData := dbServer.GetAllPublicBlog()
	if ok != true {
		return false, "获取所有公开博客失败"
	}
	fmt.Println(allBlogData)

	var publicBlogResponseList []PublicBlogsResponse

	for _, v := range allBlogData {
		ok, blogUser := dbServer.GetUserFromEmail(v.CreatorEmail)
		if ok != true {
			return false, "系统错误：存在无效的博客Email"
		}
		publicBlogResponseList = append(publicBlogResponseList, PublicBlogsResponse{
			ID:          v.ID,
			CreateTime:  v.CreateTime,
			Title:       v.Title,
			Content:     v.Content,
			Tag:         v.Tag,
			CreatorName: blogUser.Username,
		})
	}
	return true, publicBlogResponseList
}

//GetAllBlogFromUserHandler 获取当前用户所有的博客
func GetAllBlogFromUserHandler(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	vars := mux.Vars(req)
	useremail := vars["email"]
	if ok, _ := dbServer.GetUserFromEmail(useremail); ok != true {
		return false, "用户不存在"
	}

	ok, allBlogID := dbServer.GetAllBlogIDFromUserEmail(useremail)
	if ok != true {
		return false, "获取所有博客ID失败"
	}
	fmt.Println(allBlogID)
	ok, allBlogData := dbServer.GetAllBlogData(allBlogID)
	if ok != true {
		return false, "从ID获取博客数据失败"
	}
	return true, allBlogData
}

//CreateBlogHandler 提供创建博客服务
func CreateBlogHandler(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	body, err := ioutil.ReadAll(req.Body)
	vars := mux.Vars(req)
	useremail := vars["email"]
	if ok, _ := dbServer.GetUserFromEmail(useremail); ok != true {
		return false, "用户不存在"
	}
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	crateBlogRequest := CreateBlogRequest{}
	if err := json.Unmarshal(body, &crateBlogRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}

	if crateBlogRequest.Content == "" || crateBlogRequest.Title == "" {
		log.Print("博客名、博客内容不能为空")
		return false, "博客名、博客内容不能为空"
	}
	blog := model.Blog{
		CreatorEmail: useremail,
		Title:        crateBlogRequest.Title,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		Tag:          crateBlogRequest.Tag,
		ID:           lib.GetUniqueID(),
		IsPublic:     crateBlogRequest.IsPublic,
		Content:      crateBlogRequest.Content,
	}

	if ok := dbServer.UserCreateBlog(blog); ok != true {
		return false, "创建博客失败"
	}
	if ok := dbServer.AddNewBlogToUser(blog); ok != true {
		return false, "用户插入博客失败"
	}
	return true, ""
}

//CreateUserHandler 提供创建用户服务
func CreateUserHandler(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	createUserRequest := CreateUserRequest{}
	if err := json.Unmarshal(body, &createUserRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}

	// check if given information is valid
	if ok := lib.CheckEmail(createUserRequest.Email); !ok {
		return false, "无效的邮箱地址"
	}
	if ok := lib.CheckUsername(createUserRequest.Username); !ok {
		return false, "无效的用户名"
	}

	user := model.User{
		Username: createUserRequest.Username,
		Password: createUserRequest.Password,
		Email:    createUserRequest.Email,
	}
	dbServer.AddNewSignUpUser(user)
	return true, ""
}

//UserLoginHandler 提供用户登录服务
func UserLoginHandler(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	userLogInRequest := UserLogInRequest{}
	if err := json.Unmarshal(body, &userLogInRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}

	// check if given information is valid
	if ok := lib.CheckEmail(userLogInRequest.Email); !ok {
		return false, "无效的邮箱地址"
	}

	ok, user := dbServer.GetUserFromEmail(userLogInRequest.Email)
	if !ok {
		return false, "该用户不存在"
	}
	if user.Password != userLogInRequest.Password {
		return false, "密码错误"
	}
	return true, ""
}
