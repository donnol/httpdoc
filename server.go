package httpdoc

import (
	"log"
	"net/http"
)

// Param 参数
type Param struct {
	Size   int `json:"size"`
	Offset int `json:"offset"`
}

// Return 返回
type Return struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

// StartServer 启动服务器
func StartServer(addr string, mux *http.ServeMux) {
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
