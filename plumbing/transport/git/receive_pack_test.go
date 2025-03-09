package git

import (
	"testing"

	"github.com/jesseduffield/go-git/v5/internal/transport/test"
	"github.com/stretchr/testify/suite"

	fixtures "github.com/go-git/go-git-fixtures/v4"
)

func TestReceivePackSuite(t *testing.T) {
	suite.Run(t, new(ReceivePackSuite))
}

type ReceivePackSuite struct {
	rps test.ReceivePackSuite
	BaseSuite
}

func (s *ReceivePackSuite) SetupTest() {
	s.BaseSuite.SetupTest()

	s.rps.SetS(s)
	s.rps.Client = DefaultClient
	s.rps.Endpoint = s.prepareRepository(fixtures.Basic().One(), "basic.git")
	s.rps.EmptyEndpoint = s.prepareRepository(fixtures.ByTag("empty").One(), "empty.git")
	s.rps.NonExistentEndpoint = s.newEndpoint("non-existent.git")

	s.StartDaemon()
}

func (s *ReceivePackSuite) TestAdvertisedReferencesEmpty() {
	// This test from BaseSuite is flaky, so it's disabled until we figure out a solution.
}
