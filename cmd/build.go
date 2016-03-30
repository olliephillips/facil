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
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/russross/blackfriday"
	"github.com/spf13/cobra"
)

type pageContent struct {
	Path    string
	Content string
}

type navigationContent struct {
	Text  string
	Order string
	Link  string
}

type config struct {
	Domain string
	Theme  string
}

type pageConfig struct {
	Meta       meta
	Navigation navigation
	Design     design
}

type meta struct {
	Title       string `toml:"title"`
	Description string `toml:"description"`
}

type navigation struct {
	Text  string
	Order string
}

type design struct {
	Template string
}

var project string
var relPath string
var projectDir string
var conf config
var pageConf pageConfig

var partialsOutput map[string]string

var pages []pageContent

var siteMap []string
var nav []string

//var navHTML string
type navigationItems []navigationContent

var navElements navigationItems

func (slice navigationItems) Len() int {
	return len(slice)
}

func (slice navigationItems) Less(i, j int) bool {
	first, _ := strconv.Atoi(slice[i].Order)
	second, _ := strconv.Atoi(slice[j].Order)
	return int(first) < int(second)
}

func (slice navigationItems) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func processPageFile(page string, dest string) {
	var output string

	// Get the template from file
	markdown, err := ioutil.ReadFile(page)
	if err != nil {
		log.Fatal("Error page could not be built")
	}

	// Process toml
	var markdownToml = regexp.MustCompile(`\+\+\+\n*([\d\D]*)\n*\+\+\+`)
	tomlSection := markdownToml.FindStringSubmatch(string(markdown))

	if len(tomlSection) > 1 {
		if _, err := toml.Decode(tomlSection[1], &pageConf); err != nil {
			//log.Fatal("Error cannot parse page's toml config")
			log.Fatal(err)
		}

		// Add to sitemap & nav
		element := strings.Split(dest, "compiled")[1]
		element = strings.Replace(element, string(filepath.Separator)+"index.html", "/", -1)

		// Add to sitemap
		sitemapElement := conf.Domain + element
		siteMap = append(siteMap, sitemapElement)

		// Add to nav
		navElement := strings.Replace(element, ".html", "", -1)

		nav := navigationContent{
			Text:  pageConf.Navigation.Text,
			Order: pageConf.Navigation.Order,
			Link:  navElement,
		}

		navElements = append(navElements, nav)

		// Read template from theme
		pageTemplate := relPath + projectDir + string(filepath.Separator) + "theme" + string(filepath.Separator) + conf.Theme + string(filepath.Separator) + pageConf.Design.Template + ".html"

		template, err := ioutil.ReadFile(pageTemplate)
		if err != nil {
			log.Fatal("Error template could not be read")
		}

		output = string(template)

		// Merge Meta
		output = processMeta(string(markdown), output)

		// Merge Elements
		output = processElements(string(markdown), output)

		// Merge Partials (we have a map of these)
		output = processPartials(output)

		p := pageContent{
			Path:    dest,
			Content: output,
		}

		// Add this page to our slice
		pages = append(pages, p)

	}
}

func processBlog() {

}

