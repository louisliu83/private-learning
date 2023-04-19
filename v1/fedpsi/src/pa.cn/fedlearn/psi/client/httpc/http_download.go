package httpc

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func DoGetBig(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := checkResponseCode(resp); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected EOF") && len(data) != 0 {
			// igrone this error
		} else {
			return nil, err
		}
	}
	return data, nil
}

func DownloadAndSaveFile(url string, headers map[string]string, targetFileName string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := checkResponseCode(resp); err != nil {
		return err
	}

	targetFile, err := os.OpenFile(targetFileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, resp.Body)
	return err
}
