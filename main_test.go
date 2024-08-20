package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Helper functions:
// 1. Creating multipart from file upload request
func createFileUploadRequest(filename, endpoint, csvData string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, strings.NewReader(csvData))
	if err != nil {
		return nil, err
	}
	writer.Close()
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
func TestEchoHandler(t *testing.T) {
	csvData := "1,2,3\n4,5,6\n7,8,9\n"
	req, err := createFileUploadRequest("/echo", "matrix.csv", csvData)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(echoHandler)
	handler.ServeHTTP(r, req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v expected %v", status, http.StatusOK)
	}
	expected := csvData
	if r.Body.String() != expected {
		t.Errorf("unexpected result: got %v expected %v", r.Body.String(), expected)
	}
}
func TestTransposeHandler(t *testing.T) {
	csvData := "1,2,3\n4,5,6\n7,8,9\n"
	req, err := createFileUploadRequest("/invert", "matrix.csv", csvData)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(transposeHandler)
	handler.ServeHTTP(r, req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v expected %v", status, http.StatusOK)
	}
	expected := "1,4,7\n2,5,8\n3,6,9\n"
	if r.Body.String() != expected {
		t.Errorf("unexpected result: got %v expected %v", r.Body.String(), expected)
	}
}
func TestSumHandler(t *testing.T) {
	csvData := "1,2,3\n4,5,6\n7,8,9\n"
	req, err := createFileUploadRequest("/echo", "matrix.csv", csvData)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(sumHandler)
	handler.ServeHTTP(r, req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "sum:45"
	if strings.TrimSpace(r.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", r.Body.String(), expected)
	}
}
func TestMultiplyHandler(t *testing.T) {
	csvData := "1,2,3\n4,5,6\n7,8,9\n"
	req, err := createFileUploadRequest("/echo", "matrix.csv", csvData)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(multiplyHandler)
	handler.ServeHTTP(r, req)

	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "mul:362880"
	if strings.TrimSpace(r.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", r.Body.String(), expected)
	}
}
func TestMlattenHandler(t *testing.T) {
	csvData := "1,2,3\n4,5,6\n7,8,9\n"
	req, err := createFileUploadRequest("/echo", "matrix.csv", csvData)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(flattenHandler)
	handler.ServeHTTP(r, req)

	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "1,2,3,4,5,6,7,8,9"
	if strings.TrimSpace(r.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", r.Body.String(), expected)
	}
}
