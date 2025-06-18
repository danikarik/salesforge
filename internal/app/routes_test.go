package app

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danikarik/salesforge/internal/model"
	"github.com/danikarik/salesforge/internal/model/mock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	mocky "github.com/stretchr/testify/mock"
)

func performRequest(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestCreateSequence(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("CreateSequence", mocky.Anything, &model.Sequence{
			Name:                 "Test Sequence",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: false,
			Steps: []*model.Step{
				{Subject: "Step 1", Content: "Content 1"},
				{Subject: "Step 2", Content: "Content 2"},
			},
		}).Return(nil)

		service := NewService(Config{Store: store})

		req := `{
			"name": "Test Sequence",
			"openTrackingEnabled": true,
			"clickTrackingEnabled": false,
			"steps": [
				{"subject": "Step 1", "content": "Content 1"},
				{"subject": "Step 2", "content": "Content 2"}
			]
		}`

		w := performRequest(service.Handler(), "POST", "/sequences", req)
		assert.Equal(t, 201, w.Code)
	})

	t.Run("FailedValidation", func(t *testing.T) {
		store := &mock.MockStore{}
		service := NewService(Config{Store: store})

		req := `{
			"name": "Test Sequence",
			"openTrackingEnabled": true,
			"clickTrackingEnabled": false
		}`

		w := performRequest(service.Handler(), "POST", "/sequences", req)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("WithError", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("CreateSequence", mocky.Anything, &model.Sequence{
			Name:                 "Test Sequence",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: false,
			Steps: []*model.Step{
				{Subject: "Step 1", Content: "Content 1"},
				{Subject: "Step 2", Content: "Content 2"},
			},
		}).Return(errors.New("creation failed"))

		service := NewService(Config{Store: store})

		req := `{
			"name": "Test Sequence",
			"openTrackingEnabled": true,
			"clickTrackingEnabled": false,
			"steps": [
				{"subject": "Step 1", "content": "Content 1"},
				{"subject": "Step 2", "content": "Content 2"}
			]
		}`

		w := performRequest(service.Handler(), "POST", "/sequences", req)
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceCreationFailed.Error()))
	})
}

func TestFetchSequence(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("FetchSequence", mocky.Anything, uint64(1)).Return(&model.Sequence{
			ID:                   1,
			Name:                 "Test Sequence",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: false,
			Steps: []*model.Step{
				{ID: 1, Subject: "Step 1", Content: "Content 1"},
				{ID: 2, Subject: "Step 2", Content: "Content 2"},
			},
		}, nil)

		service := NewService(Config{Store: store})

		w := performRequest(service.Handler(), "GET", "/sequences/1", "")
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Test Sequence"`)
	})

	t.Run("NotFound", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("FetchSequence", mocky.Anything, uint64(1)).Return(nil, pgx.ErrNoRows)

		service := NewService(Config{Store: store})

		w := performRequest(service.Handler(), "GET", "/sequences/1", "")
		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceNotFound.Error()))
	})

	t.Run("WithError", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("FetchSequence", mocky.Anything, uint64(1)).Return(nil, errors.New("fetch failed"))

		service := NewService(Config{Store: store})

		w := performRequest(service.Handler(), "GET", "/sequences/1", "")
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceFetchingFailed.Error()))
	})
}

func TestUpdateSequence(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("UpdateSequence", mocky.Anything, uint64(1), &model.Sequence{
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: false,
		}).Return(nil)

		service := NewService(Config{Store: store})

		req := `{
			"openTrackingEnabled": true,
			"clickTrackingEnabled": false
		}`

		w := performRequest(service.Handler(), "PUT", "/sequences/1", req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("UpdateSequence", mocky.Anything, uint64(1), mocky.Anything).Return(pgx.ErrNoRows)

		service := NewService(Config{Store: store})

		req := `{
			"openTrackingEnabled": true,
			"clickTrackingEnabled": false
		}`

		w := performRequest(service.Handler(), "PUT", "/sequences/1", req)
		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceNotFound.Error()))
	})

	t.Run("WithError", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("UpdateSequence", mocky.Anything, uint64(1), mocky.Anything).Return(errors.New("update failed"))

		service := NewService(Config{Store: store})

		req := `{
			"openTrackingEnabled": true,
			"clickTrackingEnabled": false
		}`

		w := performRequest(service.Handler(), "PUT", "/sequences/1", req)
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceUpdateFailed.Error()))
	})
}

func TestUpdateStep(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("UpdateStep", mocky.Anything, uint64(1), &model.Step{
			SequenceID: 1,
			Subject:    "Updated Step",
			Content:    "Updated Content",
		}).Return(nil)

		service := NewService(Config{Store: store})

		req := `{
			"subject": "Updated Step",
			"content": "Updated Content"
		}`

		w := performRequest(service.Handler(), "PUT", "/sequences/1/steps/1", req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("UpdateStep", mocky.Anything, uint64(1), mocky.Anything).Return(pgx.ErrNoRows)

		service := NewService(Config{Store: store})

		req := `{
			"subject": "Updated Step",
			"content": "Updated Content"
		}`

		w := performRequest(service.Handler(), "PUT", "/sequences/1/steps/1", req)
		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceNotFound.Error()))
	})

	t.Run("WithError", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("UpdateStep", mocky.Anything, uint64(1), mocky.Anything).Return(errors.New("update failed"))

		service := NewService(Config{Store: store})

		req := `{
			"subject": "Updated Step",
			"content": "Updated Content"
		}`

		w := performRequest(service.Handler(), "PUT", "/sequences/1/steps/1", req)
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceUpdateFailed.Error()))
	})
}

func TestDeleteStep(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("DeleteStep", mocky.Anything, uint64(1), &model.Step{SequenceID: 1}).Return(nil)

		service := NewService(Config{Store: store})

		w := performRequest(service.Handler(), "DELETE", "/sequences/1/steps/1", "")
		assert.Equal(t, 204, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("DeleteStep", mocky.Anything, uint64(1), mocky.Anything).Return(pgx.ErrNoRows)

		service := NewService(Config{Store: store})

		w := performRequest(service.Handler(), "DELETE", "/sequences/1/steps/1", "")
		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceNotFound.Error()))
	})

	t.Run("WithError", func(t *testing.T) {
		store := &mock.MockStore{}
		store.On("DeleteStep", mocky.Anything, uint64(1), mocky.Anything).Return(errors.New("deletion failed"))

		service := NewService(Config{Store: store})

		w := performRequest(service.Handler(), "DELETE", "/sequences/1/steps/1", "")
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"error":"%s"`, ErrResourceDeletionFailed.Error()))
	})
}
