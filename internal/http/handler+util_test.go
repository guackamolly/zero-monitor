package http_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	sut "github.com/guackamolly/zero-monitor/internal/http"
	"github.com/labstack/echo/v4"
)

type TestEchoContext struct {
	echo.Context
	request     *http.Request
	redirectUri string
}

func (ectx TestEchoContext) Request() *http.Request {
	return ectx.request
}

func (ectx *TestEchoContext) Redirect(code int, path string) error {
	ectx.redirectUri = path

	return nil
}

func ContextWithRequest(req *http.Request) *TestEchoContext {
	return &TestEchoContext{
		request: req,
	}
}

func NetURL(rawUrl string) *url.URL {
	u, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}

	return u
}

func TestFromRedirectWithError(t *testing.T) {
	err := errors.New("op failed!")
	errid := sut.StoreHandlerError(err)

	testCases := []struct {
		desc   string
		input  echo.Context
		output error
	}{
		{
			desc: "returns nil error if redirect error query param is not set",
			input: ContextWithRequest(
				&http.Request{URL: NetURL("http://[::]")},
			),
			output: nil,
		},
		{
			desc: "returns error if redirect error query param is contained in errors bucket",
			input: ContextWithRequest(
				&http.Request{URL: NetURL(fmt.Sprintf("http://[::]?x-redirect-err=%s", errid))},
			),
			output: err,
		},
		{
			desc: "returns nil error if redirect error query param is set but errors bucket does not contain it",
			input: ContextWithRequest(
				&http.Request{URL: NetURL(fmt.Sprintf("http://[::]?x-redirect-err=%s", errid))},
			),
			output: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output, _ := sut.FromRedirectWithError(tC.input)
			if output != tC.output {
				t.Errorf("expected %v, but got %v", tC.output, output)
			}
		})
	}
}

func TestRedirectWithError(t *testing.T) {
	err := errors.New("op failed!")

	testCases := []struct {
		desc   string
		input  *TestEchoContext
		output string
	}{
		{
			desc: "appends redirect error query param with error uuid to url",
			input: ContextWithRequest(
				&http.Request{URL: NetURL("http://[::]")},
			),
			output: "x-redirect-err=",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_ = sut.RedirectWithError(tC.input, err)
			output := tC.input.redirectUri

			if !strings.Contains(output, tC.output) {
				t.Errorf("expected to contain %v, but got %v", tC.output, output)
			}
		})
	}
}
