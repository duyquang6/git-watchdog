package integration

import (
	"github.com/duyquang6/git-watchdog/internal/model"
	"github.com/duyquang6/git-watchdog/internal/repository"
	"github.com/duyquang6/git-watchdog/pkg/customtypes"
	"github.com/duyquang6/git-watchdog/pkg/null"
	"time"
)

func (p *MySQLRepositoryTestSuite) TestMySqlScanRepository_Create() {
	db := p.env.Database().GetDB()
	r := repository.NewRepoRepository()
	scanRepo := repository.NewScanRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().NoError(err)

		scanModel := model.Scan{
			RepositoryID: data.ID,
			Repository:   data,
			Status:       customtypes.QUEUED,
			QueuedAt:     null.NewTime(time.Now()),
		}

		err = scanRepo.Create(tx, &scanModel)
		p.Assert().NoError(err)

		res, err := scanRepo.GetByID(tx, scanModel.ID)
		p.Assert().NoError(err)
		p.Assert().Equal(res.Repository.Name, scanModel.Repository.Name)
	})
}

func (p *MySQLRepositoryTestSuite) TestMySqlScanRepository_Update() {
	db := p.env.Database().GetDB()
	r := repository.NewRepoRepository()
	scanRepo := repository.NewScanRepository()

	p.Run("Update in-progress", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().NoError(err)
		scanModel := model.Scan{
			RepositoryID: data.ID,
			Repository:   data,
			Status:       customtypes.QUEUED,
			QueuedAt:     null.NewTime(time.Now()),
		}
		err = scanRepo.Create(tx, &scanModel)
		p.Assert().NoError(err)

		scanModel.Status = customtypes.IN_PROGRESS
		scanModel.ScanningAt = null.NewTime(time.Now())
		err = scanRepo.Update(tx, &scanModel)
		p.Assert().NoError(err)
		actualRes, err := scanRepo.GetByID(tx, scanModel.ID)
		p.Assert().NoError(err)
		p.Assert().Equal(scanModel.Status, actualRes.Status)
	})

	p.Run("Update fail", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		scanModel := model.Scan{
			BaseModel:    model.BaseModel{ID: 1},
			RepositoryID: data.ID,
			Repository:   data,
			Status:       customtypes.QUEUED,
			QueuedAt:     null.NewTime(time.Now()),
		}
		err := scanRepo.Update(tx, &scanModel)
		p.Assert().Error(err)
	})
}

func (p *MySQLRepositoryTestSuite) TestMySqlScanRepository_Delete() {
	db := p.env.Database().GetDB()
	r := repository.NewRepoRepository()
	scanRepo := repository.NewScanRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().NoError(err)

		scanModel := model.Scan{
			BaseModel:    model.BaseModel{ID: 1},
			RepositoryID: data.ID,
			Repository:   data,
			Status:       customtypes.QUEUED,
			QueuedAt:     null.NewTime(time.Now()),
		}

		err = scanRepo.Create(tx, &scanModel)
		p.Assert().NoError(err)

		err = scanRepo.Delete(tx, scanModel.ID)
		p.Assert().NoError(err)

		_, err = scanRepo.GetByID(tx, scanModel.ID)
		p.Assert().Error(err)
	})
}

func (p *MySQLRepositoryTestSuite) TestMySqlScanRepository_List() {
	db := p.env.Database().GetDB()
	r := repository.NewRepoRepository()
	scanRepo := repository.NewScanRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().NoError(err)

		scanModel := model.Scan{
			RepositoryID: data.ID,
			Repository:   data,
			Status:       customtypes.QUEUED,
			QueuedAt:     null.NewTime(time.Now()),
		}

		err = scanRepo.Create(tx, &scanModel)
		p.Assert().NoError(err)

		scanModel.ID = 0
		err = scanRepo.Create(tx, &scanModel)
		p.Assert().NoError(err)

		res, err := scanRepo.List(tx, null.NewUint(data.ID), 0, 10)
		p.Assert().NoError(err)
		p.Assert().Equal(2, len(res))
	})
}
