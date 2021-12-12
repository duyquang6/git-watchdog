package service

import (
	"context"
	"fmt"
	"github.com/duyquang6/git-watchdog/internal/database"
	"github.com/duyquang6/git-watchdog/internal/model"
	"github.com/duyquang6/git-watchdog/internal/rabbitmq"
	"github.com/duyquang6/git-watchdog/internal/repository"
	"github.com/duyquang6/git-watchdog/pkg/customtypes"
	"github.com/duyquang6/git-watchdog/pkg/dto"
	"github.com/duyquang6/git-watchdog/pkg/exception"
	"github.com/duyquang6/git-watchdog/pkg/logging"
	"github.com/duyquang6/git-watchdog/pkg/null"
	scanpb "github.com/duyquang6/git-watchdog/proto/v1"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const (
	repositoryServiceLoggingFmt = "RepositoryService.%s"
)

// RepositoryService provide purchase service functionality
type RepositoryService interface {
	GetOne(ctx context.Context, id uint) (*dto.GetOneRepositoryResponse, error)
	Create(ctx context.Context, req dto.CreateRepositoryRequest) (*dto.CreateRepositoryResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateRepositoryRequest) (*dto.UpdateRepositoryResponse, error)
	Delete(ctx context.Context, id uint) (*dto.DeleteRepositoryResponse, error)
	IssueScanRepo(ctx context.Context, repoID uint) (*dto.IssueScanResponse, error)
	ListScan(ctx context.Context, repoID null.Uint, page, limit uint) (*dto.ListScanResponse, error)
}

type repoSvc struct {
	dbFactory           database.DBFactory
	repositoryRepo      repository.RepoRepository
	scanRepository      repository.ScanRepository
	scanRabbitMQChannel rabbitmq.IChannel
}

// NewRepositoryService create concrete instance which implement RepositoryService
func NewRepositoryService(dbFactory database.DBFactory,
	repoRepo repository.RepoRepository,
	scanRepository repository.ScanRepository,
	scanRabbitMQChannel rabbitmq.IChannel) *repoSvc {
	return &repoSvc{
		dbFactory:           dbFactory,
		repositoryRepo:      repoRepo,
		scanRepository:      scanRepository,
		scanRabbitMQChannel: scanRabbitMQChannel,
	}
}

// GetOne for creating purchase
func (s *repoSvc) GetOne(ctx context.Context, id uint) (*dto.GetOneRepositoryResponse, error) {
	var (
		tx       = s.dbFactory.GetDB()
		function = "GetOne"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(repositoryServiceLoggingFmt, function))
	)

	data, err := s.repositoryRepo.GetByID(tx, id)
	if err != nil {
		if database.IsNotFound(err) {
			logger.Info("get repository not found id =", id)
			return nil, exception.Wrap(exception.ErrRepositoryNotFound, err, "not found repository ")
		}
		logger.Error("get repository failed")
		return nil, exception.Wrap(exception.ErrInternalServer, err, "get repository failed")
	}

	return &dto.GetOneRepositoryResponse{
		Meta: dto.Meta{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		},
		Data: dto.Repository{
			ID:   data.ID,
			Name: data.Name,
			URL:  data.URL,
		},
	}, nil
}

func (s *repoSvc) Create(ctx context.Context, req dto.CreateRepositoryRequest) (*dto.CreateRepositoryResponse, error) {
	var (
		tx       = s.dbFactory.GetDB()
		function = "Create"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(repositoryServiceLoggingFmt, function))
	)
	repoModel := model.Repository{
		Name: req.Name,
		URL:  req.URL,
	}
	err := s.repositoryRepo.Create(tx, &repoModel)
	if err != nil {
		logger.Info("create repository failed")
		return nil, exception.Wrap(exception.ErrInternalServer, err, "create repository failed")
	}

	res := &dto.CreateRepositoryResponse{
		Meta: dto.Meta{
			Code:    http.StatusCreated,
			Message: http.StatusText(http.StatusCreated),
		},
	}
	res.Data.ID = repoModel.ID
	return res, nil
}

func (s *repoSvc) Update(ctx context.Context, id uint, req dto.UpdateRepositoryRequest) (*dto.UpdateRepositoryResponse, error) {
	var (
		tx       = s.dbFactory.GetDB()
		function = "Update"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(repositoryServiceLoggingFmt, function))
	)
	repoModel := model.Repository{
		Name:      req.Name,
		URL:       req.URL,
		BaseModel: model.BaseModel{ID: id},
	}

	err := s.repositoryRepo.Update(tx, &repoModel)
	if err != nil {
		if database.IsNotFound(err) {
			logger.Info("update repository failed, record not found")
			return nil, exception.Wrap(exception.ErrRepositoryNotFound, err, "update repository failed")
		}
		logger.Info("update repository failed")
		return nil, exception.Wrap(exception.ErrInternalServer, err, "update repository failed")
	}

	return &dto.UpdateRepositoryResponse{
		Meta: dto.Meta{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		},
	}, nil
}

