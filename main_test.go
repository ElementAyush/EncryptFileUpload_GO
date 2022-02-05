package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-services/gomod/services"
	"github.com/stretchr/testify/assert"
)

func TestEndpointUploadErr(t *testing.T) {
	log.Println("Running test")

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("testdata/testfile1.txt")
	part1, errFile1 := writer.CreateFormFile("file", filepath.Base("testdata/testfile1.txt"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		log.Println(errFile1)
		return
	}
	defer file.Close()
	_ = writer.WriteField("userId", "ayudxt")
	_ = writer.WriteField("objectName", "testfile1.txt")
	err := writer.Close()
	if err != nil {
		log.Println(err)
		return
	}
	assert := assert.New(t)
	req, err := http.NewRequest("PUT", "/upload", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(services.Upload)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusInternalServerError)
	if rr.Code == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

func TestEndpointDownload(t *testing.T) {
	assert := assert.New(t)
	req, err := http.NewRequest("GET", "/download", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(services.Download)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusOK)
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}
