package transactional

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/jesseduffield/go-git/v5/plumbing"
	"github.com/jesseduffield/go-git/v5/plumbing/cache"
	"github.com/jesseduffield/go-git/v5/plumbing/storer"
	"github.com/jesseduffield/go-git/v5/storage"
	"github.com/jesseduffield/go-git/v5/storage/filesystem"
	"github.com/jesseduffield/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommit(t *testing.T) {
	base := memory.NewStorage()
	temporal := filesystem.NewStorage(memfs.New(), cache.NewObjectLRUDefault())
	st := NewStorage(base, temporal)

	commit := base.NewEncodedObject()
	commit.SetType(plumbing.CommitObject)

	_, err := st.SetEncodedObject(commit)
	require.NoError(t, err)

	ref := plumbing.NewHashReference("refs/a", commit.Hash())
	require.NoError(t, st.SetReference(ref))

	err = st.Commit()
	require.NoError(t, err)

	ref, err = base.Reference(ref.Name())
	require.NoError(t, err)
	assert.Equal(t, commit.Hash(), ref.Hash())

	obj, err := base.EncodedObject(plumbing.AnyObject, commit.Hash())
	require.NoError(t, err)
	assert.Equal(t, commit.Hash(), obj.Hash())
}

func TestTransactionalPackfileWriter(t *testing.T) {
	base := memory.NewStorage()
	var temporal storage.Storer

	temporal = filesystem.NewStorage(memfs.New(), cache.NewObjectLRUDefault())

	st := NewStorage(base, temporal)

	_, tmpOK := temporal.(storer.PackfileWriter)
	_, ok := st.(storer.PackfileWriter)
	assert.Equal(t, tmpOK, ok)
}