func (s *repoSvc) Delete(ctx context.Context, id uint) (*dto.DeleteRepositoryResponse, error) {
	var (
		tx       = s.dbFactory.GetDB()
		function = "Delete"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(repositoryServiceLoggingFmt, function))
	)

	err := s.repositoryRepo.Delete(tx, id)
	if err != nil {
		if database.IsNotFound(err) {
			logger.Info("delete repository failed, record not found")
			return nil, exception.Wrap(exception.ErrRepositoryNotFound, err, "update repository failed")
		}
		logger.Info("delete repository failed")
		return nil, exception.Wrap(exception.ErrInternalServer, err, "delete repository failed")
	}

	return &dto.DeleteRepositoryResponse{
		Meta: dto.Meta{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		},
	}, nil
}

// IssueScanRepo ...
func (s *repoSvc) IssueScanRepo(ctx context.Context, repoID uint) (*dto.IssueScanResponse, error) {
	var (
		db       = s.dbFactory.GetDB()
		tx       = s.dbFactory.GetDBWithTx()
		function = "IssueScanRepo"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(repositoryServiceLoggingFmt, function))
	)

	scanModel := model.Scan{
		RepositoryID: repoID,
		Status:       customtypes.QUEUED,
		QueuedAt:     null.NewTime(time.Now()),
	}

	err := s.scanRepository.Create(tx, &scanModel)
	defer s.dbFactory.Rollback(tx)
	if err != nil {
		logger.Error("cannot create scan task, error:", err)
		return nil, exception.Wrap(exception.ErrInternalServer, err, "create scan failed")
	}

	msg := &scanpb.ScanTask{Id: int64(scanModel.ID)}
	bytesData, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("cannot marshal proto message, error:", err)
		return nil, exception.Wrap(exception.ErrUnknownError, err, "buf bytes marshal failed")
	}
	s.dbFactory.Commit(tx)

	err = s.scanRabbitMQChannel.PublishToQueue(amqp.Publishing{
		ContentType:  "application/protobuf",
		Body:         bytesData,
		DeliveryMode: amqp.Persistent,
	})

	if err != nil {
		// delete the created scan and return error
		defer func(scanRepository repository.ScanRepository, tx *gorm.DB, id uint) {
			err := scanRepository.Delete(tx, id)
			if err != nil {
				logger.Error("failed to delete")
			}
		}(s.scanRepository, db, scanModel.ID)
		logger.Error("publish failed, err:", err)
		return nil, exception.Wrap(exception.ErrPublishToQueueFailed, err, "publish failed")
	}

	return &dto.IssueScanResponse{Meta: dto.Meta{
		Code:    http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
	}}, nil
}

// ListScan ...
func (s *repoSvc) ListScan(ctx context.Context, repoID null.Uint, page, limit uint) (*dto.ListScanResponse, error) {
	var (
		tx       = s.dbFactory.GetDB()
		function = "ListScan"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(repositoryServiceLoggingFmt, function))
		err      error
		scans    []model.Scan
		offset   = (page - 1) * limit
	)
	scans, count, err := s.scanRepository.List(tx, repoID, offset, limit)
	if err != nil {
		logger.Error("list scan failed, err:", err)
		return nil, exception.Wrap(exception.ErrInternalServer, err, "list scan failed")
	}

	resData := make([]dto.Scan, len(scans))
	for i, scan := range scans {
		resData[i] = dto.Scan{
			ID: scan.ID,
			Repository: dto.Repository{
				ID:   scan.RepositoryID,
				Name: scan.Repository.Name,
				URL:  scan.Repository.URL,
			},
			Status:   scan.Status.String(),
			Findings: scan.Findings,
			Note:     scan.Note,
		}
		if scan.QueuedAt.Valid {
			resData[i].QueuedAt = null.NewUint(uint(scan.QueuedAt.Time.Unix()))
		}
		if scan.ScanningAt.Valid {
			resData[i].ScanningAt = null.NewUint(uint(scan.ScanningAt.Time.Unix()))
		}
		if scan.FinishedAt.Valid {
			resData[i].FinishedAt = null.NewUint(uint(scan.FinishedAt.Time.Unix()))
		}
	}

	return &dto.ListScanResponse{
		Meta: dto.PaginationMeta{Meta: dto.Meta{
			Code: http.StatusOK, Message: "OK",
		}, Total: count},
		Data: resData,
	}, nil
}
