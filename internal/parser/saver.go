package parser

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const dir = "reports"

func Save(data CarData) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}
	}

	path := filepath.Join(dir, data.Report.CarPlate)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", path)
		}
	}

	if data.Image != "" {
		file := filepath.Join(path, "preview.jpg")
		var out *os.File
		_, err := os.Stat(file)
		if !os.IsNotExist(err) {
			if err := os.Remove(file); err != nil {
				return errors.Wrap(err, "failed to remove preview.jpg")
			}
		}
		out, err = os.Create(file)
		if err != nil {
			return errors.Wrap(err, "failed to create preview.jpg")
		}
		defer out.Close()

		resp, _ := http.Get(data.Image)
		defer resp.Body.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			return errors.Wrap(err, "failed to copy image")
		}
	}

	j, err := json.MarshalIndent(data.Report, "", " ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}

	if err := os.WriteFile(filepath.Join(path, "info.json"), j, 0644); err != nil {
		return errors.Wrap(err, "failed to save json")
	}

	return nil
}
