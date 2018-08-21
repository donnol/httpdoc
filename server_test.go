package httpdoc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

var addr = ":8080"

func TestMain(m *testing.M) {
	go func() {
		mux := http.NewServeMux()
		getHandler := http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				_ = q
				// fmt.Printf("param is %+v\n", q)

				result := Return{
					Total: 10,
					List: []struct {
						Name string `json:"name"`
					}{
						{Name: "jd"},
					},
				}
				data, err := json.Marshal(result)
				if err != nil {
					panic(err)
				}

				w.Write(data)
			})
		mux.Handle("/get", Wrap(getHandler))
		addHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var param struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}
			if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
				panic(err)
			}
			// fmt.Printf("param is %+v\n", param)

			result := map[string]int{
				"id": 1,
			}

			data, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}

			w.Write(data)
		})
		mux.Handle("/add", Wrap(addHandler))
		StartServer(addr, mux)
	}()

	m.Run()
}

func TestStartServer(t *testing.T) {
	t.Run("/get", func(t *testing.T) {
		var u url.URL
		u.Host = "localhost" + addr
		u.Scheme = "http"
		u.Path = "/get"
		q := u.Query()
		q.Set("size", "10")
		q.Set("offset", "0")
		encq := q.Encode()
		u.RawQuery = encq

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			t.Fatal("bad response")
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		_ = data
		// fmt.Println("returns:")
		// var buf = new(bytes.Buffer)
		// json.Indent(buf, data, "", "\t")
		// buf.WriteTo(os.Stdout)
		// fmt.Println("")
	})

	t.Run("/add", func(t *testing.T) {
		var u url.URL
		u.Host = "localhost" + addr
		u.Scheme = "http"
		u.Path = "/add"
		var body = new(bytes.Buffer)
		body.WriteString(`{"name": "jd","age": 20}`)
		req, err := http.NewRequest(http.MethodPost, u.String(), body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			t.Fatal("bad response")
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		_ = data
		// fmt.Println("returns:")
		// var buf = new(bytes.Buffer)
		// json.Indent(buf, data, "", "\t")
		// buf.WriteTo(os.Stdout)
		// fmt.Println("")
	})
}
