package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"os/exec"
)

type Spec struct {
	Version     string      `json:"version"`
	Provider    any         `json:"provider"`
	Resources   []*Resource `json:"resources,omitempty"`
	Datasources []*Resource `json:"datasources,omitempty"`
}

type Resource struct {
	Name   string `json:"name"`
	Schema Schema `json:"schema"`
}

type Schema struct {
	Attributes []any `json:"attributes"`
}

func Main() error {
	resp, err := http.Get("https://openapi.logto.io/source.yaml")
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("expected 200 status code, got %s", resp.Status)
	}

	f, err := os.CreateTemp("", "openapi")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	f.Write(data)
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
		return err
	}

	spec, err := readSpec("provider_code_spec.json")
	if err != nil {
		return err
	}
	extra, err := readSpec("config/provider_code_extra.json")
	if err != nil {
		return err
	}

	resources := map[string]*Resource{}
	for _, r := range extra.Resources {
		resources[r.Name] = r
	}
	datasources := map[string]*Resource{}
	for _, r := range extra.Datasources {
		datasources[r.Name] = r
	}

	for _, resource := range spec.Resources {
		extraResource, found := resources[resource.Name]
		if found {
			resource.Schema.Attributes = append(
				resource.Schema.Attributes,
				extraResource.Schema.Attributes...,
			)
		}
	}
	for _, datasource := range spec.Datasources {
		extraDatasource, found := datasources[datasource.Name]
		if found {
			datasource.Schema.Attributes = append(
				datasource.Schema.Attributes,
				extraDatasource.Schema.Attributes...,
			)
		}
	}

	content, err := json.MarshalIndent(spec, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile("provider_code_spec.json", content, 0o644)
	if err != nil {
		return err
	}

	cmd = exec.Command(
		"tfplugingen-framework",
		"generate",
		"all",
		"--input=./provider_code_spec.json",
		"--output=internal/provider",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func readSpec(path string) (*Spec, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var extra Spec
	if err := json.Unmarshal(content, &extra); err != nil {
		return nil, err
	}
	return &extra, nil
}

func main() {
	err := Main()
	if err != nil {
		log.Fatalf("err: %v", err)
		os.Exit(1)
	}
}
