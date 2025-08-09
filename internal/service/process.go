package service

import (
	"context"
	"errors"
	"github.com/dnurmeden/flow-engine/internal/models"
	"github.com/dnurmeden/flow-engine/internal/repo"
)

type ProcessService struct {
	defRepo  *repo.DefinitionRepo
	instRepo *repo.InstanceRepo
}

func NewProcessService(defRepo *repo.DefinitionRepo, instRepo *repo.InstanceRepo) *ProcessService {
	return &ProcessService{defRepo: defRepo, instRepo: instRepo}
}

func (s *ProcessService) StartProcess(ctx context.Context, req models.StartProcessRequest) (*models.StartProcessResponse, error) {
	def, err := s.defRepo.GetByKeyAndVersion(ctx, req.Key, req.Version)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return nil, errors.New("process definition not found")
	}

	id, err := s.instRepo.Create(ctx, def.ID, req.Ctx)
	if err != nil {
		return nil, err
	}

	err = s.instRepo.LogEvent(ctx, id, "process.started", req.Ctx)
	if err != nil {
		return nil, err
	}

	return &models.StartProcessResponse{
		InstanceID: id,
		Status:     "running",
	}, nil
}

func (s *ProcessService) GetInstance(ctx context.Context, id int64) (*models.GetInstanceResponse, error) {
	inst, err := s.instRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inst == nil {
		return nil, nil
	}
	tasks, err := s.instRepo.ListOpenTasks(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.GetInstanceResponse{Instance: *inst, Tasks: tasks}, nil
}