func processFile(source string, dest string, contentType string) (err error) {
	dest = strings.Replace(dest, ".md", ".html", -1)

	switch contentType {
	case "page":
		processPageFile(source, dest)
	case "blog":
		//processBlogFile(sourceFile)
	}
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
		sourcefilepointer := source + string(filepath.Separator) + obj.Name()
		destinationfilepointer := dest + string(filepath.Separator) + obj.Name()

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

func reflectField(pageConf *pageConfig, field string) string {
	r := reflect.ValueOf(pageConf.Meta)
	f := reflect.Indirect(r).FieldByName("Title")
	return string(f.String())
}

func processMeta(markdown string, template string) string {

	var output string

	// Parse meta tags in template with regex
	var templateToken = regexp.MustCompile(`\[\[meta\sname\=\"([a-zA-Z0-9]*)\"\s*]]`)
	templateTokens := templateToken.FindAllStringSubmatch(string(template), -1)

	// Range over all the template tokens
	for i := range templateTokens {
		token := strings.ToLower(templateTokens[i][1])
		upperFirstToken := upperFirst(token)
		tokenValue := reflectField(&pageConf, upperFirstToken)

		replace := "[[meta name=\"" + token + "\"]]"
		template = strings.Replace(template, replace, tokenValue, -1)
	}
	output = template

	// Return a merged string
	return output
}

func processPartials(template string) string {
	var output string

	// Parse partial tags in template with regex
	var templateToken = regexp.MustCompile(`\[\[partial\sname\=\"([a-zA-Z0-9]*)\"\s*]]`)
	templateTokens := templateToken.FindAllStringSubmatch(string(template), -1)

	// We have a map loaded with all processed partials
	// Range over template tokens and find replace with corresponding value from map
	for i := range templateTokens {
		token := templateTokens[i][1]
		if partialsOutput[token] != "" {
			// We have a processed partial stored, do find replace
			replace := "[[partial name=\"" + token + "\"]]"
			template = strings.Replace(template, replace, partialsOutput[token], -1)
		}
	}
	output = template

	// Return a merged string
	return output
}

func processElements(markdown string, template string) string {
	// THis is our regex \*\*\*\s([a-zA-Z0-9]*)\s.*\n([\d\D][^\*]*)\*\*\*  (needs g modifier) to pick out the name and markdown from the mark down files
	// Merge the files!!

	var output string

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
	output = template

	// Return a merged string
	return output
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
		partialsOutput[k] = processElements(partialsMarkdown[k], partialsTemplate[k])
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

func writePages(nav string) {
	// We have a slice struct of pages (package global) with all the info we need

	// We need to add our nav, we've parsed all the pages and built it, so this is first opportunity
	// the token to replace is [[navigation]]

	// We then need to write the pages to their correct location in the compiled directory
	for i := range pages {
		content := pages[i].Content
		dest := pages[i].Path

		// Do replacement of [[navigation]]
		content = strings.Replace(content, "[[navigation]]", nav, -1)

		// Write page
		err := ioutil.WriteFile(dest, []byte(content), 0755)
		if err != nil {
			log.Fatal("Error unable to write a static file")
		}
	}

}

func makeNav() string {
	html := "<ul>\n"
	// Sort elements
	sort.Sort(navElements)

	var curLevel, prevLevel int
	prevLevel = 2
	for i := range navElements {
		// Need to reorder based on order struct properties
		text := navElements[i].Text
		link := navElements[i].Link

		linkElements := strings.Split(link, "/")
		elementsCount := len(linkElements)

		if elementsCount == 3 {
			curLevel = 3
		} else {
			curLevel = 2
		}

		if curLevel == 2 && prevLevel == 3 {
			html += "\t\t</ul>\n\t</li>\n"
			html += "\t<li><a href=\"" + link + "\">" + text + "</a></li>\n"
		}

		if curLevel == 3 && prevLevel == 2 {
			html += "\t<li><a href=\"" + link + "\">" + text + "</a>\n"
			html += "\t\t<ul>\n"

		}
		if curLevel == 2 && prevLevel == 2 {
			html += "\t<li><a href=\"" + link + "\">" + text + "</a></li>\n"
		}
		if curLevel == 3 && prevLevel == 3 {
			html += "\t\t\t<li><a href=\"" + link + "\">" + text + "</a></li>\n"
		}
		prevLevel = curLevel
	}

	if curLevel == 3 && prevLevel == 3 {
		html += "\t\t</ul>\n\t</li>\n"
	}

	html += "</ul>"

	return html
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
	processDir(relPath+projectDir+string(filepath.Separator)+"pages", relPath+projectDir+string(filepath.Separator)+"compiled", "page")

	// Traverse projects blog folder
	// Each blog file
	//processDir()

	// Make Navigation
	nav := makeNav()

	// Write pages, replacing navigation token
	writePages(nav)

	// Write a sitemap.xml
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
