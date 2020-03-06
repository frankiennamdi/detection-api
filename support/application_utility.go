package support

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type TemporaryDir struct {
	tempPath string
}

func NewTemporaryDir(dir, prefix string) *TemporaryDir {
	dir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		log.Panicf("%+v", err)
	}

	return &TemporaryDir{tempPath: dir}
}

func (tempDir *TemporaryDir) Path() string {
	return tempDir.tempPath
}

func (tempDir *TemporaryDir) Clean() {
	if err := os.RemoveAll(tempDir.tempPath); err != nil {
		log.Panicf("%+v", err)
	}
}

func GetEnvironmentValue(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func Resolve(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	_, callerFile, _, _ := runtime.Caller(0)
	callerDir := filepath.Join(filepath.Dir(callerFile))
	rootDir := filepath.Dir(callerDir)

	return fmt.Sprintf("%s/%s", rootDir, path)
}
