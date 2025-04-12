package ssh

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/toyz/ssh-thing/config"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	Config     *config.SSHServer
	SSHClient  *ssh.Client
	OutputChan chan string
	ErrChan    chan error
}

func NewClient(sshConfig *config.SSHServer) (*Client, error) {
	var authMethod ssh.AuthMethod

	if sshConfig.PrivateKeyPath != "" {
		key, err := os.ReadFile(sshConfig.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}

		authMethod = ssh.PublicKeys(signer)
	} else if sshConfig.Password != "" {
		authMethod = ssh.Password(sshConfig.Password)
	} else {
		return nil, fmt.Errorf("authentication failed: neither private key path nor password provided")
	}

	config := &ssh.ClientConfig{
		User: sshConfig.User,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	addr := fmt.Sprintf("%s:%d", sshConfig.Host, sshConfig.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &Client{
		Config:     sshConfig,
		SSHClient:  client,
		OutputChan: make(chan string),
		ErrChan:    make(chan error),
	}, nil
}

func (c *Client) RunCommand(command string) {
	go func() {
		session, err := c.SSHClient.NewSession()
		if err != nil {
			c.ErrChan <- fmt.Errorf("failed to create session: %w", err)
			return
		}
		defer session.Close()

		stdout, err := session.StdoutPipe()
		if err != nil {
			c.ErrChan <- fmt.Errorf("failed to set up stdout pipe: %w", err)
			return
		}

		stderr, err := session.StderrPipe()
		if err != nil {
			c.ErrChan <- fmt.Errorf("failed to set up stderr pipe: %w", err)
			return
		}

		if err := session.Start(command); err != nil {
			c.ErrChan <- fmt.Errorf("failed to start command: %w", err)
			return
		}

		go c.streamOutput(stdout)
		go c.streamOutput(stderr)

		err = session.Wait()
		if err != nil && err != io.EOF {
			c.ErrChan <- fmt.Errorf("command execution error: %w", err)
		}
	}()
}

func (c *Client) Close() error {
	if c.SSHClient != nil {
		return c.SSHClient.Close()
	}
	return nil
}

func (c *Client) streamOutput(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				c.ErrChan <- fmt.Errorf("read error: %w", err)
			}
			break
		}
		if n > 0 {
			c.OutputChan <- string(buf[:n])
		}
	}
}
