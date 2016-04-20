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
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var pageName, template string

func addPage() error {
	// Need to understand our path
	setRelPathProjDir()

	// Read config.toml
	tomlData, err := ioutil.ReadFile(relPath + projectDir + string(filepath.Separator) + "config.toml")
	if err != nil {
		log.Fatal("Error config.toml could not be read")
	}

	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		log.Fatal("Error cannot parse config.toml")
	}

	domain := conf.Domain
	theme := conf.Theme

	themePath := relPath + domain + string(filepath.Separator) + "theme" + string(filepath.Separator) + theme
	sitePath := relPath + string(filepath.Separator) + domain + string(filepath.Separator) + "pages"
	fmt.Println(themePath)
	if !dirExist(themePath) {
		log.Fatal("Error cannot find theme to create markdown templates")
	}
	if !dirExist(themePath + string(filepath.Separator) + template + ".html") {
		log.Fatal("Error cannot find specified theme template")
	}

	// Only create the markdown file if it does not already exist
	if !dirExist(sitePath + string(filepath.Separator) + pageName + ".md") {
		var fileOutput string

		// Open template file & get contents
		temp, err := ioutil.ReadFile(themePath + string(filepath.Separator) + template + ".html")
		if err != nil {
			return err
		}

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
		fileOutput += "template = \"" + template + "\"\n"
		fileOutput += "\n+++\n\n"

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
		err = writeFile(sitePath+string(filepath.Separator)+pageName+".md", fileOutput)
		if err != nil {
			return err
		}
	}
	return nil
}

// pageCmd represents the page command
var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "Adds a new page to the website",
	Long: `Adds a new content page to the website based on the theme tempate specified.
	
	Uses --template flag to specify page template to build the markdown page from. Uses 'default' template if omitted.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pageName = strings.Join(args, " ")

		err := addPage()
		if err != nil {
			log.Fatal("Error unable to create new page " + pageName)
		}
	},
}

func init() {
	RootCmd.AddCommand(pageCmd)
	pageCmd.Flags().StringVarP(&template, "template", "", "default", "The template to use with new page")
}
