package bootstrap

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// helper functions
func download(url string) []byte {
	println("> GET %s", url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp := must(client.Get(url))
	if sc := resp.StatusCode; sc != 200 {
		panic(fmt.Sprintf("sc: %d", resp.StatusCode))
	}

	bs := must(io.ReadAll(resp.Body))

	return bs
}

func downloadUnmarshal[T any](url string) T {
	var v T
	must0(json.Unmarshal(download(url), &v))

	return v
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func must0(err error) {
	if err != nil {
		panic(err)
	}
}

func panic(v any) {
	log.Fatal(v)
}

func println(f any, v ...any) {
	if _, ok := f.(string); !ok || len(v) == 0 {
		fmt.Println(f)
		return
	}

	fmt.Printf("%s\n", fmt.Sprintf(f.(string), v))
}
