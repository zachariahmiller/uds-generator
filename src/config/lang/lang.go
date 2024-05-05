// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023-Present The UDS Authors

// Package lang contains the language strings in english used by UDS
package lang

const (
	// uds-cli generate
	CmdGenerateShort            = "Generates a UDS Package from a helm chart"
	CmdGenerateLong             = "Generates a zarf package integrated with UDS Core based on inputs for a helm chart."
	CmdGenerateFlagChart        = "URL for the helm chart"
	CmdGenerateFlagChartGitPath = "(optional) git path for the helm chart if git based"
	CmdGenerateFlagName         = "Name of the helm chart"
	CmdGenerateFlagVersion      = "Version for the helm chart"
	CmdGenerateOutputDir        = "Output directory for the generated package"
	CmdServeShort               = `Serve API`
	CmdServeLong                = `Serve the UDS API and web interface on the specified port`
	CmdServeFlagPort            = `Port to serve the API on`
	CmdServeFlagwsPort          = `Port to serve the websocket on`
)
