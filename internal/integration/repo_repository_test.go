package integration

import (
	"github.com/duyquang6/git-watchdog/internal/model"
	"github.com/duyquang6/git-watchdog/internal/repository"
)

func (p *MySQLRepositoryTestSuite) TestMySqlRepoRepository_Create() {
	db := p.env.Database().GetDB()
	r := repository.NewRepoRepository()

	p.Run("Failed because of invalid repository name", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().Error(err)
	})

	p.Run("Failed because of invalid repository url", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().Error(err)
	})

	p.Run("OK", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().NoError(err)

		res, err := r.GetByID(tx, data.ID)
		p.Assert().NoError(err)
		p.Assert().Equal(data.URL, res.URL)
	})
}

func (p *MySQLRepositoryTestSuite) TestMySqlRepoRepository_Update() {
	db := p.env.Database().GetDB()
	r := repository.NewRepoRepository()

	p.Run("OK", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().NoError(err)

		data.URL = "https://github.com/duyquang6/duyquang6_copy"
		err = r.Update(tx, &data)
		p.Assert().NoError(err)

		data.URL = "https://github.com/duyquang6/duyquang6_copy"
		res, err := r.GetByID(tx, data.ID)
		p.Assert().NoError(err)
		p.Assert().Equal("https://github.com/duyquang6/duyquang6_copy", res.URL)
	})
	p.Run("Failed invalid url", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().NoError(err)

		data.URL = "uiuh"
		err = r.Update(tx, &data)
		p.Assert().Error(err)
	})
}

func (p *MySQLRepositoryTestSuite) TestMySqlRepoRepository_Delete() {
	db := p.env.Database().GetDB()
	r := repository.NewRepoRepository()

	p.Run("OK", func() {
		tx := db.Begin()
		defer tx.Rollback()
		data := model.Repository{
			Name: "duyquang6",
			URL:  "https://github.com/duyquang6/duyquang6",
		}
		err := r.Create(tx, &data)
		p.Assert().NoError(err)
		err = r.Delete(tx, data.ID)
		p.Assert().NoError(err)

		_, err = r.GetByID(tx, data.ID)
		p.Assert().Error(err)
	})
}
