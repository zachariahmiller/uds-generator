package common

type PackageInfo struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Version string `json:"version"`
	GitPath string `json:"gitPath"`
}

type RepoInfo struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Platform     string `json:"platform"`
	Author       string `json:"author"`
}
