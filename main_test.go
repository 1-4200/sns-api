package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestResp struct {
	Hits int
	Res  interface{}
}

var (
	url           = "http://localhost:8080/api"
	apiV1         = fmt.Sprintf("%v/%v/", url, "v1")
	startDate     = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	endDate       = time.Now().Format("2006-01-02")
	startDatetime = fmt.Sprintf("%s 00:00", startDate)
	endDatetime   = fmt.Sprintf("%s 23:59", endDate)
	// https://twitter.com/akiko_lawson/
	userID      = "115639376"
	userID2     = "12"
	userID3     = "818664358066548736"
	screenName  = "akiko_lawson"
	domain      = "www.lawson.co.jp"
	mediaType   = "2"
	keyword     = "ニュース"
	count       = "10"
	hashtag     = "天気"
	language    = "24"
	followerMin = "1"
	followerMax = "9999999"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestHealth(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "health/"), nil)
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "{\"message\":\"success\"}", rec.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestTweetsUser(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "missing required params",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/user"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", "")
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/user"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestTweetsUsers(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "missing required params",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/users"), nil)
				params := req.URL.Query()
				params.Add("user_ids", userID)
				params.Add("user_ids", userID2)
				params.Add("user_ids", userID3)
				params.Add("start_date", startDatetime)
				params.Add("end_date", "")
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/users"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("user_ids", userID2)
				params.Add("user_ids", userID3)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestTweetsDomain(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "missing required params",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/domain"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				params.Add("domain", "")
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/domain"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				params.Add("domain", domain)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestTweetsMedia(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "missing required params",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/media"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				params.Add("media_type", "")
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/media"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				params.Add("media_type", mediaType)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestTweetsTransition(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "wrong format for date param",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/transition"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "tweets/transition"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDate)
				params.Add("end_date", endDate)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestHashtags(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "hashtags/"), nil)
				params := req.URL.Query()
				params.Add("keyword", keyword)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				params.Add("count", count)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestHashtagsSearch(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "hashtags/search"), nil)
				params := req.URL.Query()
				params.Add("hashtag", hashtag)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				params.Add("count", count)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestUsersSearch(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "users/search"), nil)
				params := req.URL.Query()
				params.Add("name", screenName)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				params.Add("count", count)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
		{
			name: "ok with params",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "users/search"), nil)
				params := req.URL.Query()
				params.Add("name", screenName)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				params.Add("language", language)
				params.Add("follower_min", followerMin)
				params.Add("follower_max", followerMax)
				params.Add("count", count)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestUsersId(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "ok",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "users/id"), nil)
				params := req.URL.Query()
				params.Add("user_id", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}

func TestUsersIds(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		call func(t *testing.T)
	}{
		{
			name: "GET: ok with multiple user ids",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "users/ids"), nil)
				params := req.URL.Query()
				params.Add("user_ids", userID)
				params.Add("user_ids", userID2)
				params.Add("user_ids", userID3)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
		{
			name: "GET: ok with a user id",
			call: func(t *testing.T) {
				router, _ := setup()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiV1, "users/ids"), nil)
				params := req.URL.Query()
				params.Add("user_ids", userID)
				params.Add("start_date", startDatetime)
				params.Add("end_date", endDatetime)
				req.URL.RawQuery = params.Encode()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)
				var resp TestResp
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Errorf("error=%s", err)
				}
				assert.Equal(t, http.StatusOK, rec.Code)
				if resp.Hits <= 0 {
					t.Errorf("hits = %v, want > 0", resp.Hits)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.call)
	}
}
