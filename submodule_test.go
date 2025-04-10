package git

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/jesseduffield/go-git/v5/config"
	"github.com/jesseduffield/go-git/v5/plumbing"
	"github.com/jesseduffield/go-git/v5/storage/memory"

	fixtures "github.com/go-git/go-git-fixtures/v4"
	. "gopkg.in/check.v1"
)

type SubmoduleSuite struct {
	BaseSuite
	Worktree *Worktree
}

var _ = Suite(&SubmoduleSuite{})

func (s *SubmoduleSuite) SetUpTest(c *C) {
	path := fixtures.ByTag("submodule").One().Worktree().Root()

	dir := c.MkDir()

	r, err := PlainClone(filepath.Join(dir, "worktree"), false, &CloneOptions{
		URL: path,
	})

	c.Assert(err, IsNil)

	s.Repository = r
	s.Worktree, err = r.Worktree()
	c.Assert(err, IsNil)
}

func (s *SubmoduleSuite) TestInit(c *C) {
	sm, err := s.Worktree.Submodule("basic")
	c.Assert(err, IsNil)

	c.Assert(sm.initialized, Equals, false)
	err = sm.Init()
	c.Assert(err, IsNil)

	c.Assert(sm.initialized, Equals, true)

	cfg, err := s.Repository.Config()
	c.Assert(err, IsNil)

	c.Assert(cfg.Submodules, HasLen, 1)
	c.Assert(cfg.Submodules["basic"], NotNil)

	status, err := sm.Status()
	c.Assert(err, IsNil)
	c.Assert(status.IsClean(), Equals, false)
}

