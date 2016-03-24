// Copyright © 2016 Ollie Phillips <ollie@interject.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var project string
var relPath string
var projectDir string
var conf config

type config struct {
	Domain string
	Theme  string
}

func processPageFile() {
	// Get the template from file

	// Read template from theme

	// Merge markdown file and template file into new output

	// Write new file to complied folder
}

func processBlogFile() {

}

func processPartialFile() {

}

func copyThemeAssets() {
	err := copyDir(relPath+projectDir+string(filepath.Separator)+"theme"+string(filepath.Separator)+conf.Theme, relPath+projectDir+string(filepath.Separator)+"compiled")
	if err != nil {
		log.Fatal("Error could not build theme assets")
	}

	// Remove partials
	err = filepath.Walk(relPath+projectDir+string(filepath.Separator)+"compiled"+string(filepath.Separator)+"partials", func(path string, f os.FileInfo, _ error) error {
		// Remove html templates
		if !f.IsDir() {
			if strings.ToLower(strings.Split(f.Name(), ".")[1]) == "html" {
				_ = os.Remove(relPath + projectDir + string(filepath.Separator) + "compiled" + string(filepath.Separator) + "partials" + string(filepath.Separator) + f.Name())
			}
		}
		return nil
	})

	// Remove
	err = filepath.Walk(relPath+projectDir+string(filepath.Separator)+"compiled", func(path string, f os.FileInfo, _ error) error {
		// Remove html templates
		if !f.IsDir() {
			if strings.ToLower(strings.Split(f.Name(), ".")[1]) == "html" {
				_ = os.Remove(relPath + projectDir + string(filepath.Separator) + "compiled" + string(filepath.Separator) + f.Name())
			}
		} else if f.Name() == "partials" {
			// Remove partials folder
			_ = os.Remove(relPath + projectDir + string(filepath.Separator) + "compiled" + string(filepath.Separator) + f.Name())
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error could not build theme assets")
	}
}

func buildProject() {
	// Establish target directory based on project, check to ensure 'config.toml' exists

	// Get current directory look for toml
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error could not establish project directory")
	}

	dirs := strings.Split(dir, string(filepath.Separator))
	lastDir := dirs[len(dirs)-1]

	// In Facil directory?
	if dirExist("sites") && project != "" {
		projectDir = project
		relPath = "." + string(filepath.Separator) + "sites" + string(filepath.Separator)
	}

	// In Sites directory?
	if lastDir == "sites" && project != "" {
		projectDir = project
		relPath = "." + string(filepath.Separator)
	}

	// In site subdirectory?
	if lastDir == "pages" || lastDir == "blog" || lastDir == "partials" {
		// In subdirectory, for site folder hierarchy (we hope)
		if project == "" {
			projectDir = dirs[len(dirs)-2]
		} else {
			projectDir = project
		}
		relPath = ".." + string(filepath.Separator) + ".." + string(filepath.Separator)
	}

	if relPath == "" {
		// If relPath still unset then must assume in root of site folder
		if project == "" {
			projectDir = dirs[len(dirs)-1]
		} else {
			projectDir = project
		}
		relPath = ".." + string(filepath.Separator)
	}

	// OK so we think we have a relative path and a project, test that siteDir and config.toml exist
	if !dirExist(relPath + projectDir) {
		log.Fatal("Error project directory does not exist")
	}
	if !dirExist(relPath + projectDir + string(filepath.Separator) + "config.toml") {
		log.Fatal("Error project has no config.toml")
	}

	// Read toml config file to establish root domain and theme in use
	tomlData, err := ioutil.ReadFile(relPath + projectDir + string(filepath.Separator) + "config.toml")
	if err != nil {
		log.Fatal("Error config.toml could not be read")
	}

	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		log.Fatal("Error cannot parse config.toml")
	}

	// Copy theme assets to compiled folder, remove html templates
	copyThemeAssets()

	// Traverse projects partials folder
	// Parse each partial file

	// Traverse projects pages folder
	// Parse each page file

	// Traverse projects blog folder
	// Each blog file

}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds static website",
	Long: `Builds static website, combining markdown files with theme templates.
    
    Once built the website is available in the 'compiled' directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		project = strings.Join(args, " ")
		buildProject()
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
