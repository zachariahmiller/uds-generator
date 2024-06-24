// handler.go

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/defenseunicorns/uds-generator/src/config"
	"github.com/defenseunicorns/uds-generator/src/pkg/common"
	"github.com/defenseunicorns/uds-generator/src/pkg/generate"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/defenseunicorns/zarf/src/pkg/packager"
	"github.com/defenseunicorns/zarf/src/pkg/utils"
	"github.com/defenseunicorns/zarf/src/types"
	"github.com/goccy/go-yaml"
)

// TODO RENAME THIS STRUCT
type ResponseData struct {
	Images map[string][]string `yaml:"images"`
}

// APIHandler handles API requests
func APIHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the API!"))
}

// findImagesHandler handles the find-images API endpoint
func FindImagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var pkgInfo common.PackageInfo
	if err := json.NewDecoder(r.Body).Decode(&pkgInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Here, you should implement the logic to process the pkgInfo and generate the imgMap.
	// For demonstration, let's assume this function:
	genPath, err := processPackageInfo(pkgInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// defer os.RemoveAll(genPath)
	imagesData, found := loadImagesData(genPath)
	if !found {
		http.Error(w, "Images data not found", http.StatusNotFound)
		return
	}

	response := ResponseData{
		Images: imagesData,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// findImagesHandler handles the find-images API endpoint
func ScaffoldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var repoInfo common.RepoInfo
	if err := json.NewDecoder(r.Body).Decode(&repoInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Version is equal to %s", repoInfo.Version)
	config.GenerateChartVersion = repoInfo.Version
	log.Printf("Version is equal to %s", config.GenerateChartVersion)
	repoInfo.Version = config.GenerateUDSVersion
	log.Printf("Version is equal to %s", repoInfo.Version)
	response, err := generate.Scaffold(repoInfo)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// this is using the zarf generate now which is kind of meh... but it works
// when updated probably should return something different then a string
func processPackageInfo(pkgInfo common.PackageInfo) (string, error) {

	pkgConfig := types.PackagerConfig{}
	// zarf madness
	pkgConfig.GenerateOpts.Name = pkgInfo.Name

	// gotta move some of these things to common config (like dirs and stuff)
	pkgConfig.CreateOpts.BaseDir = "generated"                                                  // This is the base directory where the package will be created
	pkgConfig.CreateOpts.Output = "generator-tmp/output" + pkgInfo.Name + "-" + pkgInfo.Version // This is the output directory where the package will be created
	pkgConfig.FindImagesOpts.RepoHelmChartPath = pkgConfig.GenerateOpts.GitPath
	pkgConfig.FindImagesOpts.KubeVersionOverride = "v1.27.0"
	pkgConfig.FindImagesOpts.Why = ""
	pkgConfig.GenerateOpts.Version = pkgInfo.Version
	pkgConfig.GenerateOpts.GitPath = pkgInfo.GitPath
	pkgConfig.GenerateOpts.URL = pkgInfo.URL
	pkgConfig.GenerateOpts.Output = pkgConfig.CreateOpts.Output

	pkgClient := packager.NewOrDie(&pkgConfig)
	defer pkgClient.ClearTempPaths()

	if err := pkgClient.Generate(); err != nil {
		message.Fatalf(err, err.Error())
	}
	return pkgConfig.CreateOpts.Output + "/zarf.yaml", nil
}

func templateHack(pkgInfo common.PackageInfo) (string, error) {

	pkgConfig := types.PackagerConfig{}
	// zarf madness
	pkgConfig.GenerateOpts.Name = pkgInfo.Name

	// gotta move some of these things to common config (like dirs and stuff)
	pkgConfig.CreateOpts.BaseDir = "generated"                                                  // This is the base directory where the package will be created
	pkgConfig.CreateOpts.Output = "generator-tmp/output" + pkgInfo.Name + "-" + pkgInfo.Version // This is the output directory where the package will be created
	pkgConfig.FindImagesOpts.RepoHelmChartPath = pkgConfig.GenerateOpts.GitPath
	pkgConfig.FindImagesOpts.KubeVersionOverride = "v1.27.0"
	pkgConfig.FindImagesOpts.Why = ""
	pkgConfig.GenerateOpts.Version = pkgInfo.Version
	pkgConfig.GenerateOpts.GitPath = pkgInfo.GitPath
	pkgConfig.GenerateOpts.URL = pkgInfo.URL
	pkgConfig.GenerateOpts.Output = pkgConfig.CreateOpts.Output

	templates := os.Stdout
	os.Stdout = nil
	pkgClient := packager.NewOrDie(&pkgConfig)
	defer pkgClient.ClearTempPaths()
	if images, err := pkgClient.FindImages(); err != nil {
		message.Fatalf(err, err.Error())
	} else {
		utils.ColorPrintYAML(templates, nil, false)
		utils.ColorPrintYAML(images, nil, false)
	}
	return pkgConfig.CreateOpts.Output + "/zarf.yaml", nil
}

// loadImagesData loads the images data from the YAML file

func loadImagesData(filePath string) (map[string][]string, bool) {
	log.Printf("Loading YAML file from path: %s", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading YAML file: %v", err)
		return nil, false
	}

	var config map[string]interface{}
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Printf("Error unmarshaling YAML: %v", err)
		return nil, false
	}
	log.Printf("Parsed YAML content: %v", config)
	imagesData, found := findImagesKey(config)
	return imagesData, found
}

func findImagesKey(data map[string]interface{}) (map[string][]string, bool) {
	for key, value := range data {
		log.Printf("Processing key: %s with value type: %T", key, value)

		if key == "images" {
			if images, ok := value.([]interface{}); ok {
				result := make(map[string][]string)
				for _, img := range images {
					if imgStr, ok := img.(string); ok {
						parts := strings.Split(imgStr, ":")
						if len(parts) >= 2 {
							result[parts[0]] = append(result[parts[0]], parts[1])
						} else {
							result[imgStr] = []string{}
						}
					}
				}
				log.Println("Images processed successfully.")
				return result, true
			}
			log.Println("Found 'images', but could not convert to []interface{}")
			return nil, false
		}

		// Handle map structures recursively
		if subMap, ok := value.(map[string]interface{}); ok {
			log.Printf("Descending into subMap under key: %s", key)
			if result, found := findImagesKey(subMap); found {
				return result, true
			}
		} else if subList, ok := value.([]interface{}); ok { // Handle slices within the structure
			log.Printf("Descending into list under key: %s", key)
			for _, item := range subList {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if result, found := findImagesKey(itemMap); found {
						return result, true
					}
				}
			}
		}
	}
	log.Println("No 'images' key found at this level.")
	return nil, false
}
