package database

import (
	"fmt"
	"time"

	"github.com/WebBlogProjectTeamSYSU2017/backend/model"

	"errors"

	"github.com/boltdb/bolt"
)

//DBServiceInterface 定义数据库操作接口
type DBServiceInterface interface {
	GetUserFromEmail(email string) (bool, model.User)
	AddNewSignUpUser(user model.User)
	UserCreateBlog(blog model.Blog) bool
	GetUserAllBlogName(user model.User) []string
	AddNewBlogToUser(blog model.Blog) bool
	GetAllBlogIDFromUserEmail(email string) (bool, []string)
	GetAllBlogData(idlist []string) (bool, []model.Blog)
	GetAllPublicBlog() (bool, []model.Blog)
	DeleteBlogByID(id string, useremail string) error
}

//DBService 接口的实现
type DBService struct{}

//DeleteBlogByID 删除用户自己的博客
func (dbservice *DBService) DeleteBlogByID(id string, useremail string) error {
	//删除博客db内部的信息
	db, err := bolt.Open("blog.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("open failed")
		return err
	} else {
		fmt.Println("open succeed!")
	}
	defer db.Close()
	var blogTitle string
	if err := db.Update(func(tx *bolt.Tx) error {
		blog := tx.Bucket([]byte(id))
		if blog == nil {
			return errors.New("不存在当前id")
		}
		if string(blog.Get([]byte("creatoremail"))) != useremail {
			return errors.New("用户不能删除别人的博客")
		}
		blogTitle = string(blog.Get([]byte("title")))
		if err := tx.DeleteBucket([]byte(id)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	//删除用户db内部的信息
	dbUser, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("open failed")
		return err
	} else {
		fmt.Println("open succeed!")
	}
	defer dbUser.Close()
	if err := dbUser.Update(func(tx *bolt.Tx) error {
		blog := tx.Bucket([]byte(useremail))
		if blog == nil {
			return errors.New("不存在当前用户")
		}
		blogListBucket := blog.Bucket([]byte("bloglist"))
		if string(blogListBucket.Get([]byte(blogTitle))) != id {
			return errors.New("系统错误：当前用户不存在此博客")
		}
		if err := blogListBucket.Delete([]byte(blogTitle)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

//GetAllPublicBlog 获取所有用户的public博客
func (dbservice *DBService) GetAllPublicBlog() (bool, []model.Blog) {
	db, err := bolt.Open("blog.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	var bloglist []model.Blog
	if err != nil {
		fmt.Println("open failed")
		return false, bloglist
	} else {
		fmt.Println("open succeed!")
	}
	defer db.Close()
	if err := db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			if string(b.Get([]byte("ispublic"))) == "1" {
				tempBlogData := model.Blog{
					CreatorEmail: string(b.Get([]byte("creatoremail"))),
					Title:        string(b.Get([]byte("title"))),
					CreateTime:   string(b.Get([]byte("createtime"))),
					Tag:          string(b.Get([]byte("tag"))),
					ID:           string(b.Get([]byte("ID"))),
					IsPublic:     string(b.Get([]byte("ispublic"))),
					Content:      string(b.Get([]byte("content"))),
				}
				bloglist = append(bloglist, tempBlogData)
			}
			return nil
		})
		return nil
	}); err != nil {
		return false, bloglist
	}
	return true, bloglist
}

//GetAllBlogData 从用户名ID的list获取所有博客内容
func (dbservice *DBService) GetAllBlogData(idlist []string) (bool, []model.Blog) {
	db, err := bolt.Open("blog.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	var bloglist []model.Blog
	if err != nil {
		fmt.Println("open failed")
		return false, bloglist
	} else {
		fmt.Println("open succeed!")
	}

	defer db.Close()
	if err := db.View(func(tx *bolt.Tx) error {
		for _, v := range idlist {
			if v == "" {
				continue
			}
			blog := tx.Bucket([]byte(v))
			if blog == nil {
				return errors.New("不存在当前id")
			}
			tempBlogData := model.Blog{
				CreatorEmail: string(blog.Get([]byte("creatoremail"))),
				Title:        string(blog.Get([]byte("title"))),
				CreateTime:   string(blog.Get([]byte("createtime"))),
				Tag:          string(blog.Get([]byte("tag"))),
				ID:           string(blog.Get([]byte("ID"))),
				IsPublic:     string(blog.Get([]byte("ispublic"))),
				Content:      string(blog.Get([]byte("content"))),
			}
			bloglist = append(bloglist, tempBlogData)
		}
		return nil
	}); err != nil {
		return false, bloglist
	}
	return true, bloglist
}

//GetAllBlogIDFromUserEmail 为用户增加从文件名到ID的映射
func (dbservice *DBService) GetAllBlogIDFromUserEmail(email string) (bool, []string) {
	//查找对应用户
	db, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("open failed", err)
		return false, nil
	} else {
		fmt.Println("open succeed!")
	}
	var resultList []string
	defer db.Close()
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(email))
		if b == nil {
			return errors.New("不存在当前用户")
		}
		bList := b.Bucket([]byte("bloglist"))
		bList.ForEach(func(k, v []byte) error {
			//K 、 V 为从博客名 -> 博客ID的映射
			resultList = append(resultList, string(v))
			return nil
		})
		return nil
	}); err != nil {
		return false, resultList
	}
	return true, resultList
}

