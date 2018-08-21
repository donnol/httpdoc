package httpdoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// MyResponseWriter 实现ResponseWriter接口
type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

// Write 覆盖Write方法
func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}

// Wrap 包装
func Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// From https://stackoverflow.com/questions/43021058/golang-read-request-body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		// 上面将r.Body的内容读了出来，要重新给它赋值
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// 这是一个返回包装器
		mrw := &MyResponseWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
		}

		// 这里将h处理器的返回拦截，将返回内容保存到自定义的结构中
		h.ServeHTTP(mrw, r)

		httpDoc := HTTPDoc{
			Method: r.Method,
			URL:    r.URL.Path,
			Param:  getRequestParam(r, body),
			Return: json.RawMessage(mrw.buf.Bytes()),
		}
		httpDoc.JSONPrint()

		// 为了让请求正常返回，需要将拦截的内容复制到w，让调用方接收
		if _, err := io.Copy(w, mrw.buf); err != nil { // copy之后，buf变为空，应该被释放
			log.Printf("Failed to send out response: %v", err)
		}
	})
}

// HTTPDoc 文档
type HTTPDoc struct {
	Method string
	URL    string
	Param  []string
	Return json.RawMessage
}

// JSONPrint json格式打印
func (h HTTPDoc) JSONPrint() {
	fmt.Println("doc:")
	encodeData, _ := json.Marshal(h)
	buf := new(bytes.Buffer)
	json.Indent(buf, encodeData, "", "\t")
	buf.WriteTo(os.Stdout)
	fmt.Println("")
}

func getRequestParam(r *http.Request, body []byte) []string {
	// 参数
	var params []string
	switch r.Method {
	case http.MethodGet:
		values := r.URL.Query()
		for key := range values {
			params = append(params, key)
		}
	case http.MethodPost:
		// r.ParseForm()
		// values := r.PostForm
		// fmt.Printf("values %+v\n", values)
		// for key := range values {
		// 	httpDoc.Param = append(httpDoc.Param, key)
		// }
		var paramMap = make(map[string]interface{})
		if err := json.Unmarshal(body, &paramMap); err != nil {
			panic(err)
		}
		for key, value := range paramMap {
			_ = value
			// fmt.Printf("%v, %v, %v\n", key, value, interface2String(value))
			params = append(params, key)
		}
	default:
	}

	return params
}

// 将json格式数据转换为markdown
func json2Markdown(data []byte) {

}

// 将json格式数据转换为结构体
func json2Struct(data []byte) {

}

// 接口类型转为字符串
func interface2String(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case float32, float64:
		return "float"
	}
	return ""
}
