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

// This file contains most of the functionality behind 'facil start', it scaffolds a new site, creating
// the directory and file structure based on the chosen theme

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var domain, theme string

// Creates a config.toml file in the new sites root directory containing canonical domain and theme the site uses
func createConfigToml() error {
	var fileOutput string
	sitePath := basePath + string(filepath.Separator) + "sites" + string(filepath.Separator) + domain
	fileOutput += "domain = \"" + domain + "\"\n"
	fileOutput += "theme = \"" + theme + "\"\n"
	fileOutput += "https = \"off\" # Options are off, on\n"
	fileOutput += "pretty = \"off\" # Options are off, on\n"

	// Write file
	err := writeFile(sitePath+string(filepath.Separator)+"config.toml", fileOutput)
	if err != nil {
		return err
	}
	return nil
}

// Creates markdown file by parsing a template from the theme. All tokens are read from the template and setup
// in the created markdown file
func createMarkdownTemplate(template string, filename string, isPage bool) error {
	themePath := basePath + string(filepath.Separator) + "sites" + string(filepath.Separator) + domain + string(filepath.Separator) + "theme" + string(filepath.Separator) + theme
	sitePath := basePath + string(filepath.Separator) + "sites" + string(filepath.Separator) + domain

	if !dirExist(themePath) {
		log.Fatal("Error cannot find theme to create markdown templates")
	}
	if !dirExist(themePath + string(filepath.Separator) + template) {
		log.Fatal("Error cannot find specified theme template")
	}

	// Only create the markdown file if it does not already exist
	if !dirExist(sitePath + string(filepath.Separator) + filename) {
		var fileOutput string

		// Open template file & get contents
		temp, err := ioutil.ReadFile(themePath + string(filepath.Separator) + template)
		if err != nil {
			return err
		}

		// Is it a page or a partial template?
		if isPage {
			// Parse meta with regex
			var metaToken = regexp.MustCompile(`\[\[meta\sname\=\"([a-zA-Z0-9_-]*)\"]]`)
			metaTokens := metaToken.FindAllStringSubmatch(string(temp), -1)

			// Compose meta output

			fileOutput += "+++\n\n"
			fileOutput += "[Meta]\n"
			for i := range metaTokens {
				fileOutput += metaTokens[i][1] + " = \"\"\n"
			}

			// Add navigation tokens
			fileOutput += "\n[Navigation]\n"
			fileOutput += "text = \"\"\n"
			fileOutput += "order = \"99\"\n"

			// Add design tokens
			fileOutput += "\n[Design]\n"
			fileOutput += "template = \"" + strings.Replace(template, ".html", "", -1) + "\"\n"
			fileOutput += "\n+++\n\n"
		}

		// Parse element with regex
		var elementToken = regexp.MustCompile(`\[\[element\stype\=\"([a-zA-Z0-9]*)\"\sname\=\"([a-zA-Z0-9_-]*)\"\sdescription\=\"(.*)"]]`)
		elementTokens := elementToken.FindAllStringSubmatch(string(temp), -1)

		// Compose element output
		for i := range elementTokens {
			fileOutput += "***" + strings.ToUpper(elementTokens[i][1]) + "*** " + strings.Title(elementTokens[i][2]) + " (" + elementTokens[i][3] + ")\n\n"
			if elementTokens[i][1] == "html" {
				fileOutput += "# Your " + strings.Title(elementTokens[i][2]) + " markdown/html syntax here\n\n"
			} else {
				fileOutput += "Your " + strings.Title(elementTokens[i][2]) + " text syntax here\n\n"
			}

			fileOutput += "***\n\n\n\n\n"
		}

		// Write to file
		err = writeFile(sitePath+string(filepath.Separator)+filename, fileOutput)
		if err != nil {
			return err
		}
	}
	return nil
}

// Creates a new site in the sites directory
func scaffoldSite() error {
	// Partials and Blog created later if needed
	paths := []string{
		"sites" + string(filepath.Separator) + domain,
		"sites" + string(filepath.Separator) + domain + string(filepath.Separator) + "pages",
		"sites" + string(filepath.Separator) + domain + string(filepath.Separator) + "theme",
		"sites" + string(filepath.Separator) + domain + string(filepath.Separator) + "compiled",
	}

	// Create folder structure
	for i := range paths {
		// Check not exist
		if !dirExist(basePath + paths[i]) {
			err := os.Mkdir(basePath+paths[i], 0755)
			if err != nil {
				return err
			}
		}
	}

	// Copy selected theme into sites theme folder
	src := basePath + string(filepath.Separator) + "themes" + string(filepath.Separator) + theme
	dest := basePath + string(filepath.Separator) + "sites" + string(filepath.Separator) + domain + string(filepath.Separator) + "theme" + string(filepath.Separator) + theme

	// Check not exist before copy
	if !dirExist(dest) {
		err := copyDir(src, dest)
		if err != nil {
			return err
		}
	}

	// Create config.toml file - domain, theme, https, pretty
	err := createConfigToml()
	if err != nil {
		return err
	}

	// Add default markdown file and any partials - enough to render homepage
	err = createMarkdownTemplate("default.html", "pages/index.md", true)
	if err != nil {
		return err
	}

	// List partials, run createMarkdownTemplate() for each
	partialsDir := dest + string(filepath.Separator) + "partials"
	if dirExist(partialsDir) {
		// Create a site partials directory
		if !dirExist(basePath + "sites" + string(filepath.Separator) + domain + string(filepath.Separator) + "partials") {
			err := os.Mkdir(basePath+"sites"+string(filepath.Separator)+domain+string(filepath.Separator)+"partials", 0755)
			if err != nil {
				return err
			}
		}
		err := filepath.Walk(partialsDir, func(path string, f os.FileInfo, _ error) error {
			if !f.IsDir() {
				if strings.ToLower(strings.Split(f.Name(), ".")[1]) == "html" {
					tempFile := "partials" + string(filepath.Separator) + f.Name()
					mdFile := "partials" + string(filepath.Separator) + strings.Replace(f.Name(), ".html", "", -1) + ".md"

					err := createMarkdownTemplate(tempFile, mdFile, false)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Scaffolds a new website",
	Long: `Scaffolds a new website and creates markdown files based on theme chosen.
    
    Uses --theme flag to specify theme to setup site with, or the included 'default' theme
    `,
	Run: func(cmd *cobra.Command, args []string) {
		// Store domain we are scaffolding
		domain = strings.Join(args, " ")

		// Set a packpage level var for use in other functions
		setBasePath()

		// Create directories and file structure
		err := scaffoldSite()
		if err != nil {
			log.Fatal("Error unable to setup new site " + domain)
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&theme, "theme", "", "default", "The theme to use with new site")
}
