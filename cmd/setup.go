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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Creates 'themes' and 'sites' directories in current working directory
func createFolders() error {
	err := os.Mkdir(basePath+"themes", 0755)
	err = os.Mkdir(basePath+"sites", 0755)
	if err != nil {
		return err
	}
	return nil
}

// Adds a basic 'default' theme to themes directory
func addDefaultTheme() error {
	// Folder paths to create
	paths := []string{
		"themes" + string(filepath.Separator) + "default",
		"themes" + string(filepath.Separator) + "default" + string(filepath.Separator) + "partials",
		"themes" + string(filepath.Separator) + "default" + string(filepath.Separator) + "js",
		"themes" + string(filepath.Separator) + "default" + string(filepath.Separator) + "css",
		"themes" + string(filepath.Separator) + "default" + string(filepath.Separator) + "images",
	}

	// Templates
	const defaultHTML = `
<html>
    <head>
        <title>[[meta name="title"]]</title>
        <meta name="description" content="[[meta name="description"]]">
        <meta name="keywords" content="[[meta name="keywords"]]"
        <meta name="author" content="[[meta name="author"]]"
    </head>
    <body>
        <div id="nav">
            [[navigation]]
        </div>
        <div id="body">
            <div id="title">
              <p>
                	[[element type="text" name="title" description="Set the title"]]
			  </p>
            </div>
            <div id="intro">
                [[element type="html" name="introduction" description="Add an introductory paragraph"]]
            </div>
            <div id="footer">
                [[partial name="footer"]]
            </div>
        </div>
    </body>
</html>
    `

	const defaultJS = `
some js
    `
	const defaultCSS = `
some css
    `

	const partialFooter = `
<div>[[element type="html" name="footer" description="A reuseable footer element"]]</div>
    `
	// Create folder structure
	for i := range paths {
		err := os.Mkdir(basePath+paths[i], 0755)
		if err != nil {
			return err
		}
	}

	// Write 'default' theme files
	defaultThemePath := basePath + "themes" + string(filepath.Separator) + "default" + string(filepath.Separator)
	err := writeFile(defaultThemePath+"default.html", defaultHTML)
	err = writeFile(defaultThemePath+"js"+string(filepath.Separator)+"facil.js", defaultJS)
	err = writeFile(defaultThemePath+"css"+string(filepath.Separator)+"facil.css", defaultCSS)
	err = writeFile(defaultThemePath+"partials"+string(filepath.Separator)+"footer.html", partialFooter)

	if err != nil {
		return err
	}

	return nil
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Facil for first use",
	Long:  `Creates skeleton 'sites' and 'themes' directories, including a single default theme`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check not already exist
		if !dirExist(basePath+"themes") && !dirExist(basePath+"sites") {
			// Create required folders
			err := createFolders()
			if err != nil {
				log.Fatal("Error could not create 'sites' or 'themes' directories")
			}
		}

		// Check default theme not already exist
		if !dirExist(basePath + "themes" + string(filepath.Separator) + "default") {
			// Add default theme to themes folder
			err := addDefaultTheme()
			if err != nil {
				log.Fatal("Error could not create default theme in themes directory")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}
