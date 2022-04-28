package types

// GitHubRepositoryLoaderConfig is configuration for GitHubRepositoryLoader
type GitHubRepositoryLoaderConfig struct {
	DetectDefaultBranch bool   `json:"detect_default_branch" yaml:"detect_default_branch"`
	DefaultBranchName   string `json:"default_branch_name" yaml:"default_branch_name"`
}

// Configuration represents the configuration of the application.
type Configuration struct {
	DefaultRepository string                       `json:"default_repo" yaml:"default_repo"`
	GitHubLoader      GitHubRepositoryLoaderConfig `json:"github_loader" yaml:"github_loader"`
}
