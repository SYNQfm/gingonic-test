package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
)

type ExecObj struct {
	Cmd     *exec.Cmd
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	stdin   io.WriteCloser
	ReadOut []byte
	ReadErr []byte
	Err     error
}

func NewExec(command string) ExecObj {
	cmd := exec.Command(command)
	obj := ExecObj{Cmd: cmd}
	obj.Err = obj.Open()
	return obj
}

func (e *ExecObj) Open() error {
	stdin, err := e.Cmd.StdinPipe()
	e.stdin = stdin
	if err != nil {
		return err
	}
	stdout, err := e.Cmd.StdoutPipe()
	e.stdout = stdout
	if err != nil {
		return err
	}

	stderr, err := e.Cmd.StderrPipe()
	e.stderr = stderr
	if err != nil {
		return err
	}
	return nil
}

func (e *ExecObj) Close() {
	e.stderr.Close()
	e.stdout.Close()
}

func (e *ExecObj) Read() error {
	o, err := ioutil.ReadAll(e.stdout)
	if err != nil {
		return err
	}
	e.ReadOut = o
	o2, err := ioutil.ReadAll(e.stderr)
	if err != nil {
		return err
	}
	e.ReadErr = o2
	return nil
}

func (e *ExecObj) Exec(fn func(io.WriteCloser)) {
	e.Open()
	if err := e.Cmd.Start(); err != nil {
		e.Err = err
		return
	}

	fn(e.stdin)

	if err := e.stdin.Close(); err != nil {
		e.Err = err
		return
	}

	if err := e.Read(); err != nil {
		e.Err = err
		return
	}

	if err := e.Cmd.Wait(); err != nil {
		e.Err = err
		return
	}
}

func (e *ExecObj) ErrorMsg() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *ExecObj) StatusCode() int {
	if e.ErrorMsg() == "" {
		return 200
	}
	return 400
}

func (e *ExecObj) StatusBody() []byte {
	if e.ErrorMsg() == "" {
		return e.ReadOut
	}
	return e.buildJsonError()
}

func (e *ExecObj) buildJsonError() (body []byte) {
	// TODO(mastensg): Don't do string matching, but rather something with this:
	// TODO(mastensg): https://golang.org/pkg/syscall/#WaitStatus.ExitStatus
	if e.ErrorMsg() != "exit status 1" {
		log.Println("stdout:", string(e.ReadOut))
		log.Println("stderr:", string(e.ReadErr))
		return body
	}

	javascriptError := json.RawMessage(e.ReadErr)

	synqError := struct {
		Name    string           `json:"name"`
		Message string           `json:"message"`
		Url     string           `json:"url"`
		Details *json.RawMessage `json:"details"`
	}{
		Name:    "javascript_query",
		Message: "An error occured while processing your JavaScript query.",
		Url:     "http://docs.synq.fm/api/v1/errors/javascript_query",
		Details: &javascriptError,
	}

	body, err := json.MarshalIndent(synqError, "", "    ")
	if err != nil {
		return body
	}
	return body
}
