package transactional

import (
	"testing"

	"github.com/jesseduffield/go-git/v5/plumbing/format/index"
	"github.com/jesseduffield/go-git/v5/storage/memory"
	"github.com/stretchr/testify/suite"
)

func TestIndexSuite(t *testing.T) {
	suite.Run(t, new(IndexSuite))
}

type IndexSuite struct {
	suite.Suite
}

func (s *IndexSuite) TestSetIndexBase() {
	idx := &index.Index{}
	idx.Version = 2

	base := memory.NewStorage()
	err := base.SetIndex(idx)
	s.NoError(err)

	temporal := memory.NewStorage()
	cs := NewIndexStorage(base, temporal)

	idx, err = cs.Index()
	s.NoError(err)
	s.Equal(uint32(2), idx.Version)
}

func (s *IndexSuite) TestCommit() {
	idx := &index.Index{}
	idx.Version = 2

	base := memory.NewStorage()
	err := base.SetIndex(idx)
	s.NoError(err)

	temporal := memory.NewStorage()

	idx = &index.Index{}
	idx.Version = 3

	is := NewIndexStorage(base, temporal)
	err = is.SetIndex(idx)
	s.NoError(err)

	err = is.Commit()
	s.NoError(err)

	baseIndex, err := base.Index()
	s.NoError(err)
	s.Equal(uint32(3), baseIndex.Version)
}
