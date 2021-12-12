package service

import (
	"context"
	"errors"
	"github.com/duyquang6/git-watchdog/internal/database"
	_mockDB "github.com/duyquang6/git-watchdog/internal/database/mocks"
	"github.com/duyquang6/git-watchdog/internal/model"
	"github.com/duyquang6/git-watchdog/internal/rabbitmq"
	_mockRabbitMQ "github.com/duyquang6/git-watchdog/internal/rabbitmq/mocks"
	"github.com/duyquang6/git-watchdog/internal/repository"
	_mockRepo "github.com/duyquang6/git-watchdog/internal/repository/mocks"
	"github.com/duyquang6/git-watchdog/pkg/customtypes"
	"github.com/duyquang6/git-watchdog/pkg/dto"
	"github.com/duyquang6/git-watchdog/pkg/null"
	scanpb "github.com/duyquang6/git-watchdog/proto/v1"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewRepositoryService(t *testing.T) {
	t.Parallel()
	type args struct {
		dbFactory           database.DBFactory
		repositoryRepo      repository.RepoRepository
		scanRepository      repository.ScanRepository
		scanRabbitMQChannel rabbitmq.IChannel
	}
	rabbitMQChannelMock := &_mockRabbitMQ.IChannel{}
	repoRepoMock := &_mockRepo.RepoRepository{}
	scanRepoMock := &_mockRepo.ScanRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	tests := []struct {
		name string
		args args
		want RepositoryService
	}{
		{
			name: "TC1_NewPurchaseServiceSuccess",
			args: args{
				dbFactory:           dbFactoryMock,
				repositoryRepo:      repoRepoMock,
				scanRepository:      scanRepoMock,
				scanRabbitMQChannel: rabbitMQChannelMock,
			},
			want: &repoSvc{dbFactoryMock, repoRepoMock,
				scanRepoMock, rabbitMQChannelMock},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepositoryService(tt.args.dbFactory, tt.args.repositoryRepo, tt.args.scanRepository,
				tt.args.scanRabbitMQChannel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepositoryService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repoSvc_GetOne(t *testing.T) {
	t.Parallel()

	type fields struct {
		dbFactory           database.DBFactory
		repositoryRepo      repository.RepoRepository
		scanRepository      repository.ScanRepository
		scanRabbitMQChannel rabbitmq.IChannel
	}

	type args struct {
		ctx context.Context
		id  uint
	}

	//rabbitMQChannelMock := &_mockRabbitMQ.IChannel{}
	repoRepoMock := &_mockRepo.RepoRepository{}
	//scanRepoMock := &_mockRepo.ScanRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	repoModel := model.Repository{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "duyquang6",
		URL:       "https://github.com/duyquang6/duyquang6",
	}

	dbFactoryMock.On("GetDB").Return(&gorm.DB{})
	repoRepoMock.On("GetByID", mock.Anything, uint(1)).Return(&repoModel, nil)
	repoRepoMock.On("GetByID", mock.Anything, uint(2)).Return(nil, errors.New("unexpected"))
	repoRepoMock.On("GetByID", mock.Anything, uint(3)).Return(nil, gorm.ErrRecordNotFound)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.GetOneRepositoryResponse
		wantErr bool
	}{
		{
			name: "TC1_Success",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			wantErr: false,
			want: &dto.GetOneRepositoryResponse{
				Meta: dto.Meta{
					Code:    http.StatusOK,
					Message: http.StatusText(http.StatusOK),
				},
				Data: dto.Repository{
					ID:   1,
					Name: "duyquang6",
					URL:  "https://github.com/duyquang6/duyquang6",
				},
			},
		},
		{
			name: "TC2_UnexpectedError",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				id:  2,
			},
			wantErr: true,
		},
		{
			name: "TC3_RecordNotFound",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				id:  3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &repoSvc{
				dbFactory:           tt.fields.dbFactory,
				repositoryRepo:      tt.fields.repositoryRepo,
				scanRepository:      tt.fields.scanRepository,
				scanRabbitMQChannel: tt.fields.scanRabbitMQChannel,
			}
			got, err := s.GetOne(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("repoSvc.GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repoSvc.GetOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repoSvc_Create(t *testing.T) {
	t.Parallel()

	type fields struct {
		dbFactory           database.DBFactory
		repositoryRepo      repository.RepoRepository
		scanRepository      repository.ScanRepository
		scanRabbitMQChannel rabbitmq.IChannel
	}

	type args struct {
		ctx context.Context
		req dto.CreateRepositoryRequest
	}

	//rabbitMQChannelMock := &_mockRabbitMQ.IChannel{}
	repoRepoMock := &_mockRepo.RepoRepository{}
	//scanRepoMock := &_mockRepo.ScanRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	createRepoModelSuccess := model.Repository{
		Name: "duyquang6",
		URL:  "https://github.com/duyquang6/duyquang6",
	}
	createRepoModelFailed := model.Repository{
		Name: "duyquang7",
		URL:  "https://github.com/duyquang6/duyquang6",
	}
	expectedResponse := dto.CreateRepositoryResponse{
		Meta: dto.Meta{
			Code:    http.StatusCreated,
			Message: http.StatusText(http.StatusCreated),
		},
	}
	expectedResponse.Data.ID = 0

	dbFactoryMock.On("GetDB").Return(&gorm.DB{})
	repoRepoMock.On("Create", mock.Anything, &createRepoModelSuccess).
		Return(nil)
	repoRepoMock.On("Create", mock.Anything, &createRepoModelFailed).Return(errors.New("unexpected"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.CreateRepositoryResponse
		wantErr bool
	}{
		{
			name: "TC1_Success",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.CreateRepositoryRequest{
					Name: "duyquang6",
					URL:  "https://github.com/duyquang6/duyquang6",
				},
			},
			wantErr: false,
			want:    &expectedResponse,
		},
		{
			name: "TC2_UnexpectedError",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.CreateRepositoryRequest{
					Name: "duyquang7",
					URL:  "https://github.com/duyquang6/duyquang6",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &repoSvc{
				dbFactory:           tt.fields.dbFactory,
				repositoryRepo:      tt.fields.repositoryRepo,
				scanRepository:      tt.fields.scanRepository,
				scanRabbitMQChannel: tt.fields.scanRabbitMQChannel,
			}
			got, err := s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("repoSvc.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repoSvc.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repoSvc_Update(t *testing.T) {
	t.Parallel()

	type fields struct {
		dbFactory           database.DBFactory
		repositoryRepo      repository.RepoRepository
		scanRepository      repository.ScanRepository
		scanRabbitMQChannel rabbitmq.IChannel
	}

	type args struct {
		ctx context.Context
		req dto.UpdateRepositoryRequest
		id  uint
	}

	repoRepoMock := &_mockRepo.RepoRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	updateRepoModelSuccess := model.Repository{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "duyquang6",
		URL:       "https://github.com/duyquang6/duyquang6",
	}
	updateRepoModelFailed := model.Repository{
		BaseModel: model.BaseModel{ID: 2},
		Name:      "duyquang7",
		URL:       "https://github.com/duyquang6/duyquang7",
	}
	updateRepoModelNotFound := model.Repository{
		BaseModel: model.BaseModel{ID: 3},
		Name:      "duyquang7",
		URL:       "https://github.com/duyquang6/duyquang7",
	}

	expectedResponse := dto.UpdateRepositoryResponse{
		Meta: dto.Meta{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		},
	}

	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	repoRepoMock.On("Update", mock.Anything, &updateRepoModelSuccess).
		Return(nil)
	repoRepoMock.On("Update", mock.Anything, &updateRepoModelFailed).Return(errors.New("unexpected"))
	repoRepoMock.On("Update", mock.Anything, &updateRepoModelNotFound).Return(gorm.ErrRecordNotFound)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.UpdateRepositoryResponse
		wantErr bool
	}{
		{
			name: "TC1_Success",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.UpdateRepositoryRequest{
					Name: "duyquang6",
					URL:  "https://github.com/duyquang6/duyquang6",
				},
				id: 1,
			},
			wantErr: false,
			want:    &expectedResponse,
		},
		{
			name: "TC2_UnexpectedError",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.UpdateRepositoryRequest{
					Name: "duyquang7",
					URL:  "https://github.com/duyquang6/duyquang7",
				},
				id: 2,
			},
			wantErr: true,
		},
		{
			name: "TC3_ErrorNotFound",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.UpdateRepositoryRequest{
					Name: "duyquang7",
					URL:  "https://github.com/duyquang6/duyquang7",
				},
				id: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &repoSvc{
				dbFactory:           tt.fields.dbFactory,
				repositoryRepo:      tt.fields.repositoryRepo,
				scanRepository:      tt.fields.scanRepository,
				scanRabbitMQChannel: tt.fields.scanRabbitMQChannel,
			}
			got, err := s.Update(tt.args.ctx, tt.args.id, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("repoSvc.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repoSvc.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repoSvc_Delete(t *testing.T) {
	t.Parallel()

	type fields struct {
		dbFactory           database.DBFactory
		repositoryRepo      repository.RepoRepository
		scanRepository      repository.ScanRepository
		scanRabbitMQChannel rabbitmq.IChannel
	}

	type args struct {
		ctx context.Context
		id  uint
	}

	repoRepoMock := &_mockRepo.RepoRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}

	expectedResponse := dto.DeleteRepositoryResponse{
		Meta: dto.Meta{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		},
	}

	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	repoRepoMock.On("Delete", mock.Anything, uint(1)).
		Return(nil)
	repoRepoMock.On("Delete", mock.Anything, uint(2)).Return(errors.New("unexpected"))
	repoRepoMock.On("Delete", mock.Anything, uint(3)).Return(gorm.ErrRecordNotFound)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.DeleteRepositoryResponse
		wantErr bool
	}{
		{
			name: "TC1_Success",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			wantErr: false,
			want:    &expectedResponse,
		},
		{
			name: "TC2_UnexpectedError",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				id:  2,
			},
			wantErr: true,
		},
		{
			name: "TC3_ErrorNotFound",
			fields: fields{
				dbFactory:      dbFactoryMock,
				repositoryRepo: repoRepoMock,
			},
			args: args{
				ctx: context.TODO(),
				id:  3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &repoSvc{
				dbFactory:           tt.fields.dbFactory,
				repositoryRepo:      tt.fields.repositoryRepo,
				scanRepository:      tt.fields.scanRepository,
				scanRabbitMQChannel: tt.fields.scanRabbitMQChannel,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("repoSvc.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repoSvc.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repoSvc_ListScan(t *testing.T) {
	t.Parallel()

	type fields struct {
		dbFactory           database.DBFactory
		repositoryRepo      repository.RepoRepository
		scanRepository      repository.ScanRepository
		scanRabbitMQChannel rabbitmq.IChannel
	}

	type args struct {
		ctx         context.Context
		repoID      null.Uint
		page, limit uint
	}

	scanRepoMock := &_mockRepo.ScanRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}

	timeNow := time.Now()
	expectedData := []dto.Scan{
		{
			ID:       1,
			Status:   customtypes.QUEUED.String(),
			QueuedAt: null.NewUint(uint(timeNow.Unix())),
			Repository: dto.Repository{
				ID:   1,
				URL:  "1",
				Name: "1",
			},
		},
	}
	expectedResponse := dto.ListScanResponse{
		Meta: dto.PaginationMeta{
			Meta: dto.Meta{Code: http.StatusOK,
				Message: http.StatusText(http.StatusOK)},
			Total: 2,
		},
		Data: expectedData,
	}

	listRes := []model.Scan{
		{
			BaseModel:    model.BaseModel{ID: 1},
			Status:       customtypes.QUEUED,
			QueuedAt:     null.NewTime(timeNow),
			RepositoryID: 1,
			Repository: model.Repository{
				BaseModel: model.BaseModel{ID: 1},
				Name:      "1",
				URL:       "1",
			},
		},
	}

	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	scanRepoMock.On("List", mock.Anything, null.NewUint(1), uint(0), uint(1)).
		Return(listRes, uint(2), nil)
	scanRepoMock.On("List", mock.Anything, null.NewUint(2), uint(0), uint(1)).
		Return(nil, uint(0), errors.New("unexpected"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.ListScanResponse
		wantErr bool
	}{
		{
			name: "TC1_Success",
			fields: fields{
				dbFactory:      dbFactoryMock,
				scanRepository: scanRepoMock,
			},
			args: args{
				ctx:    context.TODO(),
				repoID: null.NewUint(1),
				page:   1,
				limit:  1,
			},
			wantErr: false,
			want:    &expectedResponse,
		},
		{
			name: "TC2_UnexpectedError",
			fields: fields{
				dbFactory:      dbFactoryMock,
				scanRepository: scanRepoMock,
			},
			args: args{
				ctx:    context.TODO(),
				repoID: null.NewUint(2),
				page:   1,
				limit:  1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &repoSvc{
				dbFactory:           tt.fields.dbFactory,
				repositoryRepo:      tt.fields.repositoryRepo,
				scanRepository:      tt.fields.scanRepository,
				scanRabbitMQChannel: tt.fields.scanRabbitMQChannel,
			}
			got, err := s.ListScan(tt.args.ctx, tt.args.repoID, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("repoSvc.ListScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repoSvc.ListScan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repoSvc_IssueScanRepo(t *testing.T) {
	t.Parallel()

	type fields struct {
		dbFactory           database.DBFactory
		repositoryRepo      repository.RepoRepository
		scanRepository      repository.ScanRepository
		scanRabbitMQChannel rabbitmq.IChannel
	}

	type args struct {
		ctx    context.Context
		repoID uint
	}

	scanRepoMock := &_mockRepo.ScanRepository{}
	scanRepoMockCreateFailed := &_mockRepo.ScanRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	channelMock := &_mockRabbitMQ.IChannel{}
	channelPublishFailedMock := &_mockRabbitMQ.IChannel{}

	bytesData, _ := proto.Marshal(&scanpb.ScanTask{Id: int64(0)})
	scanRepoMock.On("Create", mock.Anything, mock.Anything).
		Return(nil)
	scanRepoMockCreateFailed.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("unexpected error"))
	scanRepoMock.On("Delete", mock.Anything, uint(0)).
		Return(nil)

	channelMock.On("PublishToQueue", amqp.Publishing{
		ContentType:  "application/protobuf",
		Body:         bytesData,
		DeliveryMode: amqp.Persistent,
	}).Return(nil)

	channelPublishFailedMock.On("PublishToQueue", amqp.Publishing{
		ContentType:  "application/protobuf",
		Body:         bytesData,
		DeliveryMode: amqp.Persistent,
	}).Return(errors.New("unexpected error"))

	dbFactoryMock.On("GetDB").Return(&gorm.DB{})
	dbFactoryMock.On("GetDBWithTx").Return(&gorm.DB{})
	dbFactoryMock.On("Rollback", mock.Anything).Return()
	dbFactoryMock.On("Commit", mock.Anything).Return()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.IssueScanResponse
		wantErr bool
	}{
		{
			name: "TC1_Success",
			fields: fields{
				dbFactory:           dbFactoryMock,
				scanRepository:      scanRepoMock,
				scanRabbitMQChannel: channelMock,
			},
			args: args{
				ctx:    context.TODO(),
				repoID: 1,
			},
			wantErr: false,
			want: &dto.IssueScanResponse{Meta: dto.Meta{
				Code:    http.StatusCreated,
				Message: http.StatusText(http.StatusCreated),
			}},
		},
		{
			name: "TC2_CreateFailed",
			fields: fields{
				dbFactory:           dbFactoryMock,
				scanRepository:      scanRepoMockCreateFailed,
				scanRabbitMQChannel: channelMock,
			},
			args: args{
				ctx:    context.TODO(),
				repoID: 2,
			},
			wantErr: true,
		},
		{
			name: "TC3_PublishChannelFailed",
			fields: fields{
				dbFactory:           dbFactoryMock,
				scanRepository:      scanRepoMock,
				scanRabbitMQChannel: channelPublishFailedMock,
			},
			args: args{
				ctx:    context.TODO(),
				repoID: 2,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &repoSvc{
				dbFactory:           tt.fields.dbFactory,
				repositoryRepo:      tt.fields.repositoryRepo,
				scanRepository:      tt.fields.scanRepository,
				scanRabbitMQChannel: tt.fields.scanRabbitMQChannel,
			}
			got, err := s.IssueScanRepo(tt.args.ctx, tt.args.repoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("repoSvc.IssueScanRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repoSvc.IssueScanRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
