// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023-Present The UDS Authors

// Package lang contains the language strings in english used by UDS
package lang

const (
	// uds-cli generate
	CmdGenerateShort = "Generates a UDS Package from a helm chart"
	CmdGenerateLong  = `Generates a zarf package integrated with UDS Core based on inputs for a helm chart. Supports OCI, Git, and Helm chart repositories.
Example:
  uds-generator generate -c https://charts.bitnami.com/bitnami/nginx-8.9.0.tgz -n nginx -v 8.9.0
`

	CmdGenerateFlagChart        = "URL for the helm chart"
	CmdGenerateFlagChartGitPath = "(optional) git path for the helm chart if git based"
	CmdGenerateFlagName         = "Name of the helm chart"
	CmdGenerateFlagVersion      = "Version for the helm chart"
	CmdGenerateOutputDir        = "Output directory for the generated package"

	// Server
	CmdServeShort      = `Serve API`
	CmdServeLong       = `Serve the UDS API and web interface on the specified port`
	CmdServeFlagPort   = `Port to serve the API on`
	CmdServeFlagwsPort = `Port to serve the websocket on`

	// scaffold
	CmdScaffoldShort = "Scaffolds a new UDS Package repository"
	CmdScaffoldLong  = `Scaffolds a new UDS Package repository with the given name, organization, author, and platform

Use "uds-generator [command] --help" for more information about a command.

Example: 
  uds-generator scaffold -n harbor -a "the UDS Authors" -o "defenseunicorns" -p github.com
`
	CmdScaffoldFlagName         = "Name of the package to scaffold"
	CmdScaffoldFlagOrganization = "Organization of the package"
	CmdScaffoldFlagAuthor       = "Author of the package"
	CmdScaffoldFlagPlatform     = "Platform for the package eg github.com, gitlab.com"
)
