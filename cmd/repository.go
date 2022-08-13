package cmd

// Repository contains the metadata for an ADR repository (e.g. configurations, path, etc)
type Repository struct {
	Path string `yaml:"path"`
}
