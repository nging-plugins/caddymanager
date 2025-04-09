package github

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	download "github.com/admpub/go-download/v2"
	downloadProgressbar "github.com/admpub/go-download/v2/progressbar"
	"github.com/admpub/log"
	"github.com/webx-top/restyclient"
)

type GithubRelease struct {
	TagName string               `json:"tag_name"`
	Assets  []GithubReleaseAsset `json:"assets"`
}

type GithubReleaseAsset struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	URL  string `json:"browser_download_url"`
}

func GetReleaseURL(orig, repo string) (*GithubRelease, error) {
	apiURL := `https://api.github.com/repos/` + orig + `/` + repo + `/releases/latest`
	client := restyclient.Retryable()
	result := &GithubRelease{}
	resp, err := client.SetResult(result).Get(apiURL)
	if err != nil {
		return result, err
	}
	if resp.IsError() {
		return result, fmt.Errorf("%s http error: %s", apiURL, resp.Status())
	}
	return result, nil
}

func Download(orig, repo string) (string, error) {
	r, err := GetReleaseURL(orig, repo)
	if err != nil {
		return ``, err
	}
	if len(r.Assets) == 0 {
		return ``, ErrNoAssetsFound
	}
	osName := runtime.GOOS
	osArch := runtime.GOARCH
	extension := `.tar.gz`
	suffix := fmt.Sprintf(`_%s_%s%s`, osName, osArch, extension)
	for _, asset := range r.Assets {
		if !strings.HasSuffix(asset.URL, suffix) {
			continue
		}
		var saveDir string
		saveDir, err = os.MkdirTemp(``, `nging-`+repo)
		if err != nil {
			return ``, fmt.Errorf("mkdir temp dir error: %w", err)
		}
		savePath := filepath.Join(saveDir, filepath.Base(asset.Name))
		var n int64
		log.Infof(`Downloading %s %s`, repo, asset.URL)
		opt := &download.Options{}
		prog := downloadProgressbar.New(opt)
		defer prog.Wait()
		n, err = download.Download(asset.URL, savePath, opt)
		if err != nil {
			return ``, fmt.Errorf("download error: %w", err)
		}
		if n != asset.Size {
			return ``, fmt.Errorf("download size mismatch: %d != %d", n, asset.Size)
		}
		return savePath, nil
	}
	return ``, ErrNoAssetsFound
}
