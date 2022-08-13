package cmd

type Config struct {
	WorkingDirectory string `yaml:"-"`
	CfgFile          string `yaml:"-"`
	CfgFileName      string `yaml:"-"`
	CfgFileExt       string `yaml:"-"`
	UserHome         string `yaml:"-"`
	UsingLocalConfig bool   `yaml:"-"`
	*Repository
}
