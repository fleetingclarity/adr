package cmd

// Repository contains the metadata for an ADR repository (e.g. configurations, path, etc)
type Repository struct {
	RelativePath string `yaml:"path"`
}
