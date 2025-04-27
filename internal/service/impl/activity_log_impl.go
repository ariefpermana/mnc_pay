package impl

import (
	"context"
	"encoding/json"
	"mnc/internal/model"
	"mnc/internal/repository"
	"mnc/internal/service"
)

func NewLogServiceImp(repo repository.ActivityLogRepository) service.LogService {
	return &logServiceImpl{ActivityLogRepository: repo}
}

type logServiceImpl struct {
	repository.ActivityLogRepository
}

func (s *logServiceImpl) LogActivity(ctx context.Context, userID int, action string, req, res interface{}, status string) {
	reqJSON, _ := json.Marshal(req)
	resJSON, _ := json.Marshal(res)

	log := model.ActivityLog{
		UserID:   userID,
		Action:   action,
		Request:  string(reqJSON),
		Response: string(resJSON),
		Status:   status,
	}

	_ = s.Create(ctx, log) // Optional: handle error
}
