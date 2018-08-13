package httpdoc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

var addr = ":8080"

func TestMain(m *testing.M) {
	go func() {
		mux := http.NewServeMux()
		handler := http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				log.Printf("%+v\n", q)

				result := Return{
					Total: 10,
					List:  []string{"jd"},
				}
				data, err := json.Marshal(result)
				if err != nil {
					panic(err)
				}
				w.Write(data)
			})
		mux.Handle("/", Wrap(handler))
		StartServer(addr, mux)
	}()

	m.Run()
}

func TestStartServer(t *testing.T) {
	var u url.URL
	u.Host = "localhost" + addr
	u.Scheme = "http"
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
	var buf = new(bytes.Buffer)
	json.Indent(buf, data, "", "\t")
	buf.WriteTo(os.Stdout)
}
