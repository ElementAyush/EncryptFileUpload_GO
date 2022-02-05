package services

import (
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func TestFileEncryption(t *testing.T) {
	// uploadEncryptedFile := func(w http.ResponseWriter, file []byte, filename string, userId string, fileSize int64, unfile []byte) bool {
	// 	return true
	// }
	theServiceMock := serviceMock{}
	theServiceMock.On("encryptFileWithKey", mock.Anything).Return(true)

	var w http.ResponseWriter
	var file multipart.File
	isTrue := encryptFileWithKey(w, file, "testfile.txt", "ayudxt", 65)
	assert := assert.New(t)
	assert.Equal(isTrue, true)
}
