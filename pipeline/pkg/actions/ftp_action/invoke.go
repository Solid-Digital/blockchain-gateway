package ftp_action

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/errors"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

const (
	// config
	Address   = "address"
	Username  = "username"
	Password  = "password"
	ServerKey = "serverKey"
	// functions
	Function = "function"
	Retrieve = "retrieve"
	Delete   = "delete"
	Move     = "move"
	// parameters
	FilePath    = "filepath"
	Directory   = "dir"
	Target      = "target"
	OutputFiles = "outputFiles"
)

func Invoke(stub domain.Stub, input map[string]interface{}) (output map[string]interface{}, err error) {
	addr, ok := input[Address].(string)
	if !ok {
		return nil, errors.New("could not cast address to string")
	}
	username, ok := input[Username].(string)
	if !ok {
		return nil, errors.New("could not cast username to string")
	}
	password, ok := input[Password].(string)
	if !ok {
		return nil, errors.New("could not cast password to string")
	}
	serverKey, ok := input[ServerKey].(string)
	if !ok {
		return nil, errors.New("could not cast server key to string")
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: trustedHostKeyCallback(serverKey),
		Timeout:         5 * time.Second,
	}
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	stub.Debugf("connected to sftp server on %s", addr)
	defer client.Close()

	function, ok := input[Function].(string)
	if !ok {
		return nil, errors.New("could not cast function to string")
	}

	switch function {
	case Retrieve:
		dir, ok := input[Directory].(string)
		if !ok {
			return nil, errors.New("could not cast dir to string")
		}
		walker := client.Walk(dir)

		files := make(map[string][]byte)
		for walker.Step() {
			if walker.Err() != nil {
				return nil, err
			}

			// skip if step is dir
			fi := walker.Stat()
			if fi.IsDir() {
				continue
			}

			fileBytes, err := handleFile(stub, client, walker.Path(), fi)
			if err != nil {
				return nil, err
			}
			files[walker.Path()] = fileBytes
		}
		client.Close()

		return map[string]interface{}{
			OutputFiles: files,
		}, nil
	case Delete:
		filepath, ok := input[FilePath].(string)
		if !ok {
			return nil, errors.New("could not cast filepath to string")
		}
		err := client.Remove(filepath)
		if err != nil {
			return nil, err
		}
		err = client.Close()
		if err != nil {
			stub.Errorf("%s", err)
		}

		return nil, nil
	case Move:
		filepath, ok := input[FilePath].(string)
		if !ok {
			return nil, errors.New("could not cast filepath to string")
		}
		target, ok := input[Target].(string)
		if !ok {
			return nil, errors.New("could not cast target to string")
		}
		err := client.Rename(filepath, target)
		if err != nil {
			return nil, err
		}
		client.Close()
		return nil, nil
	default:
		return nil, errors.New("unknown function, please refer to the docs")
	}
}

// create human-readable SSH-key strings
func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal()) // e.g. "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTY...."
}

func trustedHostKeyCallback(trustedKey string) ssh.HostKeyCallback {

	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			return nil
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedKey, ks)
		}

		return nil
	}
}

func handleFile(stub domain.Stub, client *sftp.Client, path string, fi os.FileInfo) ([]byte, error) {
	// get file
	file, err := client.Open(path)
	if err != nil {
		return nil, err
	}

	fileBytes := make([]byte, fi.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
