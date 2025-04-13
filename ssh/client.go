package ssh

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/toyz/ssh-thing/config"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	Config      *config.SSHServer
	SSHClient   *ssh.Client
	session     *ssh.Session
	OutputChan  chan string
	ErrChan     chan error
	stdin       io.WriteCloser
	isLastCmd   bool
	initialized bool
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
		isLastCmd:  false,
	}, nil
}

func (c *Client) initSession() error {
	if c.session != nil {
		return nil
	}

	var err error
	c.session, err = c.SSHClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := c.session.RequestPty("xterm", 80, 40, modes); err != nil {
		c.session.Close()
		c.session = nil
		return fmt.Errorf("request for pseudo terminal failed: %w", err)
	}

	stdout, err := c.session.StdoutPipe()
	if err != nil {
		c.session.Close()
		c.session = nil
		return fmt.Errorf("failed to set up stdout pipe: %w", err)
	}

	stderr, err := c.session.StderrPipe()
	if err != nil {
		c.session.Close()
		c.session = nil
		return fmt.Errorf("failed to set up stderr pipe: %w", err)
	}

	c.stdin, err = c.session.StdinPipe()
	if err != nil {
		c.session.Close()
		c.session = nil
		return fmt.Errorf("failed to set up stdin pipe: %w", err)
	}

	go c.streamOutput(stdout)
	go c.streamOutput(stderr)

	if err := c.session.Shell(); err != nil {
		c.session.Close()
		c.session = nil
		return fmt.Errorf("failed to start shell: %w", err)
	}

	time.Sleep(500 * time.Millisecond)

	c.isLastCmd = false
	if _, err := c.stdin.Write([]byte("clear\n")); err != nil {
		c.ErrChan <- fmt.Errorf("warning: failed to clear terminal: %w", err)
	}

	time.Sleep(300 * time.Millisecond)

	c.initialized = true
	return nil
}

func (c *Client) RunCommand(command string) {
	go func() {
		if err := c.initSession(); err != nil {
			c.ErrChan <- err
			return
		}

		if !strings.HasSuffix(command, "\n") {
			command = command + "\n"
		}

		if _, err := c.stdin.Write([]byte(command)); err != nil {
			c.ErrChan <- fmt.Errorf("failed to send command: %w", err)

			c.Close()
			c.session = nil
			return
		}
	}()
}

func (c *Client) RunCommands(commands []string) {
	if len(commands) == 0 {
		return
	}

	go func() {
		if err := c.initSession(); err != nil {
			c.ErrChan <- err
			return
		}

		if !c.initialized {
			time.Sleep(200 * time.Millisecond)
		}

		for i, cmd := range commands {
			isLastCmd := i == len(commands)-1

			c.isLastCmd = isLastCmd

			if !strings.HasSuffix(cmd, "\n") {
				cmd = cmd + "\n"
			}

			if i > 0 {
				time.Sleep(500 * time.Millisecond)
			}

			if _, err := c.stdin.Write([]byte(cmd)); err != nil {
				c.ErrChan <- fmt.Errorf("failed to send command: %w", err)

				c.Close()
				c.session = nil
				return
			}
		}
	}()
}

func (c *Client) Close() error {
	if c.session != nil {
		c.session.Close()
		c.session = nil
	}

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
			if c.isLastCmd {
				c.OutputChan <- string(buf[:n])
			}
		}
	}
}
