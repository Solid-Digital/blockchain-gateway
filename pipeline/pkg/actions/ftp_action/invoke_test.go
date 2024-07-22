package ftp_action_test

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/ftp_action"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/helper"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

func (s *TestSuite) TestFtpAction_InvokeRetrieve() {
	// start local ftp server
	s.createTestFile("./testdata/files/test.txt")
	go startSftpServer()
	time.Sleep(1 * time.Second)

	cases := map[string]struct {
		Stub           domain.Stub
		Input          map[string]interface{}
		Success        bool
		ExpectedOutput map[string]interface{}
	}{
		"invoke ftp action retrieve successfully": {
			s.logger,
			map[string]interface{}{
				"address":   "localhost:2022",
				"username":  "user",
				"password":  "password",
				"function":  "retrieve",
				"dir":       "./testdata/files",
				"serverKey": s.getServerKey(),
			},
			true,
			map[string]interface{}{"outputFiles": map[string][]uint8{"testdata/files/test.txt": []uint8{0x41, 0x42, 0x43}}},
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			output, err := ftp_action.Invoke(tc.Stub, tc.Input)

			if tc.Success {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedOutput, output)
			} else {
				require.Error(t, err)
			}
		})
	}
	s.removeDir("./testdata/files")
}

func (s *TestSuite) TestFtpAction_InvokeRemove() {
	s.createTestFile("./testdata/files/delme.txt")

	// start local ftp server
	go startSftpServer()
	time.Sleep(1 * time.Second)

	cases := map[string]struct {
		Stub           domain.Stub
		Input          map[string]interface{}
		Success        bool
		ExpectedOutput map[string]interface{}
	}{
		"invoke ftp action delete successfully": {
			s.logger,
			map[string]interface{}{
				"address":   "localhost:2022",
				"username":  "user",
				"password":  "password",
				"function":  "delete",
				"filepath":  "./testdata/files/delme.txt",
				"serverKey": s.getServerKey(),
			},
			true,
			nil,
		},
		"invoke ftp action delete errors with unknown file": {
			s.logger,
			map[string]interface{}{
				"address":   "localhost:2022",
				"username":  "user",
				"password":  "password",
				"function":  "delete",
				"filepath":  "./testdata/files/random_file.txt",
				"serverKey": s.getServerKey(),
			},
			false,
			nil,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			output, err := ftp_action.Invoke(tc.Stub, tc.Input)

			if tc.Success {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedOutput, output)
			} else {
				require.Error(t, err)
			}
		})
	}

	s.removeDir("./testdata/files")
}

func (s *TestSuite) TestFtpAction_InvokeMove() {
	s.createTestFile("./testdata/files/failed.txt")

	// start local ftp server
	go startSftpServer()
	time.Sleep(1 * time.Second)

	cases := map[string]struct {
		Stub           domain.Stub
		Input          map[string]interface{}
		Success        bool
		ExpectedOutput map[string]interface{}
	}{
		"invoke ftp action move successfully": {
			s.logger,
			map[string]interface{}{
				"address":   "localhost:2022",
				"username":  "user",
				"password":  "password",
				"function":  "move",
				"filepath":  "./testdata/files/failed.txt",
				"target":    "./testdata/files/failed/failed.txt",
				"serverKey": s.getServerKey(),
			},
			true,
			nil,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			output, err := ftp_action.Invoke(tc.Stub, tc.Input)

			if tc.Success {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedOutput, output)
			} else {
				require.Error(t, err)
			}
		})
	}

	s.removeDir("./testdata/files")
}

func (s *TestSuite) createTestFile(testFile string) {
	err := os.Mkdir("./testdata/files", 0755)
	require.NoError(s.T(), err, "could not create dir testdata/files")
	err = os.Mkdir("./testdata/files/failed", 0755)
	require.NoError(s.T(), err, "could not create dir testdata/files/failed")

	file, err := os.Create(testFile)
	require.NoError(s.T(), err, "could not create file %s", testFile)

	_, err = file.Write([]byte("ABC"))
	require.NoError(s.T(), err, "could not write to testfile")

	err = file.Close()
	require.NoError(s.T(), err, "could not close testfile")
}

func (s *TestSuite) removeDir(dir string) {
	err := os.RemoveAll(dir)
	require.NoError(s.T(), err, "could not remove dir %s", dir)
}

// func startSftpServer(stopChannel chan bool) {
func startSftpServer() {
	debugStream := ioutil.Discard

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			fmt.Fprintf(debugStream, "Login: %s\n", c.User())
			if c.User() == "user" && string(pass) == "password" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	privateBytes, err := ioutil.ReadFile("./testdata/id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	defer listener.Close()

	nConn, err := listener.Accept()
	if err != nil {
		log.Fatal("failed to accept incoming connection", err)
	}

	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Fatal("failed to handshake", err)
	}
	fmt.Fprintf(debugStream, "SSH server established\n")

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of an SFTP session, this is "subsystem"
		// with a payload string of "<length=4>sftp"
		fmt.Fprintf(debugStream, "Incoming channel: %s\n", newChannel.ChannelType())
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			fmt.Fprintf(debugStream, "Unknown channel type: %s\n", newChannel.ChannelType())
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Fatal("could not accept channel.", err)
		}
		fmt.Fprintf(debugStream, "Channel accepted\n")

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "subsystem" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Fprintf(debugStream, "Request: %v\n", req.Type)
				ok := false
				switch req.Type {
				case "subsystem":
					fmt.Fprintf(debugStream, "Subsystem: %s\n", req.Payload[4:])
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				fmt.Fprintf(debugStream, " - accepted: %v\n", ok)
				req.Reply(ok, nil)
			}
		}(requests)

		serverOptions := []sftp.ServerOption{
			sftp.WithDebug(debugStream),
		}

		server, err := sftp.NewServer(
			channel,
			serverOptions...,
		)
		if err != nil {
			log.Fatal(err)
		}
		defer server.Close()

		if err := server.Serve(); err == io.EOF {
			log.Print("sftp client exited session.")
			server.Close()
		} else if err != nil {
			log.Fatal("sftp server completed with error:", err)
		}
	}
}

func (s *TestSuite) getServerKey() string {
	helper := helper.NewHelper(&s.Suite, s.logger)
	bytes := helper.BytesFromFile("./testdata/id_rsa.pub")
	return string(bytes)
}
