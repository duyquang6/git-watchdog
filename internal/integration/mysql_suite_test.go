//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/duyquang6/git-watchdog/internal/database"
	"github.com/duyquang6/git-watchdog/internal/serverenv"
	"github.com/duyquang6/git-watchdog/internal/setup"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MySQLRepositoryTestSuite struct {
	env    *serverenv.ServerEnv
	config database.Config
	suite.Suite
}

func (p *MySQLRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	var config database.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		panic(err)
	}
	p.env = env
}

func TestMySqlRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &MySQLRepositoryTestSuite{})
}

func (p *MySQLRepositoryTestSuite) SetupTest() {
	ctx := context.Background()
	if err := p.env.Database().Migrate(ctx); err != nil {
		panic(err)
	}
}

func (p *MySQLRepositoryTestSuite) TearDownTest() {
	ctx := context.Background()
	if err := p.env.Database().MigrateDown(ctx); err != nil {
		panic(err)
	}
}
