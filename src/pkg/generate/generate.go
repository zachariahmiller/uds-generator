package generate

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/pkg/helpers"
	"github.com/defenseunicorns/uds-generator/src/config"
	"github.com/defenseunicorns/uds-generator/src/pkg/packager"
	"github.com/defenseunicorns/zarf/src/pkg/utils"
	"github.com/defenseunicorns/zarf/src/types"
	goyaml "github.com/goccy/go-yaml"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

//go:embed chart/*
var chart embed.FS

//go:embed github/*
var github embed.FS

//go:embed tasks.yaml
var tasks embed.FS

var kubeVersionOverride = "1.28.0"

type Generator struct {
	pkg       types.PackagerConfig
	component types.ZarfComponent
}

var generator Generator
var required bool = true

// Generate a UDS Package from a given helm chart in the config
func Generate() (types.ZarfPackage, error) {
	log.Println("Starting Generate function")
	// Generate the metadata
	fmt.Printf("Generating UDS Package for %s\n", config.GenerateChartName)
	metadata := types.ZarfMetadata{
		Name:    config.GenerateChartName,
		Version: config.GenerateUDSVersion,
		URL:     config.GenerateChartUrl,
		Authors: config.GenerateAuthors,
	}

	// Generate the config chart zarf yaml
	configChart := types.ZarfChart{
		Name:      "uds-config",
		Namespace: config.GenerateChartName,
		LocalPath: "chart",
		Version:   "0.1.0",
	}

	// Generate the upstream chart zarf yaml
	upstreamChart := types.ZarfChart{
		Name:      config.GenerateChartName,
		Namespace: config.GenerateChartName,
		URL:       config.GenerateChartUrl,
		Version:   config.GenerateChartVersion,
		GitPath:   config.GenerateChartGitPath,
	}

	// Generate the component
	generator.component = types.ZarfComponent{
		Name:     config.GenerateChartName,
		Required: &required,
		Charts:   []types.ZarfChart{configChart, upstreamChart},
		Only: types.ZarfComponentOnlyTarget{
			Flavor: "upstream",
		},
	}
	components := []types.ZarfComponent{generator.component}

	// Generate the package
	packageInstance := types.ZarfPackage{
		Kind:       types.ZarfPackageConfig,
		Metadata:   metadata,
		Components: components,
	}
	log.Println("Created package instance")

	// Create generated directory if it doesn't exist
	if err := os.MkdirAll(config.GenerateOutputDir, 0755); err != nil {
		log.Println("Failed to create output directory:", err)
		return packageInstance, err
	}
	zarfPath := filepath.Join(config.GenerateOutputDir, "zarf.yaml")

	// Write in progress zarf yaml to a file
	text, _ := goyaml.Marshal(packageInstance)
	if err := os.WriteFile(zarfPath, text, 0644); err != nil {
		log.Println("Failed to write zarf.yaml:", err)
		return packageInstance, err
	}
	log.Println("Wrote zarf.yaml to output directory")

	// Copy template chart to destination
	writeEmbeddedFolder(chart, "", "")
	// Copy github folder to destination
	writeEmbeddedFolder(github, "", "")
	// Manipulate chart
	log.Println("Entering manipulatePackage")
	if err := manipulatePackage(); err != nil {
		log.Println("Error in manipulatePackage:", err)
		return packageInstance, err
	}

	if err := writeTasks(tasks); err != nil {
		log.Println("Error in writeTasks:", err)
		return packageInstance, err
	}

	// Find images to add to the component
	generator.pkg = types.PackagerConfig{
		CreateOpts: types.ZarfCreateOptions{
			Flavor:  "upstream",
			BaseDir: config.GenerateOutputDir,
		},
		// TODO: Why is this needed?
		FindImagesOpts: types.ZarfFindImagesOptions{
			KubeVersionOverride: kubeVersionOverride,
		},
	}
	log.Println("Set up packager configuration")

	packager := packager.NewOrDie(&generator.pkg)
	defer packager.ClearTempPaths()

	//stdout := os.Stdout
	//os.Stdout = nil
	log.Println("Finding images...")
	images, err := packager.FindImages()
	if err != nil {
		log.Println("Error finding images:", err)
		return packageInstance, err
	}
	//os.Stdout = stdout
	log.Println("Found images")

	log.Println("Finding manifest images...")
	manifestImages, manifestErr := packager.GetManifests()
	if manifestErr != nil {
		log.Println("Error finding images:", manifestErr)
		return packageInstance, manifestErr
	}
	//os.Stdout = stdout
	log.Println("Found manifestImages", manifestImages)

	// TODO: Strip off cosign signatures/attestations?
	components[0].Images = images[config.GenerateChartName]

	utils.ColorPrintYAML(packageInstance, nil, false)

	// Write final zarf yaml to a file
	text, _ = goyaml.Marshal(packageInstance)
	if err := os.WriteFile(zarfPath, text, 0644); err != nil {
		log.Println("Failed to write final zarf.yaml:", err)
		return packageInstance, err
	}

	log.Println("Generate function completed successfully")
	return packageInstance, nil
}

