package consumer

import (
	"context"
	"encoding/json"
	"github.com/duyquang6/git-watchdog/internal/core"
	"github.com/duyquang6/git-watchdog/internal/database"
	"github.com/duyquang6/git-watchdog/internal/model"
	"github.com/duyquang6/git-watchdog/internal/rabbitmq"
	"github.com/duyquang6/git-watchdog/internal/repository"
	"github.com/duyquang6/git-watchdog/pkg/customtypes"
	"github.com/duyquang6/git-watchdog/pkg/null"
	scanpb "github.com/duyquang6/git-watchdog/proto/v1"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var _ Consumer = (*ScanConsumer)(nil)

type ScanConsumer struct {
	BaseConsumer
	logger         *zap.SugaredLogger
	gitScan        core.GitScan
	dbFactory      database.DBFactory
	scanRepository repository.ScanRepository
}

func NewScanConsumer(
	logger *zap.SugaredLogger,
	dbFactory database.DBFactory,
	scanRepository repository.ScanRepository,
	gitScan core.GitScan,
	conn *amqp.Connection,
	channel rabbitmq.IChannel,
	tag string,
	done chan error,
	queueName string) *ScanConsumer {
	return &ScanConsumer{
		BaseConsumer: BaseConsumer{
			conn:      conn,
			channel:   channel,
			tag:       tag,
			done:      done,
			queueName: queueName,
		},
		logger:         logger,
		gitScan:        gitScan,
		dbFactory:      dbFactory,
		scanRepository: scanRepository,
	}
}

func (c *ScanConsumer) handle(ctx context.Context, deliveries <-chan amqp.Delivery) {

	for delivery := range deliveries {
		d := rabbitmq.NewDelivery(delivery)

		err := c.processingMessage(ctx, d)
		if err != nil {
			c.logger.Errorf("failed to handle message, "+
				"move message tag %d to dead letter queue", d.DeliveryTag())
			d.Nack(false, false)
		} else {
			d.Ack(false)
		}
	}

	c.logger.Info("handle: deliveries channel closed")
	c.done <- nil
}

func (c *ScanConsumer) processingMessage(ctx context.Context, message rabbitmq.Delivery) error {
	var (
		tx       = c.dbFactory.GetDB()
		msgProto scanpb.ScanTask
	)

	err := proto.Unmarshal(message.Body(), &msgProto)
	if err != nil {
		return err
	}

	c.logger.Infof(
		"got %dB delivery: [%v] content-type: %s, consumer-tag: %s, task_id: %d",
		len(message.Body()),
		message.DeliveryTag(),
		message.ContentType(),
		message.ConsumerTag(),
		msgProto.Id,
	)

	scan, err := c.scanRepository.GetByID(tx, uint(msgProto.Id))
	if err != nil {
		c.logger.Error("failed to get scan information, error:", err)
		return err
	}

	defer func() {
		if err != nil {
			err := c.updateFailed(&scan, err, tx)
			if err != nil {
				c.logger.Error("failed to update failed status, error:", err)
			}
		}
	}()

	err = c.updateInProgress(&scan, err, tx)
	if err != nil {
		c.logger.Error("failed to update in-progress status, error:", err)
		return err
	}

	findings, err := c.gitScan.Scan(scan.Repository.Name, scan.Repository.URL)

	if err != nil {
		c.logger.Error("scan failed, error:", err)
		return err
	}

	err = c.updateSuccess(&scan, err, tx, findings)
	if err != nil {
		c.logger.Error("update success status failed, error:", err)
		return err
	}

	c.logger.Infof("success processing message, deli_tag %d, task_id %d", message.DeliveryTag(), msgProto.Id)
	return nil
}

func (c *ScanConsumer) updateInProgress(scan *model.Scan, err error, tx *gorm.DB) error {
	scan.Status = customtypes.IN_PROGRESS
	scan.UpdatedAt = time.Now()
	scan.ScanningAt = null.NewTime(time.Now())
	err = c.scanRepository.Update(tx, scan)
	return err
}

func (c *ScanConsumer) updateFailed(scan *model.Scan, err error, tx *gorm.DB) error {
	scan.Status = customtypes.FAILURE
	scan.UpdatedAt = time.Now()
	scan.FinishedAt = null.NewTime(time.Now())
	scan.Note = null.NewString("got error: " + err.Error())
	err = c.scanRepository.Update(tx, scan)
	return err
}

func (c *ScanConsumer) updateSuccess(scan *model.Scan, err error, tx *gorm.DB, findings []core.Finding) error {
	scan.Status = customtypes.SUCCESS
	scan.UpdatedAt = time.Now()
	scan.FinishedAt = null.NewTime(time.Now())
	bytes, _ := json.Marshal(findings)
	scan.Findings = bytes
	err = c.scanRepository.Update(tx, scan)
	return err
}
