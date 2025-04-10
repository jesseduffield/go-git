package merkletrie_test

import (
	"github.com/jesseduffield/go-git/v5/utils/merkletrie"
	"github.com/jesseduffield/go-git/v5/utils/merkletrie/internal/fsnoder"
	"github.com/jesseduffield/go-git/v5/utils/merkletrie/noder"

	. "gopkg.in/check.v1"
)

type ChangeSuite struct{}

var _ = Suite(&ChangeSuite{})

func (s *ChangeSuite) TestActionString(c *C) {
	action := merkletrie.Insert
	c.Assert(action.String(), Equals, "Insert")

	action = merkletrie.Delete
	c.Assert(action.String(), Equals, "Delete")

	action = merkletrie.Modify
	c.Assert(action.String(), Equals, "Modify")
}

func (s *ChangeSuite) TestUnsupportedAction(c *C) {
	a := merkletrie.Action(42)
	c.Assert(a.String, PanicMatches, "unsupported action.*")
}

func (s ChangeSuite) TestEmptyChanges(c *C) {
	ret := merkletrie.NewChanges()
	p := noder.Path{}

	err := ret.AddRecursiveInsert(p)
	c.Assert(err, Equals, merkletrie.ErrEmptyFileName)

	err = ret.AddRecursiveDelete(p)
	c.Assert(err, Equals, merkletrie.ErrEmptyFileName)
}

func (s ChangeSuite) TestNewInsert(c *C) {
	tree, err := fsnoder.New("(a(b(z<>)))")
	c.Assert(err, IsNil)
	path := find(c, tree, "z")
	change := merkletrie.NewInsert(path)
	c.Assert(change.String(), Equals, "<Insert a/b/z>")

	shortPath := noder.Path([]noder.Noder{path.Last()})
	change = merkletrie.NewInsert(shortPath)
	c.Assert(change.String(), Equals, "<Insert z>")
}

func (s ChangeSuite) TestNewDelete(c *C) {
	tree, err := fsnoder.New("(a(b(z<>)))")
	c.Assert(err, IsNil)
	path := find(c, tree, "z")
	change := merkletrie.NewDelete(path)
	c.Assert(change.String(), Equals, "<Delete a/b/z>")

	shortPath := noder.Path([]noder.Noder{path.Last()})
	change = merkletrie.NewDelete(shortPath)
	c.Assert(change.String(), Equals, "<Delete z>")
}

func (s ChangeSuite) TestNewModify(c *C) {
	tree1, err := fsnoder.New("(a(b(z<>)))")
	c.Assert(err, IsNil)
	path1 := find(c, tree1, "z")

	tree2, err := fsnoder.New("(a(b(z<1>)))")
	c.Assert(err, IsNil)
	path2 := find(c, tree2, "z")

	change := merkletrie.NewModify(path1, path2)
	c.Assert(change.String(), Equals, "<Modify a/b/z>")

	shortPath1 := noder.Path([]noder.Noder{path1.Last()})
	shortPath2 := noder.Path([]noder.Noder{path2.Last()})
	change = merkletrie.NewModify(shortPath1, shortPath2)
	c.Assert(change.String(), Equals, "<Modify z>")
}

func (s ChangeSuite) TestMalformedChange(c *C) {
	change := merkletrie.Change{}
	c.Assert(change.String, PanicMatches, "malformed change.*")
}
