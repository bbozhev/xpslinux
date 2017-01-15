package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// sh runs a shell command. Any output will be displayed on the screen.
func sh(args ...string) error {
	if len(args) == 0 {
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run: %v: %s", args, err)
	}
	return nil
}

// getSh runs a shell command. Stdout and stderr are returned as strings.
func getSh(args ...string) (string, string, error) {
	if len(args) == 0 {
		return "", "", nil
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to run: %v: %s", args, err)
	}

	return outBuf.String(), errBuf.String(), nil
}

// cp copies a file. If dst is a directory, then src will be copied into that directory.
func cp(src, dst string) error {
	dstFi, err := os.Stat(dst)
	if err == nil {
		if dstFi.IsDir() {
			dst = filepath.Join(dst, filepath.Base(src))
		}
	} else if os.IsNotExist(err) {
		// Fine. Keep going.
	} else if err != nil {
		return err
	}

	sfd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sfd.Close()

	dfd, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dfd.Close()

	if _, err = io.Copy(dfd, sfd); err != nil {
		return err
	}

	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sfi.Mode())
}

// fwrite writes data at a given path. Data is appended with a newline. If name
// exists, it will be overwritten.
func fwrite(name, data string) error {
	return ioutil.WriteFile(name, []byte(data+"\n"), 0644)
}

// fappend writes data at the end of a given file path. If name doesn't exist,
// then a new file is created. fappend adds newlines to data automatically.
func fappend(name, data string) error {
	fi, err := os.Stat(name)
	if os.IsNotExist(err) {
		return fwrite(name, data)
	} else if err != nil {
		return err
	}

	fdata, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(fdata)
	if !bytes.HasSuffix(fdata, []byte("\n")) {
		buf.WriteString("\n")
	}
	buf.WriteString(data)
	buf.WriteString("\n")

	return ioutil.WriteFile(name, buf.Bytes(), fi.Mode())
}

// curlo downloads a file to a given directory.
func curlo(outPath, u string) error {
	ofi, err := os.Stat(outPath)
	if os.IsNotExist(err) {
		// Fine.
	} else if err != nil {
		return err
	} else if ofi.IsDir() {
		return fmt.Errorf("failed to create file %q: is a directory", outPath)
	}

	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func clear() error {
	return sh("clear")
}

func chownR(root string, uid, gid int) error {
	return filepath.Walk(root, func(name string, fi os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}
