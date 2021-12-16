package consumer

import (
	"context"
	"errors"
	"github.com/duyquang6/git-watchdog/internal/core"
	_mockCore "github.com/duyquang6/git-watchdog/internal/core/mocks"
	"github.com/duyquang6/git-watchdog/internal/database"
	_mockDB "github.com/duyquang6/git-watchdog/internal/database/mocks"
	"github.com/duyquang6/git-watchdog/internal/model"
	_mockRabbitMQ "github.com/duyquang6/git-watchdog/internal/rabbitmq/mocks"
	"github.com/duyquang6/git-watchdog/pkg/null"
	scanpb "github.com/duyquang6/git-watchdog/proto/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	"github.com/duyquang6/git-watchdog/internal/rabbitmq"
	"github.com/duyquang6/git-watchdog/internal/repository"
	_mockRepo "github.com/duyquang6/git-watchdog/internal/repository/mocks"
	"github.com/duyquang6/git-watchdog/pkg/logging"
	"gorm.io/gorm"
	"testing"
)

func Test_scanConsumer_processingMessage(t *testing.T) {
	t.Parallel()

	type fields struct {
		BaseConsumer
		gitScan        core.GitScan
		dbFactory      database.DBFactory
		scanRepository repository.ScanRepository
	}

	type args struct {
		ctx     context.Context
		message rabbitmq.IDelivery
	}

	scanRepoMock := &_mockRepo.ScanRepository{}
	scanRepoGetByIDFailedMock := &_mockRepo.ScanRepository{}
	scanRepoUpdateFailedMock := &_mockRepo.ScanRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	gitScanMock := &_mockCore.GitScan{}
	gitScanMockFailed := &_mockCore.GitScan{}
	deliveryMock := &_mockRabbitMQ.IDelivery{}

	repoModel := model.Repository{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "duyquang6",
		URL:       "https://github.com/duyquang6/duyquang6",
	}

	mockModelSuccess := &model.Scan{
		BaseModel:    model.BaseModel{ID: 1},
		RepositoryID: 0,
		Repository:   repoModel,
		Status:       0,
		QueuedAt:     null.Time{},
		ScanningAt:   null.Time{},
		FinishedAt:   null.Time{},
		Findings:     nil,
		Note:         null.String{},
	}

	mockFindings := []core.Finding{
		{
			Type:   "ast",
			RuleID: "ast303",
			Location: core.Location{
				Path:     "pkg/linguyen.go",
				Position: core.Position{Begin: core.Begin{Line: 1}},
			},
		},
	}

	bytesData, _ := proto.Marshal(&scanpb.ScanTask{Id: int64(1)})
	deliveryMock.On("Body").Return(bytesData)
	deliveryMock.On("DeliveryTag").Return(uint64(1))
	deliveryMock.On("ContentType").Return("application/protobuf")
	deliveryMock.On("ConsumerTag").Return("consumer-tag")

	scanRepoMock.On("GetByID", mock.Anything, uint(1)).Return(mockModelSuccess, nil)
	scanRepoMock.On("Update", mock.Anything, mock.Anything).Return(nil)

	scanRepoGetByIDFailedMock.On("GetByID", mock.Anything, uint(1)).
		Return(nil, errors.New("unexpected"))
	scanRepoGetByIDFailedMock.On("Update", mock.Anything, mock.Anything).Return(nil)

	scanRepoUpdateFailedMock.On("GetByID", mock.Anything, uint(1)).Return(mockModelSuccess, nil)
	scanRepoUpdateFailedMock.On("Update", mock.Anything, mock.Anything).Return(errors.New("unexpected"))

	gitScanMock.On("Scan", repoModel.Name, repoModel.URL).Return(mockFindings, nil)

	gitScanMockFailed.On("Scan", repoModel.Name, repoModel.URL).
		Return(nil, errors.New("unexpected"))

	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TC1_Success",
			fields: fields{
				dbFactory:      dbFactoryMock,
				scanRepository: scanRepoMock,
				gitScan:        gitScanMock,
				BaseConsumer:   BaseConsumer{logger: logging.DefaultLogger()},
			},
			args: args{
				ctx:     context.TODO(),
				message: deliveryMock,
			},
			wantErr: false,
		},
		{
			name: "TC2_GetByIDFailed",
			fields: fields{
				dbFactory:      dbFactoryMock,
				scanRepository: scanRepoGetByIDFailedMock,
				gitScan:        gitScanMock,
				BaseConsumer:   BaseConsumer{logger: logging.DefaultLogger()},
			},
			args: args{
				ctx:     context.TODO(),
				message: deliveryMock,
			},
			wantErr: true,
		},
		{
			name: "TC3_UpdateFailed",
			fields: fields{
				dbFactory:      dbFactoryMock,
				scanRepository: scanRepoUpdateFailedMock,
				gitScan:        gitScanMock,
				BaseConsumer:   BaseConsumer{logger: logging.DefaultLogger()},
			},
			args: args{
				ctx:     context.TODO(),
				message: deliveryMock,
			},
			wantErr: true,
		},
		{
			name: "TC4_GitScanFailed",
			fields: fields{
				dbFactory:      dbFactoryMock,
				scanRepository: scanRepoMock,
				gitScan:        gitScanMockFailed,
				BaseConsumer:   BaseConsumer{logger: logging.DefaultLogger()},
			},
			args: args{
				ctx:     context.TODO(),
				message: deliveryMock,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScanConsumer{
				BaseConsumer:   tt.fields.BaseConsumer,
				gitScan:        tt.fields.gitScan,
				dbFactory:      tt.fields.dbFactory,
				scanRepository: tt.fields.scanRepository,
			}
			err := s.processingMessage(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanConsumer.processingMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
