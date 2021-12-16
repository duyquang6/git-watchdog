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
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"time"
)

var _ Consumer = (*ScanConsumer)(nil)

type ScanConsumer struct {
	BaseConsumer
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
	scanConsumer := &ScanConsumer{
		gitScan:        gitScan,
		dbFactory:      dbFactory,
		scanRepository: scanRepository,
	}
	scanConsumer.BaseConsumer = BaseConsumer{
		conn:      conn,
		channel:   channel,
		tag:       tag,
		done:      done,
		queueName: queueName,
		logger:    logger,
		consumer:  scanConsumer,
	}
	return scanConsumer
}

func (c *ScanConsumer) handle(ctx context.Context, deliveries <-chan amqp.Delivery) {

	for delivery := range deliveries {
		d := rabbitmq.NewDelivery(delivery)

		err := c.processingMessage(ctx, d)
		if err != nil {
			c.logger.Errorf("failed to handle message, "+
				"move message tag %d to dead letter queue", d.DeliveryTag())
			if d.Nack(false, false) != nil {
				c.logger.Error("cannot nack message with delivery_tag:", d.DeliveryTag())
			}
		} else {
			if d.Ack(false) != nil {
				c.logger.Error("cannot ack message with delivery_tag:", d.DeliveryTag())
			}
		}
	}

	c.logger.Info("handle: deliveries channel closed")
	c.done <- nil
}

func (c *ScanConsumer) processingMessage(ctx context.Context, message rabbitmq.IDelivery) error {
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
			err := c.updateFailed(scan, err, tx)
			if err != nil {
				c.logger.Error("failed to update failed status, error:", err)
			}
		}
	}()

	err = c.updateInProgress(scan, tx)
	if err != nil {
		c.logger.Error("failed to update in-progress status, error:", err)
		return err
	}

	findings, err := c.gitScan.Scan(scan.Repository.Name, scan.Repository.URL)

	if err != nil {
		c.logger.Error("scan failed, error:", err)
		return err
	}

	err = c.updateSuccess(scan, tx, findings)
	if err != nil {
		c.logger.Error("update success status failed, error:", err)
		return err
	}

	c.logger.Infof("success processing message, deli_tag %d, task_id %d", message.DeliveryTag(), msgProto.Id)
	return nil
}

func (c *ScanConsumer) updateInProgress(scan *model.Scan, tx *gorm.DB) error {
	scan.Status = customtypes.IN_PROGRESS
	scan.UpdatedAt = time.Now()
	scan.ScanningAt = null.NewTime(time.Now())
	return c.scanRepository.Update(tx, scan)
}

func (c *ScanConsumer) updateFailed(scan *model.Scan, err error, tx *gorm.DB) error {
	scan.Status = customtypes.FAILURE
	scan.UpdatedAt = time.Now()
	scan.FinishedAt = null.NewTime(time.Now())
	scan.Note = null.NewString("got error: " + err.Error())
	err = c.scanRepository.Update(tx, scan)
	return err
}

func (c *ScanConsumer) updateSuccess(scan *model.Scan, tx *gorm.DB, findings []core.Finding) error {
	bytes, err := json.Marshal(findings)
	if err != nil {
		return err
	}
	scan.Status = customtypes.SUCCESS
	scan.UpdatedAt = time.Now()
	scan.FinishedAt = null.NewTime(time.Now())
	scan.Findings = bytes
	err = c.scanRepository.Update(tx, scan)
	return err
}
