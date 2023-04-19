package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"fedlearn/psi/common/log"

	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context,
	binpath,
	psiProtocol,
	role,
	inputFilePath,
	ip,
	port,
	resultFilePath string,
	limit int64,
	serverAcceptTimeout int32,
	ch chan int) error {
	if serverAcceptTimeout <= 0 {
		serverAcceptTimeout = 300
	}

	log.Infof(ctx, "psi start: role=%s, input=%s, ip=%s, port=%s, result=%s", role, inputFilePath, ip, port, resultFilePath)
	//check role
	if role == "server" {
		role = "0"
	} else if role == "client" {
		role = "1"
	} else {
		ch <- -1
		return errors.New("Param 'role' cannot be recognized, we only accept 'client' or 'server'.")
	}

	//check inputFilePath
	if inputFilePath == "" {
		ch <- -1
		return errors.New("Empty inputFilePath.")
	} else {
		_, err := os.Lstat(inputFilePath)
		if os.IsNotExist(err) {
			ch <- -1
			return errors.New("Input file does not exists!")
		}
	}

	//check resultFilePath
	if resultFilePath == "" {
		ch <- -1
		return errors.New("Empty resultFilePath.")
	}

	//check ip
	if ip == "" {
		ip = "127.0.0.1"
	}

	//check port
	if port == "" {
		port = "7766"
	}

	// do not support 3600+
	if serverAcceptTimeout > 3600 {
		serverAcceptTimeout = 3600
	}

	// at least 1 seconds
	if serverAcceptTimeout < 0 {
		serverAcceptTimeout = 1
	}

	cmd := exec.Command(binpath,
		"-r", role,
		"-p", psiProtocol,
		"-f", inputFilePath,
		"-l", fmt.Sprintf("%d", limit),
		"-a", ip,
		"-o", port,
		"-u", fmt.Sprintf("%d", serverAcceptTimeout),
		resultFilePath)

	log.Infoln(ctx, cmd.Path)
	log.Infoln(ctx, cmd.Args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	done := make(chan struct{})
	errCh := make(chan error)

	//due to cmd.Run() will block this goroutine
	//we cannot return pid of cmd until it ends
	//so we put it in another goroutine
	go func() {
		//after Stop() executed, cmd.Run() will be finished and return an error
		//this goroutine will quit safely without worry about goroutine leak
		err := cmd.Run()

		errCh <- err
		done <- struct{}{}
	}()

	//waiting the cmd.Run() finish its init work
	for cmd.Process == nil {
		select {
		case err := <-errCh:
			if err != nil {
				ch <- -1
				return err
			}
		default:
		}
	}
	ch <- cmd.Process.Pid

	//waiting the cmd.Run() finish all work
	err := <-errCh
	if err != nil {
		//we can't pass -1 to ch like below, or we'll hang forever here.
		//ch <- -1
		return err
	}
	<-done

	// Write stdout into result file
	if f, err := os.OpenFile(resultFilePath, os.O_WRONLY|os.O_CREATE, 0776); err != nil {
		logrus.Errorln("Open result file error:", err)
	} else {
		f.Write(stdout.Bytes())
		f.Close()
	}

	// Write stderr into error file
	if errFile, err := os.OpenFile(filepath.Join(filepath.Dir(resultFilePath), "error.out"), os.O_WRONLY|os.O_CREATE, 0776); err != nil {
		logrus.Errorln("Open err file error:", err)
	} else {
		errFile.Write(stderr.Bytes())
		errFile.Close()
	}

	return nil
}

func Stop(ctx context.Context, pid int) error {
	log.Infof(ctx, "psi stop: pid=%d", pid)
	//check process existence
	if err := syscall.Kill(pid, 0); err != nil {
		return errors.New("this process does not exists!")
	}

	err := syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		return errors.New("failed to kill this process")
	}

	return nil
}
