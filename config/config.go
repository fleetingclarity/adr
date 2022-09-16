package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
)

const (
	DefaultRepositoryDir = "docs/decisions"
	DefaultConfigExt     = "yaml"
	DefaultConfigName    = ".adr"
)

type Config struct {
	WorkingDirectory string `yaml:"-"`
	CfgFileName      string `yaml:"-"`
	CfgFileExt       string `yaml:"-"`
	*Repository
	*ADR
}

// EnsureRepositoryExists creates the repository directory if it doesn't exist. ADRs will be stored in this directory
func (c *Config) EnsureRepositoryExists() error {
	if _, err := os.Stat(c.Repository.Path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path.Join(c.WorkingDirectory, c.Repository.Path), os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("warning: unable to create the repository directory %s. You will likely need to create it manually", path.Join(c.WorkingDirectory, c.Repository.Path)))
		}
	}
	return nil
}

// CreateAndWrite creates a file from the current structs values and then calls Write with it
func (c *Config) CreateAndWrite() error {
	f, err := os.Create(path.Join(c.WorkingDirectory, c.CfgFileName+"."+c.CfgFileExt))
	defer f.Close()
	if err != nil {
		return err
	}
	err = c.Write(f)
	if err != nil {
		return err
	}
	return nil
}

// Write writes the current struct to the given io.Writer
func (c *Config) Write(w io.Writer) error {
	o, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(o)
	return err
}

// Managed checks if the current working directory is already managed by this tool
func (c *Config) Managed() (bool, error) {
	cfgFile := path.Join(c.WorkingDirectory, c.CfgFileName+"."+c.CfgFileExt)
	if _, err := os.Stat(cfgFile); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func NewDefaultConfig() *Config {
	wd, _ := os.Getwd()
	return &Config{
		WorkingDirectory: wd,
		CfgFileName:      DefaultConfigName,
		CfgFileExt:       DefaultConfigExt,
		Repository: &Repository{
			Path: DefaultRepositoryDir,
		},
		ADR: &ADR{
			FormatName:    "Nygard",
			Sections:      []string{"title", "date", "status", "context", "decision", "consequences"},
			TitleTemplate: defaultTitleTemplate,
			BodyTemplate:  defaultBodyTemplate,
		},
	}
}
