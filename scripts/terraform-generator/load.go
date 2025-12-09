package tfgen

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"
)

func Load(ctx context.Context) (*spec.Specification, *spec.Specification, error) {
	resp, err := http.Get("https://openapi.logto.io/source.yaml")
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("expected 200 status code, got %s", resp.Status)
	}

	f, err := os.CreateTemp("", "openapi")
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(f.Name())

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	_, err = f.Write(data)
	if err != nil {
		return nil, nil, err
	}
	f.Close()

	cmd := exec.Command(
		"tfplugingen-openapi",
		"generate",
		"--config=./config/generator_config.yml",
		"--output=./provider_code_spec.json",
		f.Name(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}

	content, err := os.ReadFile("provider_code_spec.json")
	if err != nil {
		return nil, nil, err
	}
	specification, err := spec.Parse(ctx, content)
	if err != nil {
		return nil, nil, err
	}

	content, err = os.ReadFile("config/provider_code_extra.json")
	if err != nil {
		return nil, nil, err
	}
	extra, err := spec.Parse(ctx, content)
	if err != nil {
		return nil, nil, err
	}

	return &specification, &extra, nil
}
