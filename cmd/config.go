package cmd

type Config struct {
	WorkingDirectory string `yaml:"-"`
	CfgFile          string `yaml:"-"`
	*Repository
}
