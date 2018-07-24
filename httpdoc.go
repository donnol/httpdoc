package httpdoc

import (
	"log"
	"net/http"
)

// GenerateHTTPDoc 生成HTTP接口文档
// 对于每个接口，必然要有对应的接口文档。如果全靠个人手写，则必然过于繁琐。哪怕是只改一个字段类型，也要在代码和文档之间同步。久而久之，忘了更新其中一边，就是家常便饭了。
// 在测试接口时自动生成文档，可以大大减少工作量。
// 开发流程变为：接口定义 -> 数据建模 -> 代码开发 -> 功能测试(同时生成文档) -> 代码修改 -> 功能测试(同时生成文档)
func GenerateHTTPDoc(w http.ResponseWriter, r *http.Request) {
	// 收集参数
	var param struct{}
	switch r.Method {
	case http.MethodGet:
		// r.URL.Query()
	case http.MethodPost:
		// r.Body
	default:
		log.Printf("Not support method, %s\n", r.Method)
		return
	}

	// 执行业务
	do := func(struct{}) (result interface{}) { return }
	result := do(param)

	// 收集返回
	_ = result
}

// Wrap 包装
func Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do something

		// 继续执行
		h.ServeHTTP(w, r)
	})
}
