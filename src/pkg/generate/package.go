package generate

import (
	"github.com/defenseunicorns/uds-generator/src/config"
	"github.com/defenseunicorns/zarf/src/types"
)

func generateCommonPackage() (types.ZarfPackage, error) {
	type Generator struct {
		pkg       types.PackagerConfig
		component types.ZarfComponent
	}
	// Generate the upstream chart zarf yaml
	commonChart := types.ZarfChart{
		Name:      config.GenerateChartName,
		Namespace: config.GenerateChartName,
		URL:       config.GenerateChartUrl,
		Version:   config.GenerateChartVersion,
		GitPath:   config.GenerateChartGitPath,
	}
	metadata := types.ZarfMetadata{
		Name:        config.GenerateChartName + "-common",
		Description: "UDS" + SimpleTitleCase(config.GenerateChartName) + "Common Package",
		Version:     config.GenerateUDSVersion,
		URL:         config.GenerateChartUrl,
		Authors:     config.GenerateAuthors,
	}

	// Generate the config chart zarf yaml
	configChart := types.ZarfChart{
		Name:      "uds-" + config.GenerateChartName + "-config",
		Namespace: config.GenerateChartName,
		LocalPath: "chart",
		Version:   "0.1.0",
	}

	generator.component = types.ZarfComponent{
		Name:   config.GenerateChartName,
		Charts: []types.ZarfChart{configChart, commonChart},
	}
	components := []types.ZarfComponent{generator.component}
	// Generate the package
	packageInstance := types.ZarfPackage{
		Kind:       types.ZarfPackageConfig,
		Metadata:   metadata,
		Components: components,
	}
	return packageInstance, nil
}

func generateFlavoredPackage() (types.ZarfPackage, error) {
	flavor := "upstream"
	var required bool = true

	type Generator struct {
		pkg       types.PackagerConfig
		component types.ZarfComponent
	}
	// Generate the upstream chart zarf yaml
	flavoredChart := types.ZarfChart{
		Name:        config.GenerateChartName,
		ValuesFiles: []string{"../values/" + flavor + "-values.yaml"},
	}
	metadata := types.ZarfMetadata{
		Name:    config.GenerateChartName,
		Version: config.GenerateUDSVersion,
		URL:     config.GenerateChartUrl,
		Authors: config.GenerateAuthors,
	}

	generator.component = types.ZarfComponent{
		Name:     config.GenerateChartName,
		Required: &required,
		Import: types.ZarfComponentImport{
			Path: "common",
		},
		Charts: []types.ZarfChart{flavoredChart},
		Only: types.ZarfComponentOnlyTarget{
			Flavor: flavor,
		},
	}
	components := []types.ZarfComponent{generator.component}
	// Generate the package
	packageInstance := types.ZarfPackage{
		Kind:       types.ZarfPackageConfig,
		Metadata:   metadata,
		Components: components,
	}
	return packageInstance, nil
}
