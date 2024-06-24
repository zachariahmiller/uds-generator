// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The UDS Authors

// Package cmd contains the CLI commands for UDS.
package cmd

import (
	"log"

	"github.com/defenseunicorns/uds-generator/src/config"
	"github.com/defenseunicorns/uds-generator/src/config/lang"
	"github.com/defenseunicorns/uds-generator/src/pkg/common"
	"github.com/defenseunicorns/uds-generator/src/pkg/generate"
	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:     "scaffold",
	Aliases: []string{"s"},
	Short:   lang.CmdScaffoldShort,
	Long:    lang.CmdScaffoldLong,
	Run: func(_ *cobra.Command, _ []string) {
		repoInfo := common.RepoInfo{
			Name:         config.ScaffoldRepoName,
			Organization: config.ScaffoldOrganization,
			Author:       config.ScaffoldAuthor,
			Platform:     config.ScaffoldPlatform,
			Version:      config.GenerateUDSVersion,
		}
		generate.Scaffold(repoInfo)
		log.Printf("Scaffolded %s", repoInfo.Name)
	},
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)
	scaffoldCmd.Flags().StringVarP(&config.ScaffoldRepoName, "name", "n", "example", lang.CmdScaffoldFlagName)
	scaffoldCmd.Flags().StringVarP(&config.ScaffoldOrganization, "organization", "o", "defenseunicorns", lang.CmdScaffoldFlagOrganization)
	scaffoldCmd.Flags().StringVarP(&config.ScaffoldAuthor, "author", "a", "The UDS Authors", lang.CmdScaffoldFlagName)
	scaffoldCmd.Flags().StringVarP(&config.ScaffoldPlatform, "platform", "p", "github.com", lang.CmdScaffoldFlagPlatform)
	scaffoldCmd.Flags().StringVarP(&config.GenerateOutputDir, "output", "d", "generated", lang.CmdGenerateOutputDir)
	scaffoldCmd.Flags().StringVarP(&config.GenerateChartVersion, "version", "v", "0.1.0", lang.CmdGenerateFlagVersion)
	scaffoldCmd.MarkFlagRequired("name")
	// scaffoldCmd.MarkFlagRequired("organization")
	// scaffoldCmd.MarkFlagRequired("platform")
	// scaffoldCmd.MarkFlagRequired("author")
}
