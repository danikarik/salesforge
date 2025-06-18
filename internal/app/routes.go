package app

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/danikarik/salesforge/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var (
	ErrResourceNotFound       = errors.New("resource not found")
	ErrResourceCreationFailed = errors.New("resource creation failed")
	ErrResourceFetchingFailed = errors.New("resource fetching failed")
	ErrResourceUpdateFailed   = errors.New("resource update failed")
	ErrResourceDeletionFailed = errors.New("resource deletion failed")
)

type CreateStepRequest struct {
	Subject string `json:"subject" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CreateSequenceRequest struct {
	Name                 string              `json:"name" binding:"required"`
	OpenTrackingEnabled  bool                `json:"openTrackingEnabled"`
	ClickTrackingEnabled bool                `json:"clickTrackingEnabled"`
	Steps                []CreateStepRequest `json:"steps" binding:"required,dive"`
}

func (s *Service) createSequence(c *gin.Context) {
	var data CreateSequenceRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sequence := &model.Sequence{
		Name:                 data.Name,
		OpenTrackingEnabled:  data.OpenTrackingEnabled,
		ClickTrackingEnabled: data.ClickTrackingEnabled,
		Steps:                make([]*model.Step, len(data.Steps)),
	}
	for i, step := range data.Steps {
		sequence.Steps[i] = &model.Step{
			Subject: step.Subject,
			Content: step.Content,
		}
	}

	if err := s.store.CreateSequence(c.Request.Context(), sequence); err != nil {
		log.Printf("Failed to create sequence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrResourceCreationFailed.Error()})
		return
	}

	c.JSON(http.StatusCreated, sequence)
}

func (s *Service) fetchSequence(c *gin.Context) {
	id, err := fetchResourceID(c, "id", "Invalid sequence ID")
	if err != nil {
		return
	}

	sequence, err := s.store.FetchSequence(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrResourceNotFound.Error()})
			return
		}
		log.Printf("Failed to fetch sequence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrResourceFetchingFailed.Error()})
		return
	}

	c.JSON(http.StatusOK, sequence)
}

type UpdateSequenceRequest struct {
	OpenTrackingEnabled  bool `json:"openTrackingEnabled"`
	ClickTrackingEnabled bool `json:"clickTrackingEnabled"`
}

type UpdateSequenceResponse struct {
	ID                   uint64    `json:"id"`
	OpenTrackingEnabled  bool      `json:"openTrackingEnabled"`
	ClickTrackingEnabled bool      `json:"clickTrackingEnabled"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

func (s *Service) updateSequence(c *gin.Context) {
	id, err := fetchResourceID(c, "id", "Invalid sequence ID")
	if err != nil {
		return
	}

	var data UpdateSequenceRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sequence := &model.Sequence{
		OpenTrackingEnabled:  data.OpenTrackingEnabled,
		ClickTrackingEnabled: data.ClickTrackingEnabled,
	}
	if err := s.store.UpdateSequence(c.Request.Context(), id, sequence); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrResourceNotFound.Error()})
			return
		}
		log.Printf("Failed to update sequence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrResourceUpdateFailed.Error()})
		return
	}

	c.JSON(http.StatusOK, UpdateSequenceResponse{
		ID:                   sequence.ID,
		OpenTrackingEnabled:  sequence.OpenTrackingEnabled,
		ClickTrackingEnabled: sequence.ClickTrackingEnabled,
		UpdatedAt:            sequence.UpdatedAt,
	})
}

type UpdateStepRequest struct {
	Subject string `json:"subject" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (s *Service) updateStep(c *gin.Context) {
	sequenceID, err := fetchResourceID(c, "id", "Invalid sequence ID")
	if err != nil {
		return
	}
	id, err := fetchResourceID(c, "step_id", "Invalid step ID")
	if err != nil {
		return
	}

	var data UpdateStepRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	step := &model.Step{
		SequenceID: sequenceID,
		Subject:    data.Subject,
		Content:    data.Content,
	}
	if err := s.store.UpdateStep(c.Request.Context(), id, step); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrResourceNotFound.Error()})
			return
		}
		log.Printf("Failed to update step: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrResourceUpdateFailed.Error()})
		return
	}

	c.JSON(http.StatusOK, step)
}

func (s *Service) deleteStep(c *gin.Context) {
	sequenceID, err := fetchResourceID(c, "id", "Invalid sequence ID")
	if err != nil {
		return
	}
	id, err := fetchResourceID(c, "step_id", "Invalid step ID")
	if err != nil {
		return
	}

	step := &model.Step{SequenceID: sequenceID}
	if err := s.store.DeleteStep(c.Request.Context(), id, step); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrResourceNotFound.Error()})
			return
		}
		log.Printf("Failed to delete step: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrResourceDeletionFailed.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func fetchResourceID(c *gin.Context, param, errorMessage string) (uint64, error) {
	raw := c.Param(param)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return 0, err
	}
	return uint64(id), nil
}
