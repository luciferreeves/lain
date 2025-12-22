package email

import (
	"crypto/tls"
	"fmt"
	"lain/config"
	"lain/types"

	"github.com/emersion/go-imap/client"
)

func ConnectIMAP(email, password string) (*types.EmailClient, error) {
	address := fmt.Sprintf("%s:%d", config.MailServer.IMAPHost, config.MailServer.IMAPPort)

	var c *client.Client
	var err error

	if config.MailServer.IMAPTLS {
		c, err = client.DialTLS(address, &tls.Config{
			ServerName: config.MailServer.IMAPHost,
		})
	} else {
		c, err = client.Dial(address)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server: %w", err)
	}

	if err := c.Login(email, password); err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	return &types.EmailClient{Client: c}, nil
}

func DisconnectIMAP(c *types.EmailClient) error {
	return c.Logout()
}
