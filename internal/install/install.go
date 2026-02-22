package install

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func TargetDir(root, version string) string {
	return filepath.Join(root, version)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func PromptYesNo(msg string) (bool, error) {
	fmt.Printf("%s [y/N]: ", msg)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes", nil
}

func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

func TempFile() (string, error) {
	f, err := os.CreateTemp("", "goversion-*.tar.gz")
	if err != nil {
		return "", err
	}
	name := f.Name()
	f.Close()
	return name, nil
}

func Download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func ExtractTarGz(src, destRoot string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Strip the leading "go/" from the tar path
		name := strings.TrimPrefix(hdr.Name, "go/")
		if name == hdr.Name { // wasn't a "go/" prefix
			continue // skip non-go files (like README)
		}

		targetPath := filepath.Join(destRoot, name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		default:
			// skip other types for now
		}
	}

	return nil
}

func VerifyChecksum(filePath, expected string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	calculated := hex.EncodeToString(h.Sum(nil))
	if calculated != expected {
		return fmt.Errorf("checksum mismatch: got %s, want %s", calculated, expected)
	}
	return nil
}
