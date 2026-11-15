package imap

import (
	"fmt"
	"net"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

// errNoConnection reports a command issued after a failed reconnect left the
// client without a connection. It wraps net.ErrClosed so isNetworkError
// treats it like the dropped connection it follows: callers reconnect or
// skip the mailbox instead of dereferencing a nil client.
var errNoConnection = fmt.Errorf("IMAP connection is not open: %w", net.ErrClosed)

// The conn* helpers run one command on c.conn, or fail with errNoConnection
// when a failed reconnect left it nil. Caller must hold mu.

func (c *Client) connCaps() imap.CapSet {
	if c.conn == nil {
		return nil
	}
	return c.conn.Caps()
}

func (c *Client) connSelect(mailbox string, options *imap.SelectOptions) (*imap.SelectData, error) {
	if c.conn == nil {
		return nil, errNoConnection
	}
	return c.conn.Select(mailbox, options).Wait()
}

func (c *Client) connList(options *imap.ListOptions) ([]*imap.ListData, error) {
	if c.conn == nil {
		return nil, errNoConnection
	}
	return c.conn.List("", "*", options).Collect()
}

func (c *Client) connStatus(mailbox string, options *imap.StatusOptions) (*imap.StatusData, error) {
	if c.conn == nil {
		return nil, errNoConnection
	}
	return c.conn.Status(mailbox, options).Wait()
}

func (c *Client) connUIDSearch(criteria *imap.SearchCriteria) (*imap.SearchData, error) {
	if c.conn == nil {
		return nil, errNoConnection
	}
	return c.conn.UIDSearch(criteria, nil).Wait()
}

func (c *Client) connFetch(uids imap.NumSet, options *imap.FetchOptions) ([]*imapclient.FetchMessageBuffer, error) {
	if c.conn == nil {
		return nil, errNoConnection
	}
	return c.conn.Fetch(uids, options).Collect()
}

func (c *Client) connEnable() (*imapclient.EnableData, error) {
	if c.conn == nil {
		return nil, errNoConnection
	}
	return c.conn.Enable(imap.CapQResync).Wait()
}
