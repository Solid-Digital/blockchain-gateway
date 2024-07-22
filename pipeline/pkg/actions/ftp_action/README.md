# FTP Action

FTP action is a wrapper for an FTP client with a structured interface so it can be conveniently used in pipelines.

## Functions

- Retrieve all files from directory
- Delete file
- Move file

## Input

- address
The address of the ftp server including port, e.g. `ftp.example.com:22`
- username
The username for logging in.
- password
The password used for logging in.
- serverKey
The SSH public key of the server.
- function
The function to use, options: `retrieve`, `move` and `delete`.
- dir
The path to the directory you want to retrieve files from. Only applicable to the retrieve function.
- filepath
The filepath of the file you want to move or delete.
- target
The target filepath you want to move the source file to. Only applicable to `move`.

## Output
Only `retrieve` outputs a map of filename to file content as:
```
map[string]interface{}{
	"exampleFile.txt": []byte{A, B, C}
}

