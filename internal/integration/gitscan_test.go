package integration

import (
	"context"
	"github.com/duyquang6/git-watchdog/internal/configuration"
	"github.com/duyquang6/git-watchdog/internal/core"
	"github.com/duyquang6/git-watchdog/pkg/logging"
	"github.com/duyquang6/git-watchdog/pkg/ulid"
	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type GitScanTestSuite struct {
	config  *configuration.Config
	gitScan core.GitScan
	suite.Suite
}

func (p *GitScanTestSuite) SetupSuite() {
	ctx := context.Background()
	var config configuration.Config
	if err := envconfig.ProcessWith(ctx, &config, envconfig.OsLookuper()); err != nil {
		panic(err)
	}
	p.config = &config
}

func TestGitScanTestSuite(t *testing.T) {
	suite.Run(t, &GitScanTestSuite{})
}

func (p *GitScanTestSuite) SetupTest() {
	// create root temp dir
	if err := os.MkdirAll(p.config.TempRootDir, os.ModePerm); err != nil {
		panic(err)
	}
	logger := logging.DefaultLogger()
	gitScan := core.NewGitScan(logger, p.config, p.config.RuleFilePath)
	p.gitScan = gitScan
}

func (p *GitScanTestSuite) TearDownTest() {
	// drop temp dir
	os.RemoveAll(p.config.TempRootDir)
}

func (p *GitScanTestSuite) TestGitScan_Scan() {
	randomID := ulid.GetUniqueID()

	p.Run("Found zero issue", func() {
		findings, err := p.gitScan.Scan("duyquang6", "https://github.com/duyquang6/duyquang6")
		p.Assert().NoError(err)
		p.Assert().Equal(0, len(findings))
	})

	p.Run("Found one issue", func() {
		findings, err := p.gitScan.Scan("vulnerability-code", "https://gitlab.com/nguyenduyquang06/vulnerability-code.git")
		p.Assert().NoError(err)
		p.Assert().Equal(1, len(findings))
	})

	p.Run("Invalid repository", func() {
		findings, err := p.gitScan.Scan("vulnerability-code", "https://gitlab.com/nguyenduyquang06/"+randomID)
		p.Assert().Error(err)
		p.Assert().Equal(0, len(findings))
	})

	p.Run("Repo no commit yet", func() {
		findings, err := p.gitScan.Scan("vulnerability-code", "https://github.com/duyquang6/test-repo-nocommit.git")
		p.Assert().Error(err)
		p.Assert().Equal(0, len(findings))
	})
}
