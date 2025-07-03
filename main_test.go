package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCafeWhenBad(t *testing.T) {
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
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}
	for _, tc := range requests {
		url := fmt.Sprintf("/cafe?city=moscow&count=%d", tc.count)
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		result := strings.Split(response.Body.String(), ",")
		if tc.want == 0 && len(result[0]) == 0 {
			result = []string{} // если вернулась пустая строка — это 0 записей
		}
		assert.Equal(t, tc.want, len(result), fmt.Sprintf("unexpected cafe count. Count=%d", tc.count))
	}

}
