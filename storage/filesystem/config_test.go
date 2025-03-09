package filesystem

import (
	"os"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/jesseduffield/go-git/v5/config"
	"github.com/jesseduffield/go-git/v5/storage/filesystem/dotgit"
	"github.com/stretchr/testify/suite"
)

type ConfigFixtureSuite struct {
	fixtures.Suite
}

type ConfigSuite struct {
	suite.Suite
	ConfigFixtureSuite

	dir  *dotgit.DotGit
	path string
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) SetupTest() {
	tmp, err := util.TempDir(osfs.Default, "", "go-git-filestystem-config")
	s.NoError(err)

	s.dir = dotgit.New(osfs.New(tmp))
	s.path = tmp
}

func (s *ConfigSuite) TestRemotes() {
	dir := dotgit.New(fixtures.Basic().ByTag(".git").One().DotGit())
	storer := &ConfigStorage{dir}

	cfg, err := storer.Config()
	s.NoError(err)

	remotes := cfg.Remotes
	s.Len(remotes, 1)
	remote := remotes["origin"]
	s.Equal("origin", remote.Name)
	s.Equal([]string{"https://github.com/git-fixtures/basic"}, remote.URLs)
	s.Equal([]config.RefSpec{config.RefSpec("+refs/heads/*:refs/remotes/origin/*")}, remote.Fetch)
}

func (s *ConfigSuite) TearDownTest() {
	defer os.RemoveAll(s.path)
}
