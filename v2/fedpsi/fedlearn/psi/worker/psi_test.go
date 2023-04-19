package worker

import (
    "testing"
    "time"
)

func TestStartNormalCase(t *testing.T) {
    serverInputFile := "./test_data/a.txt"
    clientInputFile := "./test_data/b.txt"
    serverResultFile := "./test_data/server_result.txt"
    clientResultFile := "./test_data/client_result.txt"
    ip := ""
    port := ""
    serverPidCh := make(chan int)
    serverErrCh := make(chan error)
    clientPidCh := make(chan int)
    clientErrCh := make(chan error)
    
    //start server
    go func() {
        err := Start("server", serverInputFile, ip, port, serverResultFile, serverPidCh)
        serverErrCh <- err
    }()
    serverPid := <-serverPidCh

    //we will not wait serverErrCh, otherwise client can't be launched
    var serverErr error
    select {
    case serverErr = <- serverErrCh:
    default:
    }
    if serverPid == -1 || serverErr != nil {
        t.Errorf("Failed to start server.")
    }

    //start client
    go func() {
        err := Start("client", clientInputFile, ip, port, clientResultFile, clientPidCh)
        clientErrCh <- err
    }()

    clientPid := <-clientPidCh
    println(clientPid)
    clientErr := <- clientErrCh

    if clientPid == -1 || clientErr != nil {
        t.Errorf("Failed to start client.")
    }
}

func TestStartWithWrongInputFilePath(t *testing.T) {
    clientInputFile := "../wrong/path/b.txt"
    clientResultFile := "./test_data/client_result.txt"
    ip := ""
    port := ""
    ch := make(chan int)
    errCh := make(chan error)

    go func() {
        err := Start("client", clientInputFile, ip, port, clientResultFile, ch)
        errCh <- err
    }()

    pid := <-ch

    err := <-errCh
    if pid != -1 || err == nil {
        t.Errorf("Failed to check input file which does not exists.")
    }
}

func TestStartWithWrongResultFilePath(t *testing.T) {
    clientInputFile := "./test_data/b.txt"
    clientResultFile := "../wrong/path/result.txt"
    ip := ""
    port := ""
    ch := make(chan int)
    errCh := make(chan error)

    go func() {
        err := Start("client", clientInputFile, ip, port, clientResultFile, ch)
        errCh <- err
    }()

    pid := <-ch

    err := <-errCh
    if pid != -1 || err == nil {
        t.Errorf("Failed to check result file which does not exists.")
    }
}

func TestStopSuccessfulCase(t *testing.T) {
    serverInputFile := "./test_data/a.txt"
    serverResultFile := "./test_data/server_result.txt"
    ip := ""
    port := ""
    serverPidCh := make(chan int)
    serverErrCh := make(chan error)
    
    //start server
    go func() {
        err := Start("server", serverInputFile, ip, port, serverResultFile, serverPidCh)
        serverErrCh <- err
    }()
    serverPid := <-serverPidCh

    //we will not wait serverErrCh, otherwise Stop() cannot be tested
    var serverErr error
    select {
    case serverErr = <- serverErrCh:
    default:
    }
    if serverPid == -1 || serverErr != nil {
        t.Errorf("Failed to start server.")
    }

    err := Stop(serverPid)
    if err != nil {
        t.Errorf("Failed to stop server.")
    }
}

//stop a server have already be stopped
func TestStopFailedCase1(t *testing.T) {
    serverInputFile := "./test_data/a.txt"
    serverResultFile := "./test_data/server_result.txt"
    ip := ""
    port := ""
    serverPidCh := make(chan int)
    serverErrCh := make(chan error)
    
    //start server
    go func() {
        err := Start("server", serverInputFile, ip, port, serverResultFile, serverPidCh)
        serverErrCh <- err
    }()
    serverPid := <-serverPidCh

    //we will not wait serverErrCh, otherwise Stop() cannot be tested
    var serverErr error
    select {
    case serverErr = <- serverErrCh:
    default:
    }
    if serverPid == -1 || serverErr != nil {
        t.Errorf("Failed to start server.")
    }

    err := Stop(serverPid)
    if err != nil {
        t.Errorf("Failed to stop server.")
    }

    time.Sleep(time.Second)

    err = Stop(serverPid)
    if err == nil {
        t.Errorf("Failed to detect the scene that stopping a stopped server.")
    }
}

//stop a process which does not exists
func TestStopFailedCase2(t *testing.T) {
    serverInputFile := "./test_data/a.txt"
    serverResultFile := "./test_data/server_result.txt"
    ip := ""
    port := ""
    serverPidCh := make(chan int)
    serverErrCh := make(chan error)
    
    //start server
    go func() {
        err := Start("server", serverInputFile, ip, port, serverResultFile, serverPidCh)
        serverErrCh <- err
    }()
    serverPid := <-serverPidCh

    //we will not wait serverErrCh, otherwise Stop() cannot be tested
    var serverErr error
    select {
    case serverErr = <- serverErrCh:
    default:
    }
    if serverPid == -1 || serverErr != nil {
        t.Errorf("Failed to start server.")
    }

    wrongPid := 12345678
    err := Stop(wrongPid)
    if err == nil {
        t.Errorf("Failed to detect the scene that stopping a stopped server.")
    }
}
