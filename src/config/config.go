package config

import (
	"time"

	"github.com/defenseunicorns/uds-cli/src/types"
)

var (
	// CommonOptions tracks user-defined values that apply across commands.
	CommonOptions types.BundleCommonOptions

	// CLIVersion track the version of the CLI
	CLIVersion = "unset"

	// CLIArch is the computer architecture of the device executing the CLI commands
	CLIArch string

	// SkipLogFile is a flag to skip logging to a file
	SkipLogFile bool

	// ListTasks is a flag to print available tasks in a TaskFileLocation
	ListTasks bool

	// TaskFileLocation is the location of the tasks file to run
	TaskFileLocation string

	// SetRunnerVariables is a map of the run time variables defined using --set
	SetRunnerVariables map[string]string

	// GenerateChartUrl is a URL for the helm chart to generate a UDS Package based off of
	GenerateChartUrl string

	// GenerateChartName is the name of the helm chart to generate a UDS Package based off of
	GenerateChartName string

	// GenerateChartVersion is the version of the helm chart to generate a UDS Package based off of
	GenerateChartVersion string

	// GenerateChartGitPath is the git path of the helm chart to generate a UDS Package based off of
	GenerateChartGitPath string

	// GenerateOutputDir is the directory to output the generated package to
	GenerateOutputDir string

	// Name of repo/package/application for scaffolding used in repoPath like github.com/organization/uds-package-repo
	ScaffoldRepoName string

	// Organization of package for scaffolding used in repoPath eg github.com/organization/repo
	ScaffoldOrganization string

	// Author of package for scaffolding
	ScaffoldAuthor string
	// Git Platform for package for scaffolding eg gitlab.com, github.com
	ScaffoldPlatform string

	GenerateUDSVersion = GenerateChartVersion + "-uds.0"

	GenerateAuthors = "The UDS Authors"
	// HelmTimeout is the default timeout for helm deploys
	HelmTimeout = 15 * time.Minute

	ApiPort string

	WebsocketPort string
)
