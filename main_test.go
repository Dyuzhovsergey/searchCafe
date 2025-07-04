package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		fmt.Println(response.Body.String())
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	type requestTest struct {
		request string
		status  int
		message string
	}
	requests := []requestTest{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)
		assert.Equal(t, v.status, response.Code)

		resResp := strings.TrimSpace(response.Body.String())
		assert.Equal(t, v.message, resResp, fmt.Sprintf("unexpected message='%s'", v.message))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: min(100, len(cafeList["moscow"]))},
	}
	for _, tc := range requests {
		url := fmt.Sprintf("/cafe?city=moscow&count=%d", tc.count)
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		result := strings.Split(response.Body.String(), ",")
		if tc.want == 0 && len(result[0]) == 0 {
			result = []string{}
		}
		assert.Len(t, result, tc.want, fmt.Sprintf("unexpected cafe count. Count=%d", tc.count))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	for _, tc := range requests {
		url := fmt.Sprintf("/cafe?city=moscow&search=%s", tc.search)

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		result := strings.Split(response.Body.String(), ",")
		if tc.wantCount == 0 && len(result[0]) == 0 {
			result = []string{}
		}
		assert.Len(t, result, tc.wantCount, fmt.Sprintf("unexpected count cafe for query='%s'", tc.search))

		for _, cafe := range result {
			assert.Contains(
				t,
				strings.ToLower(cafe),
				strings.ToLower(tc.search),
				fmt.Sprintf("cafe name %s does not contain search query %s", cafe, tc.search))
		}

	}
}
