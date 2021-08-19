package flarum

import (
	"log"
	"os"
	"testing"
	"zoe/util"
)

func TestGetExtensionMetadata(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	t.Error(path)
	logger := util.GetLogger()
	logger.Error(path)
	_, _ = ReadExtensionMetadata("../../view/extensions")
}
