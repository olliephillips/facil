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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/russross/blackfriday"
	"github.com/spf13/cobra"
)

var project string
var relPath string
var projectDir string
var conf config

var partialsOutput map[string]string

var siteMap []string
var nav []string

type config struct {
	Domain string
	Theme  string
}

func processPageFile(page string) {
	// Get the template from file

	// Read template from theme

	// Merge markdown file and template file into new output

	// Write new file to complied folder
}

func processBlog() {

}

func addToSitemapAndNav(dest string) {
	element := strings.Split(dest, "compiled")[1]
	element = strings.Replace(element, string(filepath.Separator)+"index.html", "", -1)

	// Add to sitemap
	sitemapElement := conf.Domain + element
	siteMap = append(siteMap, sitemapElement)

	// Add to nav
	navElement := strings.Replace(element, ".html", "", -1)
	nav = append(nav, navElement)
}

func processFile(source string, dest string, contentType string) (err error) {
	dest = strings.Replace(dest, ".md", ".html", -1)

	// Add to sitemap and nav
	addToSitemapAndNav(dest)

	switch contentType {
	case "pages":
		processPageFile(source)
	case "blog":
		//processBlog(sourceFile)
	}

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	// At this point we have an empty file we can write to

	/*
		_, err = io.Copy(destfile, sourcefile)
		if err == nil {
			sourceinfo, err := os.Stat(source)
			if err != nil {
				err = os.Chmod(dest, sourceinfo.Mode())
			}
		}*/
	return
}

func processDir(source string, dest string, contentType string) (err error) {
	// Dest properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// Create sub-directories - recursively
			err = processDir(sourcefilepointer, destinationfilepointer, contentType)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Perform copy
			err = processFile(sourcefilepointer, destinationfilepointer, contentType)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}

func buildPages() {

	// Get markdown pages iterate subdirectories to understand nav hierarchy

	processDir(relPath+projectDir+string(filepath.Separator)+"pages", relPath+projectDir+string(filepath.Separator)+"compiled", "page")

	// Build a special nav/menu partial

	// Read the toml config of each markdown file

	// Read all the markdown tokens to map

	// Get the template to use

	// Process any partials

	// Process token replacements

	// Write a page in the format directory_path/page_title/index.html
}

func processPartial(filename string, markdown string, template string) string {
	// THis is our regex \*\*\*\s([a-zA-Z0-9]*)\s.*\n([\d\D][^\*]*)\*\*\*  (needs g modifier) to pick out the name and markdown from the mark down files
	// Merge the files!!

	var partialOutput string

	// Parse element tags in template with regex
	var templateToken = regexp.MustCompile(`\[\[element\sname\=\"([a-zA-Z0-9]*)\"\sdescription\=\"(.*)"]]`)
	templateTokens := templateToken.FindAllStringSubmatch(string(template), -1)

	// Parse element tags in markdown file with regex
	var markdownToken = regexp.MustCompile(`\*\*\*\s([a-zA-Z0-9]*)\s.*\n([\d\D][^\*]*)\*\*\*`)
	markdownTokens := markdownToken.FindAllStringSubmatch(string(markdown), -1)

	// Range over all the markdown tokens
	for i := range markdownTokens {
		token := strings.ToLower(markdownTokens[i][1])

		// Process Markdown content ready for inclusion
		htmlContent := string(blackfriday.MarkdownCommon([]byte(markdownTokens[i][2])))

		if token == strings.ToLower(templateTokens[i][1]) {
			// We have a match
			description := templateTokens[i][2]
			// What to replace
			replace := "[[element name=\"" + token + "\" description=\"" + description + "\"]]"
			// Replace
			template = strings.Replace(template, replace, htmlContent, -1)
		}
	}
	partialOutput = template

	// Return a merged string
	return partialOutput
}

func buildPartials() {
	// Get partials
	partialsMarkdown := make(map[string]string)
	partialsTemplate := make(map[string]string)
	partialsOutput = make(map[string]string) // Package namespace

	err := filepath.Walk(relPath+projectDir+string(filepath.Separator)+"partials", func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			//go func() {
			filename := strings.ToLower(strings.Split(f.Name(), ".")[0])
			extension := strings.ToLower(strings.Split(f.Name(), ".")[1])
			if extension == "md" {
				// Get the markdown file
				md, err := ioutil.ReadFile(relPath + projectDir + string(filepath.Separator) + "partials" + string(filepath.Separator) + f.Name())
				if err != nil {
					log.Fatal("Error reading a partial markdown file")
				}

				tmp, err := ioutil.ReadFile(relPath + projectDir + string(filepath.Separator) + "theme" + string(filepath.Separator) + theme +
					string(filepath.Separator) + "partials" + string(filepath.Separator) + filename + ".html")

				if err != nil {
					log.Fatal("Error reading a partial template file")
				}

				// Store to our maps
				partialsMarkdown[filename] = string(md)
				partialsTemplate[filename] = string(tmp)
			}
			//}()
		}
		return nil
	})

	if err != nil {
		log.Fatal("Error unable to build partials")
	}

	// Range over one of the maps, pass name and both the markdown and template as args to processPartial()
	for k := range partialsMarkdown {
		partialsOutput[k] = processPartial(k, partialsMarkdown[k], partialsTemplate[k])
	}
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

	// Wipe entire "compiled" directory
	deleteDirectoryContents(relPath + projectDir + string(filepath.Separator) + "compiled")

	// Copy theme assets to compiled folder, remove html templates
	copyThemeAssets()

	// Build partials
	buildPartials()
	//fmt.Println(partialsOutput)

	// Build Pages
	buildPages()

	// Traverse projects blog folder
	// Each blog file
	//go buildBlogs()

	// Should have a map/slice of pages to include navigation on
	// Do this for each page and write it to compiled directory
	// WHat about nav order override?
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
