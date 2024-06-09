package toolkit

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
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

func TestTools_CreateDirIfNotExist(t *testing.T) {
	var testTool Tools

	err := testTool.CreateDirIfNotExist("./testdata/myDir")
	if err != nil {
		t.Error(err)
	}

	err = testTool.CreateDirIfNotExist("./testdata/myDir")
	if err != nil {
		t.Error(err)
	}

	_ = os.Remove("./testdata/myDir")
}

func TestTools_CreateDirIfNotExistInvalidDirectory(t *testing.T) {
	var testTool Tools

	// we should not be able to create a directory at the root level (no permissions)
	err := testTool.CreateDirIfNotExist("/mydir")
	if err == nil {
		t.Error(errors.New("able to create a directory where we should not be able to"))
	}
}
func TestTools_DownloadLargeStaticFile(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	var testTools Tools

	testTools.DownloadStaticFile(rr, req, "./testdata", "senna.png", "w11.jpeg")

	res := rr.Result()
	defer res.Body.Close()

	if res.Header["Content-Length"][0] != "8829" {
		t.Error("wrong content length of", res.Header["Content-Length"][0])
	}

	if res.Header["Content-Disposition"][0] != "attachment; filename=\"w11.jpeg\"" {
		t.Error("wrong content disposition of", res.Header["Content-Disposition"][0])
	}

	_, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
}
func TestTools_Slugify(t *testing.T) {
	var testTool Tools

	for _, e := range slugTests {
		slug, err := testTool.Slugify(e.s)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: error received when none expected: %s", e.name, err.Error())
		}

		if !e.errorExpected && slug != e.expected {
			t.Errorf("%s: wrong slug returned; expected %s but got %s", e.name, e.expected, slug)
		}
	}
}

var slugTests = []struct {
	name          string
	s             string
	expected      string
	errorExpected bool
}{
	{name: "valid string", s: "now is the time", expected: "now-is-the-time", errorExpected: false},
	{name: "empty string", s: "", expected: "", errorExpected: true},
	{name: "complex string", s: "Now is the time for all GOOD men! + Fish & such &^?123", expected: "now-is-the-time-for-all-good-men-fish-such-123", errorExpected: false},
	{name: "japanese string", s: "こんにちは世界", expected: "", errorExpected: true},
	{name: "japanese string plus roman characters", s: "こんにちは世界 hello world", expected: "hello-world", errorExpected: false},
}
