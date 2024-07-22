package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/unchainio/pkg/xexec"

	"github.com/unchainio/pkg/iferr"
)

func main() {
	var err error
	if len(os.Args) < 3 {
		log.Fatal("expecting 2 arguments: the path to the generated handlers and the path to the implementations")
	}
	handlersGenPath := os.Args[1]
	handlersImplPath := os.Args[2]
	handlersOldPath := handlersGenPath + "_old"

	_, err = xexec.Run("mkdir -p %s", xexec.WithFormat(handlersOldPath))
	iferr.Panic(err)

	_, err = xexec.Run("mkdir -p %s", xexec.WithFormat(handlersImplPath))
	iferr.Panic(err)

	//os.Open()
	err = filepath.Walk(handlersGenPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fmt.Println(path)
		relGen, err := filepath.Rel(handlersGenPath, path)
		if err != nil {
			return err
		}

		fmt.Println(relGen)

		implPath := filepath.Join(handlersImplPath, relGen)
		oldPath := filepath.Join(handlersOldPath, relGen)
		fmt.Println(implPath)

		_, err = xexec.Run("touch %s %s", xexec.WithFormat(implPath, oldPath))
		if err != nil {
			return err
		}

		//out, err := xexec.Run("diff3 %s %s %s --easy-only --merge > %s", xexec.WithFormat(implPath, oldPath, path, implPath))
		out, err := xexec.Run("git merge-file %s %s %s --union", xexec.WithFormat(implPath, oldPath, path))
		if err != nil {
			return err
		}

		fmt.Printf(string(out))

		return nil
	})

	iferr.Panic(err)

	_, err = xexec.Run("rm -rf %s", xexec.WithFormat(handlersOldPath))
	iferr.Panic(err)

	_, err = xexec.Run("cp -r %s %s", xexec.WithFormat(handlersGenPath, handlersOldPath))
	iferr.Panic(err)

	fmt.Println(handlersGenPath)
	fmt.Println(handlersImplPath)
}
