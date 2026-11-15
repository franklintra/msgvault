package imap

import (
	"log/slog"
	"testing"

	"github.com/emersion/go-imap/v2"
	"github.com/stretchr/testify/assert"
)

// A failed reconnect leaves c.conn nil. The label map then moved on to the
// next mailbox and SELECT dereferenced the nil client, which panicked the
// whole sync (seen with DavMail restarting under a first sync).
func TestCommandsWithoutConnectionFailAsNetworkErrors(t *testing.T) {
	assert := assert.New(t)
	client := &Client{logger: slog.Default()}

	err := client.selectMailbox("INBOX")
	assert.ErrorIs(err, errNoConnection)
	assert.True(isNetworkError(err), "a missing connection must read as a network error")

	_, err = client.connUIDSearch(&imap.SearchCriteria{})
	assert.ErrorIs(err, errNoConnection)
	_, err = client.connFetch(imap.UIDSet{}, &imap.FetchOptions{})
	assert.ErrorIs(err, errNoConnection)
	_, err = client.connStatus("INBOX", &imap.StatusOptions{})
	assert.ErrorIs(err, errNoConnection)
	_, err = client.connList(client.listOptionsLocked())
	assert.ErrorIs(err, errNoConnection)
	assert.False(client.enableQresync())
}