//AddNewBlogToUser 为用户增加从文件名到ID的映射
func (dbservice *DBService) AddNewBlogToUser(blog model.Blog) bool {
	//查找对应用户
	db, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("open failed", err)
		return false
	} else {
		fmt.Println("open succeed!")
	}
	defer db.Close()
	if err := db.Update(func(tx *bolt.Tx) error {
		blogerOwner := tx.Bucket([]byte(blog.CreatorEmail))
		if blogerOwner == nil {
			return errors.New("不存在此创建者")
		}
		blogList := blogerOwner.Bucket([]byte("bloglist"))
		if blogList == nil {
			return errors.New("该用户不存在博客列表")
		}
		if err := blogList.Put([]byte(blog.Title), []byte(blog.ID)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false
	}
	return true
}

//UserCreateBlog 用户新建博客
func (dbservice *DBService) UserCreateBlog(blog model.Blog) bool {
	db, err := bolt.Open("blog.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("open failed")
		return false
	} else {
		fmt.Println("open succeed!")
	}
	defer db.Close()
	fmt.Println("is going to insert ", blog)
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(blog.ID)); err != nil {
			fmt.Println("create failed")
			return err
		}
		blogBucket := tx.Bucket([]byte(blog.ID))
		if err := blogBucket.Put([]byte("creatoremail"), []byte(blog.CreatorEmail)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		if err := blogBucket.Put([]byte("title"), []byte(blog.Title)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		if err := blogBucket.Put([]byte("createtime"), []byte(blog.CreateTime)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		if err := blogBucket.Put([]byte("tag"), []byte(blog.Tag)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		if err := blogBucket.Put([]byte("ID"), []byte(blog.ID)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		if err := blogBucket.Put([]byte("ispublic"), []byte(string(blog.IsPublic))); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		if err := blogBucket.Put([]byte("content"), []byte(blog.Content)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		return nil
	}); err != nil {
		fmt.Println("插入博客信息失败")
	}
	return true
}

//GetUserAllBlogName 从
func (dbservice *DBService) GetUserAllBlogName(user model.User) []string {
	var blogList []string
	return blogList
}

func worker(user model.User) {
	db, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("open failed")
		return
	} else {
		fmt.Println("open succeed!")
	}
	time.Sleep(1 * time.Second)
	db.Close()
	return
}

//GetUserFromEmail 根据邮箱获取用户信息
func (dbservice *DBService) GetUserFromEmail(email string) (bool, model.User) {
	db, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	var resultUser model.User
	if err != nil {
		fmt.Println("open failed", err)
		return false, model.User{}
	} else {
		fmt.Println("open succeed!")
	}
	defer db.Close()
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(email))
		if b == nil {
			return errors.New("不存在当前用户")
		}

		resultUser = model.User{
			Password: string(b.Get([]byte("password"))),
			Email:    string(b.Get([]byte("email"))),
			Username: string(b.Get([]byte("username"))),
		}
		bList := b.Bucket([]byte("bloglist"))
		bList.ForEach(func(k, v []byte) error {
			//K 、 V 为从博客名 -> 博客ID的映射
			resultUser.BlockList = append(resultUser.BlockList, string(k))
			return nil
		})
		return nil
	}); err != nil {
		return false, resultUser
	}
	return true, resultUser
}

//AddNewSignUpUser 添加博客网站注册用户
func (dbservice *DBService) AddNewSignUpUser(user model.User) {
	db, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("open failed")
		return
	} else {
		fmt.Println("open succeed!")
	}
	defer db.Close()
	fmt.Println("is going to insert ", user)
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(user.Email)); err != nil {
			fmt.Println("create failed")
			return err
		}
		userBucket := tx.Bucket([]byte(user.Email))
		if err := userBucket.Put([]byte("email"), []byte(user.Email)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		if err := userBucket.Put([]byte("username"), []byte(user.Username)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		if err := userBucket.Put([]byte("password"), []byte(user.Password)); err != nil {
			fmt.Println("insert k/v pair failed")
			return err
		}
		blogList, err := userBucket.CreateBucketIfNotExists([]byte("bloglist"))
		if err != nil {
			return err
		}
		if err := blogList.Put([]byte("BLOG HEAD"), []byte("")); err != nil {
			return err
		}

		return nil
	}); err != nil {
		fmt.Println("update failed!")
	}
	return
}
