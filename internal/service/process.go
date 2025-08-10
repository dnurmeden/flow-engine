package service

import (
	"context"
	"errors"
	"github.com/dnurmeden/flow-engine/internal/models"
	"github.com/dnurmeden/flow-engine/internal/repo"
	"log"
)

type ProcessService struct {
	defRepo  *repo.DefinitionRepo
	instRepo *repo.InstanceRepo
	taskRepo *repo.TaskRepo
}

func NewProcessService(defRepo *repo.DefinitionRepo, instRepo *repo.InstanceRepo, taskRepo *repo.TaskRepo) *ProcessService {
	return &ProcessService{defRepo: defRepo, instRepo: instRepo, taskRepo: taskRepo}
}

func (s *ProcessService) StartProcess(ctx context.Context, req models.StartProcessRequest) (*models.StartProcessResponse, error) {
	def, err := s.defRepo.GetByKeyAndVersion(ctx, req.Key, req.Version)
	if err != nil {
		log.Println("1: ", err)
		return nil, err
	}
	if def == nil {
		log.Println("2: ", err)
		return nil, errors.New("process definition not found")
	}

	id, err := s.instRepo.Create(ctx, def.ID, req.Ctx)
	if err != nil {
		log.Println("3: ", err)
		return nil, err
	}

	err = s.instRepo.LogEvent(ctx, id, "process.started", req.Ctx)
	if err != nil {
		log.Println("4: ", err)
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

func (s *ProcessService) ClaimTask(ctx context.Context, taskID int64, user string) error {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.New("task not found")
	}

	if err := s.taskRepo.Claim(ctx, taskID, user); err != nil {
		return err
	}
	_ = s.taskRepo.LogEvent(ctx, t.InstanceID, "task.claimed", map[string]any{"task_id": taskID, "user": user})
	return nil
}

func (s *ProcessService) CompleteTask(ctx context.Context, taskID int64, user string, output map[string]any) error {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.New("task not found")
	}

	if err := s.taskRepo.Complete(ctx, taskID, user, output); err != nil {
		return err
	}
	_ = s.taskRepo.LogEvent(ctx, t.InstanceID, "task.completed", map[string]any{"task_id": taskID, "user": user})

	// На будущее — продвижение процесса
	_ = s.taskRepo.AdvanceAfterComplete(ctx, t.InstanceID)
	return nil
}
