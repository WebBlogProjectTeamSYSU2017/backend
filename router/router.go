package router

import (
	"encoding/json"
	"net/http"

	"github.com/WebBlogProjectTeamSYSU2017/backend/server"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

//ResponseHandler 用于返回处理函数接口
type ResponseHandler func(w http.ResponseWriter, r *http.Request) (bool, interface{})

func (h ResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		OK   bool        `json:"ok"`
		Data interface{} `json:"data"`
	}
	ok, ret := h(w, r)
	res := Response{ok, ret}
	byteRes, err := json.Marshal(&res)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(byteRes)
}

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.Handle("/signup", ResponseHandler(server.CreateUserHandler)).Methods("POST")
	mx.Handle("/login", ResponseHandler(server.UserLoginHandler)).Methods("POST")
	mx.Handle("/{email}/createblog", ResponseHandler(server.CreateBlogHandler)).Methods("POST")
	mx.Handle("/{email}/bloghome", ResponseHandler(server.GetAllBlogFromUserHandler)).Methods("GET")
	mx.Handle("/{email}/blogground", ResponseHandler(server.GetAllBlogPublic)).Methods("GET")
	mx.Handle("/{email}/bloghome", ResponseHandler(server.DeleteBlogHandler)).Methods("DELETE")

}
