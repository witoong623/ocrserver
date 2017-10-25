package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/otiai10/gosseract"
)

type ocrContext struct {
	ocrClient *gosseract.Client
	tempPath  string
}

// OCRReturnResult is type holding return OCR result and will be parse to JSON string
type OCRReturnResult struct {
	Result string
	Error  string
}

func main() {
	listenAddr, ok := os.LookupEnv("SERVER_ADDR")
	if !ok {
		// In case local development
		listenAddr = ":8080"
	}
	tempPath, ok := os.LookupEnv("TEMP_DIR")
	if !ok {
		tempPath = ""
	}

	ocrClient, err := gosseract.NewClient()
	if err != nil {
		log.Fatal("Cannot instantiate OCR Client. ", err.Error())
	}
	ocrContext := &ocrContext{ocrClient: ocrClient, tempPath: tempPath}

	http.Handle("/gettext", handlePostImage(ocrContext))

	log.Println("Server running.")
	err = http.ListenAndServe(listenAddr, nil)
	log.Println("Server is stoped.")
	log.Fatal(err)
}

func handlePostImage(c *ocrContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		uploadedFile, _, err := r.FormFile("image")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		defer uploadedFile.Close()

		tempFile, err := ioutil.TempFile(c.tempPath, "ocr-request-")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		if _, err := io.Copy(tempFile, uploadedFile); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		result, err := c.ocrClient.Src(tempFile.Name()).Out()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OCRReturnResult{Result: result})
	})
}
