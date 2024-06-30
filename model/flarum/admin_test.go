package flarum

import (
	"log"
	"os"
	"testing"

	"github.com/corvofeng/go-flarum/util"
)

func TestGetExtensionMetadata(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	logger := util.GetLogger()
	logger.Error(path)
	_, _ = ReadExtensionMetadata("../../view/extensions")
}