func (s *SubmoduleSuite) TestUpdate(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	sm, err := s.Worktree.Submodule("basic")
	c.Assert(err, IsNil)

	err = sm.Update(&SubmoduleUpdateOptions{
		Init: true,
	})

	c.Assert(err, IsNil)

	r, err := sm.Repository()
	c.Assert(err, IsNil)

	ref, err := r.Reference(plumbing.HEAD, true)
	c.Assert(err, IsNil)
	c.Assert(ref.Hash().String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")

	status, err := sm.Status()
	c.Assert(err, IsNil)
	c.Assert(status.IsClean(), Equals, true)
}

func (s *SubmoduleSuite) TestRepositoryWithoutInit(c *C) {
	sm, err := s.Worktree.Submodule("basic")
	c.Assert(err, IsNil)

	r, err := sm.Repository()
	c.Assert(err, Equals, ErrSubmoduleNotInitialized)
	c.Assert(r, IsNil)
}

func (s *SubmoduleSuite) TestUpdateWithoutInit(c *C) {
	sm, err := s.Worktree.Submodule("basic")
	c.Assert(err, IsNil)

	err = sm.Update(&SubmoduleUpdateOptions{})
	c.Assert(err, Equals, ErrSubmoduleNotInitialized)
}

func (s *SubmoduleSuite) TestUpdateWithNotFetch(c *C) {
	sm, err := s.Worktree.Submodule("basic")
	c.Assert(err, IsNil)

	err = sm.Update(&SubmoduleUpdateOptions{
		Init:    true,
		NoFetch: true,
	})

	// Since we are not fetching, the object is not there
	c.Assert(err, Equals, plumbing.ErrObjectNotFound)
}

func (s *SubmoduleSuite) TestUpdateWithRecursion(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	sm, err := s.Worktree.Submodule("itself")
	c.Assert(err, IsNil)

	err = sm.Update(&SubmoduleUpdateOptions{
		Init:              true,
		RecurseSubmodules: 2,
	})

	c.Assert(err, IsNil)

	fs := s.Worktree.Filesystem
	_, err = fs.Stat(fs.Join("itself", "basic", "LICENSE"))
	c.Assert(err, IsNil)
}

func (s *SubmoduleSuite) TestUpdateWithInitAndUpdate(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	sm, err := s.Worktree.Submodule("basic")
	c.Assert(err, IsNil)

	err = sm.Update(&SubmoduleUpdateOptions{
		Init: true,
	})
	c.Assert(err, IsNil)

	idx, err := s.Repository.Storer.Index()
	c.Assert(err, IsNil)

	for i, e := range idx.Entries {
		if e.Name == "basic" {
			e.Hash = plumbing.NewHash("b029517f6300c2da0f4b651b8642506cd6aaf45d")
		}

		idx.Entries[i] = e
	}

	err = s.Repository.Storer.SetIndex(idx)
	c.Assert(err, IsNil)

	err = sm.Update(&SubmoduleUpdateOptions{})
	c.Assert(err, IsNil)

	r, err := sm.Repository()
	c.Assert(err, IsNil)

	ref, err := r.Reference(plumbing.HEAD, true)
	c.Assert(err, IsNil)
	c.Assert(ref.Hash().String(), Equals, "b029517f6300c2da0f4b651b8642506cd6aaf45d")

}

func (s *SubmoduleSuite) TestSubmodulesInit(c *C) {
	sm, err := s.Worktree.Submodules()
	c.Assert(err, IsNil)

	err = sm.Init()
	c.Assert(err, IsNil)

	sm, err = s.Worktree.Submodules()
	c.Assert(err, IsNil)

	for _, m := range sm {
		c.Assert(m.initialized, Equals, true)
	}
}

func (s *SubmoduleSuite) TestGitSubmodulesSymlink(c *C) {
	f, err := s.Worktree.Filesystem.Create("badfile")
	c.Assert(err, IsNil)
	defer func() { _ = f.Close() }()

	err = s.Worktree.Filesystem.Remove(gitmodulesFile)
	c.Assert(err, IsNil)

	err = s.Worktree.Filesystem.Symlink("badfile", gitmodulesFile)
	c.Assert(err, IsNil)

	_, err = s.Worktree.Submodules()
	c.Assert(err, Equals, ErrGitModulesSymlink)
}

func (s *SubmoduleSuite) TestSubmodulesStatus(c *C) {
	sm, err := s.Worktree.Submodules()
	c.Assert(err, IsNil)

	status, err := sm.Status()
	c.Assert(err, IsNil)
	c.Assert(status, HasLen, 2)
}

func (s *SubmoduleSuite) TestSubmodulesUpdateContext(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	sm, err := s.Worktree.Submodules()
	c.Assert(err, IsNil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = sm.UpdateContext(ctx, &SubmoduleUpdateOptions{Init: true})
	c.Assert(err, NotNil)
}

func (s *SubmoduleSuite) TestSubmodulesFetchDepth(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	sm, err := s.Worktree.Submodule("basic")
	c.Assert(err, IsNil)

	err = sm.Update(&SubmoduleUpdateOptions{
		Init:  true,
		Depth: 1,
	})
	c.Assert(err, IsNil)

	r, err := sm.Repository()
	c.Assert(err, IsNil)

	lr, err := r.Log(&LogOptions{})
	c.Assert(err, IsNil)

	commitCount := 0
	for _, err := lr.Next(); err == nil; _, err = lr.Next() {
		commitCount++
	}
	c.Assert(err, IsNil)

	c.Assert(commitCount, Equals, 1)
}

func (s *SubmoduleSuite) TestSubmoduleParseScp(c *C) {
	repo := &Repository{
		Storer: memory.NewStorage(),
		wt:     memfs.New(),
	}
	worktree := &Worktree{
		Filesystem: memfs.New(),
		r:          repo,
	}
	submodule := &Submodule{
		initialized: true,
		c:           nil,
		w:           worktree,
	}

	submodule.c = &config.Submodule{
		URL: "git@github.com:username/submodule_repo",
	}

	_, err := submodule.Repository()
	c.Assert(err, IsNil)
}
