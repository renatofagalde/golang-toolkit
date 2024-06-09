package toolkit

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

const PATH_TEST_PNG = "./testdata/senna.png"

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)
	require.Len(t, s, 10, "Wrong length for RandomString")
}

var uploadTests = []struct {
	name          string
	allowedType   []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedType: []string{"image/jpeg", "image/png"}, renameFile: false, errorExpected: false},
	{name: "allowed rename", allowedType: []string{"image/jpeg", "image/png"}, renameFile: true, errorExpected: false},
	{name: "not allowed", allowedType: []string{"image/jpeg", "image/gif"}, renameFile: false, errorExpected: true},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, test := range uploadTests {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			part, err := writer.CreateFormFile("file", PATH_TEST_PNG)
			if err != nil {
				t.Error(err)
			}
			f, err := os.Open(PATH_TEST_PNG)
			if err != nil {
				t.Error(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error(err)
			}
			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}

		}()

		//read from the pipe
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFilesTypes = test.allowedType

		uploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploads/", test.renameFile)
		if err != nil && !test.errorExpected {
			require.NoError(t, err)
		}

		if test.errorExpected {
			fmt.Println("error ", err)
			require.Error(t, err)
		} else {
			require.FileExists(t, fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
		}

		wg.Wait()
	}
}

func TestTools_UploadOneFile(t *testing.T) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := writer.CreateFormFile("file", PATH_TEST_PNG)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(PATH_TEST_PNG)
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			t.Error(err)
		}
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}

	}()

	//read from the pipe
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var testTools Tools

	uploadedFiles, err := testTools.UploadOneFile(request, "./testdata/uploads/", true)
	if err != nil {
		t.Error(err)
	}

	require.NoError(t, err)
	require.FileExists(t, fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName))

	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName))

}
