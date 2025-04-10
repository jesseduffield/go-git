package packp

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/jesseduffield/go-git/v5/plumbing"

	. "gopkg.in/check.v1"
)

type ServerResponseSuite struct{}

var _ = Suite(&ServerResponseSuite{})

func (s *ServerResponseSuite) TestDecodeNAK(c *C) {
	raw := "0008NAK\n"

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), false)
	c.Assert(err, IsNil)

	c.Assert(sr.ACKs, HasLen, 0)
}

func (s *ServerResponseSuite) TestDecodeNewLine(c *C) {
	raw := "\n"

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), false)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "invalid pkt-len found")
}

func (s *ServerResponseSuite) TestDecodeEmpty(c *C) {
	raw := ""

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), false)
	c.Assert(err, IsNil)
}

func (s *ServerResponseSuite) TestDecodePartial(c *C) {
	raw := "000600\n"

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), false)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, fmt.Sprintf("unexpected content %q", "00"))
}

func (s *ServerResponseSuite) TestDecodeACK(c *C) {
	raw := "0031ACK 6ecf0ef2c2dffb796033e5a02219af86ec6584e5\n"

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), false)
	c.Assert(err, IsNil)

	c.Assert(sr.ACKs, HasLen, 1)
	c.Assert(sr.ACKs[0], Equals, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))
}

func (s *ServerResponseSuite) TestDecodeMultipleACK(c *C) {
	raw := "" +
		"0031ACK 1111111111111111111111111111111111111111\n" +
		"0031ACK 6ecf0ef2c2dffb796033e5a02219af86ec6584e5\n" +
		"00080PACK\n"

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), false)
	c.Assert(err, IsNil)

	c.Assert(sr.ACKs, HasLen, 2)
	c.Assert(sr.ACKs[0], Equals, plumbing.NewHash("1111111111111111111111111111111111111111"))
	c.Assert(sr.ACKs[1], Equals, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))
}

func (s *ServerResponseSuite) TestDecodeMultipleACKWithSideband(c *C) {
	raw := "" +
		"0031ACK 1111111111111111111111111111111111111111\n" +
		"0031ACK 6ecf0ef2c2dffb796033e5a02219af86ec6584e5\n" +
		"00080aaaa\n"

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), false)
	c.Assert(err, IsNil)

	c.Assert(sr.ACKs, HasLen, 2)
	c.Assert(sr.ACKs[0], Equals, plumbing.NewHash("1111111111111111111111111111111111111111"))
	c.Assert(sr.ACKs[1], Equals, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))
}

func (s *ServerResponseSuite) TestDecodeMalformed(c *C) {
	raw := "0029ACK 6ecf0ef2c2dffb796033e5a02219af86ec6584e\n"

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), false)
	c.Assert(err, NotNil)
}

// multi_ack isn't fully implemented, this ensures that Decode ignores that fact,
// as in some circumstances that's OK to assume so.
//
// TODO: Review as part of multi_ack implementation.
func (s *ServerResponseSuite) TestDecodeMultiACK(c *C) {
	raw := "" +
		"0031ACK 1111111111111111111111111111111111111111\n" +
		"0031ACK 6ecf0ef2c2dffb796033e5a02219af86ec6584e5\n" +
		"00080PACK\n"

	sr := &ServerResponse{}
	err := sr.Decode(bufio.NewReader(bytes.NewBufferString(raw)), true)
	c.Assert(err, IsNil)

	c.Assert(sr.ACKs, HasLen, 2)
	c.Assert(sr.ACKs[0], Equals, plumbing.NewHash("1111111111111111111111111111111111111111"))
	c.Assert(sr.ACKs[1], Equals, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))
}