func manipulatePackage() error {
	log.Println("Starting manipulatePackage function")
	var udsPackage Package
	packagePath := filepath.Join(config.GenerateOutputDir, "chart", "templates", "uds-package.yaml")
	packageYaml, err := os.ReadFile(packagePath)
	if err != nil {
		log.Println("Failed to read uds-package.yaml:", err)
		return err
	}
	if err := goyaml.Unmarshal(packageYaml, &udsPackage); err != nil {
		log.Println("Failed to unmarshal uds-package.yaml:", err)
		return err
	}
	log.Println("Unmarshalled uds-package.yaml")

	udsPackage.ObjectMeta.Name = config.GenerateChartName
	udsPackage.ObjectMeta.Namespace = config.GenerateChartName

	expose, err := findHttpServices()
	if err != nil {
		log.Println("Error in findHttpServices:", err)
	} else {
		fmt.Printf("Found services to expose: %v", expose)
		if expose != nil {
			udsPackage.Spec.Network.Expose = expose
		}
	}

	text, _ := goyaml.Marshal(udsPackage)
	if err := os.WriteFile(packagePath, text, 0644); err != nil {
		log.Println("Failed to write uds-package.yaml:", err)
		return err
	}
	log.Println("Wrote uds-package.yaml")

	log.Println("manipulatePackage function completed successfully")
	return nil
}

func writeTasks(tasks embed.FS) error {
	log.Println("Starting writeTasks function")
	fileName := "tasks.yaml"

	// Open the embedded file.
	fileData, err := tasks.Open(fileName)
	if err != nil {
		log.Println("Failed to open tasks.yaml:", err)
		return err
	}
	defer fileData.Close()

	// Create a new file in the target directory.
	targetPath := config.GenerateOutputDir + "/" + fileName
	outFile, err := os.Create(targetPath)
	if err != nil {
		log.Println("Failed to create tasks.yaml:", err)
		return err
	}
	defer outFile.Close()

	// Copy the content of the embedded file to the new file.
	if _, err := io.Copy(outFile, fileData); err != nil {
		log.Println("Failed to copy tasks.yaml:", err)
		return err
	}

	fmt.Printf("File %s copied to %s successfully.", fileName, config.GenerateOutputDir)
	log.Println("writeTasks function completed successfully")
	return nil
}

func findHttpServices() ([]Expose, error) {
	log.Println("Starting findHttpServices function")
	var exposeList []Expose
	chartName := config.GenerateChartName
	chartVersion := config.GenerateChartVersion
	repoURL := config.GenerateChartUrl

	settings := cli.New()
	actionConfig := new(action.Configuration)
	actionConfig.Init(settings.RESTClientGetter(), chartName, "", log.Printf)

	pull := action.NewPull()
	pull.Settings = cli.New()

	chartDownloader := downloader.ChartDownloader{
		Out:            nil,
		RegistryClient: nil,
		Verify:         downloader.VerifyNever,
		Getters:        getter.All(pull.Settings),
		Options: []getter.Option{
			getter.WithInsecureSkipVerifyTLS(config.CommonOptions.Insecure),
		},
	}
	log.Println("Set up chart downloader")

	temp := filepath.Join(config.GenerateOutputDir, "temp")
	if err := helpers.CreateDirectory(temp, 0700); err != nil {
		log.Println("Failed to create temporary directory:", err)
		return nil, err
	}
	defer os.RemoveAll(temp)
	log.Println("Created temporary directory")

	chartURL, _ := repo.FindChartInAuthRepoURL(repoURL, "", "", chartName, chartVersion, pull.CertFile, pull.KeyFile, pull.CaFile, getter.All(pull.Settings))
	log.Println("Created temporary directory")
	saved, _, err := chartDownloader.DownloadTo(chartURL, pull.Version, temp)
	if err != nil {
		log.Println("Failed to download chart:", err)
		return nil, err
	}
	log.Println("Downloaded chart")

	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.Replace = true // Skip the name check.
	client.ClientOnly = true
	client.IncludeCRDs = true
	client.Verify = false
	client.KubeVersion, _ = chartutil.ParseKubeVersion(kubeVersionOverride)
	client.InsecureSkipTLSverify = config.CommonOptions.Insecure
	client.ReleaseName = chartName
	client.Namespace = chartName

	loadedChart, err := loader.Load(saved)
	if err != nil {
		log.Println("Failed to load chart:", err)
		return nil, err
	}
	log.Println("Loaded chart")

	templatedChart, err := client.Run(loadedChart, nil)
	if err != nil {
		log.Println("Failed to run Helm client:", err)
		return nil, err
	}
	template := templatedChart.Manifest
	yamls, _ := utils.SplitYAML([]byte(template))
	var resources []*unstructured.Unstructured
	resources = append(resources, yamls...)

	for _, resource := range resources {
		if resource.GetKind() == "Service" {
			contents := resource.UnstructuredContent()
			var service v1.Service
			runtime.DefaultUnstructuredConverter.FromUnstructured(contents, &service)
			for _, port := range service.Spec.Ports {
				// Guess that we want to expose any ports named "http"
				if port.Name == "http" {
					expose := Expose{
						Gateway:  "tenant",
						Host:     service.ObjectMeta.Name,
						Port:     int(port.Port),
						Selector: service.Spec.Selector,
						Service:  service.ObjectMeta.Name,
						// TODO: Target Port
					}
					exposeList = append(exposeList, expose)
				}
			}
		}
	}
	log.Println("findHttpServices function completed successfully")
	return exposeList, nil
}
