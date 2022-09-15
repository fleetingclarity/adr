package cmd

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
)

type Config struct {
	WorkingDirectory string `yaml:"-"`
	CfgFile          string `yaml:"-"`
	CfgFileName      string `yaml:"-"`
	CfgFileExt       string `yaml:"-"`
	UserHome         string `yaml:"-"`
	UsingLocalConfig bool   `yaml:"-"`
	*Repository
	*ADRConfig
}

func (c *Config) EnsureRepositoryExists() error {
	if _, err := os.Stat(c.Repository.Path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path.Join(c.WorkingDirectory, c.Repository.Path), os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("warning: unable to create the repository directory %s. You will likely need to create it manually", path.Join(c.WorkingDirectory, c.Repository.Path)))
		}
	}
	return nil
}

func (c *Config) Write(w io.Writer) error {
	o, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(o)
	return err
}
