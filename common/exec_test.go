package common

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupPid() string {
	pidFile := "test.pid"
	os.Remove(pidFile)
	return pidFile
}

func TestCheckPid(t *testing.T) {
	pidFile := setupPid()
	assert := require.New(t)
	// no such file should be fine
	pid, err := CheckPid(pidFile)
	assert.Nil(err)
	assert.Equal(os.Getpid(), pid)
	// same pid shoudl be allowed
	pid2, err := CheckPid(pidFile)
	assert.Nil(err)
	assert.Equal(pid2, pid)
	// running pid is not ok
	ioutil.WriteFile(pidFile, []byte("0"), 0755)
	_, err = CheckPid(pidFile)
	assert.NotNil(err)
	assert.Equal("Pid '0' already exists, will not run", err.Error())
}
