package restserver

import (
	"github.com/Chetan177/pismo/pkg/db/mock"
	"github.com/Chetan177/pismo/pkg/db/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const fail = "fail"

var testServer *RestServer

func startServer() *RestServer {
	s := &RestServer{}
	s.db = mock.NewMockDB(map[string]interface{}{})
	s.config = &Config{port: 2030}
	s.Start()
	return s
}

func TestMain(m *testing.M) {
	testServer = startServer()
	code := m.Run()
	os.Exit(code)
}

func TestCreateAccount(t *testing.T) {
	tables := []struct {
		jsonStrReq   string
		httpResponse int
	}{
		{`{"document_number" : "102392932231231"}`, 200},
		{`{"document_number" : ""}`, 400},
		{`{"document_number" : "102392932231231-fail"}`, 500},
	}

	for i, table := range tables {
		rec := httptest.NewRecorder()
		res := rec.Result()
		defer res.Body.Close()
		req, err := http.NewRequest("post", "/v1/accounts", strings.NewReader(table.jsonStrReq))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)
		c := testServer.server.NewContext(req, rec)
		t.Logf("entering test case: %d", i)
		if assert.NoError(t, testServer.createAccount(c)) {
			assert.Equal(t, table.httpResponse, rec.Code)
		}
		t.Logf("exiting test case: %d", i)
	}

}

func TestGetAccount(t *testing.T) {
	tables := []struct {
		accId        string
		httpResponse int
	}{
		{"65d683625777994656e7dedf", 200},
		{"", 400},
		{"65d683625777994656e7dedf-fail", 500},
	}

	mockData := map[string]interface{}{
		"65d683625777994656e7dedf": &model.Account{DocumentNumber: "102392932231231"},
	}

	testServer.db = mock.NewMockDB(mockData)
	for i, table := range tables {
		rec := httptest.NewRecorder()
		res := rec.Result()
		defer res.Body.Close()
		req, err := http.NewRequest("get", "/v1/accounts/:accId", nil)
		assert.NoError(t, err)
		c := testServer.server.NewContext(req, rec)
		c.SetParamNames("accId")
		c.SetParamValues(table.accId)
		t.Logf("entering test case: %d", i)
		if assert.NoError(t, testServer.getAccount(c)) {
			assert.Equal(t, table.httpResponse, rec.Code)
		}
		t.Logf("exiting test case: %d", i)
	}

}

func TestCreatTransaction(t *testing.T) {
	tables := []struct {
		jsonStrReq   string
		httpResponse int
	}{
		{`{"account_id":"65d683625777994656e7dedf","operation_type_id":2,"amount":20}`, 200},
		{`{"account_id":"","operation_type_id":2,"amount":20}`, 400},
		{`{"account_id":"65d683625777994656e7dedf","operation_type_id":10,"amount":20}`, 400},
		{`{"account_id":"65d683625777994656e7dedf-fail","operation_type_id":2,"amount":20}`, 400},
	}

	mockData := map[string]interface{}{
		"65d683625777994656e7dedf": &model.Account{DocumentNumber: "102392932231231"},
	}

	testServer.db = mock.NewMockDB(mockData)
	for i, table := range tables {
		rec := httptest.NewRecorder()
		res := rec.Result()
		defer res.Body.Close()
		req, err := http.NewRequest("post", "/v1/transactions", strings.NewReader(table.jsonStrReq))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)
		c := testServer.server.NewContext(req, rec)
		t.Logf("entering test case: %d", i)
		if assert.NoError(t, testServer.createTransaction(c)) {
			assert.Equal(t, table.httpResponse, rec.Code)
		}
		t.Logf("exiting test case: %d", i)
	}
}
