package core

import (
	"github.com/duyquang6/git-watchdog/internal/configuration"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var _ GitScan = (*gitScan)(nil)

type gitScan struct {
	logger    *zap.SugaredLogger
	appConfig *configuration.Config
	rules     []Rule
}

func NewGitScan(logger *zap.SugaredLogger,
	appConfig *configuration.Config, ruleFilePath string) *gitScan {
	return &gitScan{logger: logger, appConfig: appConfig, rules: NewRulesFromFile(ruleFilePath)}
}

type GitScan interface {
	Scan(repoName, repoURL string) ([]Finding, error)
}

func (g *gitScan) Scan(repoName, repoURL string) ([]Finding, error) {
	var (
		findings []Finding
	)
	g.logger.Infof("Scanning repo %s ....", repoName)
	// setup temp folder
	g.logger.Info("Preparing temporary dir to clone....")
	folder, err := ioutil.TempDir(g.appConfig.TempRootDir, repoName+"-scan")
	absPath := g.appConfig.TempRootDir + folder
	defer os.RemoveAll(absPath)
	if err != nil {
		g.logger.Error("scan error:", err)
		return nil, err
	}

	g.logger.Info("Clone source...")
	gitRepo, err := git.PlainClone(absPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	// Checkout latest commit source codes
	g.logger.Info("Scanning...")
	commits, err := gitRepo.CommitObjects()
	if err != nil {
		g.logger.Error("scan error:", err)
		return nil, err
	}
	defer commits.Close()
	latestCommit, err := commits.Next()
	if err != nil {
		g.logger.Error("get latest commit error:", err)
		return nil, err
	}
	fIter, err := latestCommit.Files()
	if err != nil {
		g.logger.Error("get file iter error:", err)
		return nil, err
	}
	defer fIter.Close()

	// loop through all files
	err = fIter.ForEach(func(file *object.File) error {
		path := file.Name
		lines, _ := file.Lines()
		for i, line := range lines {
			for _, rule := range g.rules {
				if rule.Regexp.MatchString(line) {
					findings = append(findings, Finding{
						Type:   rule.Type,
						RuleID: rule.ID,
						Location: Location{
							Path:     path,
							Position: Position{Begin: Begin{Line: uint(i + 1)}},
						},
						Metadata: FindingMetadata{
							Description: rule.Description,
							Severity:    rule.Severity,
						},
					})
				}
			}
		}
		return nil
	})

	if err != nil {
		g.logger.Error("scan error:", err)
		return nil, err
	}

	g.logger.Infof("Done scan repo %s ....", repoName)
	return findings, nil
}
