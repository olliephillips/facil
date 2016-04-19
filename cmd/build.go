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
	"time"

	"github.com/BurntSushi/toml"
	"github.com/joeguo/sitemap"
	"github.com/russross/blackfriday"
	"github.com/spf13/cobra"
)

type (
	config struct {
		Domain string
		Theme  string
		Https  string
		Pretty string
	}

	pageConfig struct {
		Meta       meta
		Navigation navigation
		Design     design
	}

	pageContent struct {
		Path    string
		Content string
	}

	// TOML parsing structs
	meta struct {
		Title       string `toml:"title"`
		Description string `toml:"description"`
		Keywords    string `toml:"keywords"`
		Author      string `toml:"author"`
	}

	navigation struct {
		Text  string
		Order string
	}

	design struct {
		Template string
	}

	// Navigation building
	navigationContent struct {
		Text        string
		Order       string
		Link        string
		NaturalLink string
	}

	navigationItems []navigationContent
)

var (
	project        string
	relPath        string
	projectDir     string
	conf           config
	pageConf       pageConfig
	partialsOutput map[string]string
	pages          []pageContent
	siteMap        []string
	nav            []string
	navElements    navigationItems
)

// Implement sort interface on navigationItems
func (slice navigationItems) Len() int {
	return len(slice)
}

func (slice navigationItems) Less(i, j int) bool {
	first, _ := strconv.Atoi(slice[i].Order)
	second, _ := strconv.Atoi(slice[j].Order)
	return first < second
}

func (slice navigationItems) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func processPageFile(page string, dest string) {
	var output, writeEl, navEl, navNatEl, sitemapEl, el string

	// Get the template from file
	markdown, err := ioutil.ReadFile(page)
	if err != nil {
		log.Fatal("Error page could not be built")
	}

	// Process toml
	var markdownToml = regexp.MustCompile(`\+\+\+\n*([\d\D]*)\n*\+\+\+`)
	tomlSection := markdownToml.FindStringSubmatch(string(markdown))

	if len(tomlSection) > 1 {
		// Read TOML
		if _, err := toml.Decode(tomlSection[1], &pageConf); err != nil {
			log.Fatal(err)
		}

		// Read template from theme
		pageTemplate := relPath + projectDir + string(filepath.Separator) + "theme" + string(filepath.Separator) + conf.Theme + string(filepath.Separator) + pageConf.Design.Template + ".html"

		template, err := ioutil.ReadFile(pageTemplate)
		if err != nil {
			log.Fatal("Error template could not be read")
		}

		// Initialise into output var
		output = string(template)

		// Merge Meta
		output = processMeta(string(markdown), output)

		// Merge Elements
		output = processElements(string(markdown), output)

		// Merge Partials (we have a map of these)
		output = processPartials(output)

		// Logic is going to branch here, depending on whether pretty URLs are in use

		switch conf.Pretty {
		case "on":
			// On, directory for page name and file is index.html
			splitDest := strings.Replace(strings.Split(dest, "compiled")[1], string(filepath.Separator), "", 1)
			splitDest2 := strings.Split(splitDest, string(filepath.Separator))

			if len(splitDest2) > 1 {
				el = splitDest2[1]
			} else {
				el = splitDest2[0]
			}

			writeEl = strings.Replace(dest, el, "", -1)
			var newDir string
			if el != "index.md" {
				newDir = strings.Replace(el, ".md", "", -1)
				os.MkdirAll(writeEl+string(filepath.Separator)+newDir, 0755)
			}
			if newDir != "" {
				writeEl = writeEl + newDir + string(filepath.Separator) + "index.html"
			} else {
				writeEl = writeEl + "index.html"
			}

			navEl = strings.Replace(strings.Split(writeEl, "compiled")[1], "index.html", "", -1)
			navNatEl = strings.Replace(strings.Split(strings.Replace(dest, ".md", ".html", -1), "compiled")[1], "index.html", "", -1)
			sitemapEl = conf.Domain + navEl

		default:
			// We should have conf.Pretty="off" but set as default
			writeEl = strings.Replace(dest, ".md", ".html", -1)
			navEl = strings.Replace(strings.Split(writeEl, "compiled")[1], "index.html", "", -1)
			navNatEl = strings.Replace(strings.Split(writeEl, "compiled")[1], "index.html", "", -1)
			sitemapEl = conf.Domain + navEl
		}

		// Add to sitemap
		siteMap = append(siteMap, sitemapEl)

		// Add to nav
		nav := navigationContent{
			Text:        pageConf.Navigation.Text,
			Order:       pageConf.Navigation.Order,
			Link:        navEl,
			NaturalLink: navNatEl,
		}
		navElements = append(navElements, nav)

		// Add tp pages slice
		p := pageContent{
			Path:    writeEl,
			Content: output,
		}
		pages = append(pages, p)
	}
}

func processFile(source string, dest string, contentType string) (err error) {

	switch contentType {
	case "page":
		processPageFile(source, dest)
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
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

func processMeta(markdown string, template string) string {

	var output string

	// Parse meta tags in template with regex
	var templateToken = regexp.MustCompile(`\[\[meta\sname\=\"([a-zA-Z0-9_-]*)\"\s*]]`)
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
	var templateToken = regexp.MustCompile(`\[\[partial\sname\=\"([a-zA-Z0-9_-]*)\"\s*]]`)
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
	var templateToken = regexp.MustCompile(`\[\[element\stype\=\"([a-zA-Z0-9]*)\"\sname\=\"([a-zA-Z0-9_-]*)\"\sdescription\=\"(.*)"]]`)
	templateTokens := templateToken.FindAllStringSubmatch(string(template), -1)

	// Parse element tags in markdown file with regex
	var markdownToken = regexp.MustCompile(`\*\*\*([a-zA-Z0-9]*)\*\*\*\s([a-zA-Z0-9_-]*)\s.*\n([\d\D][^\*]*)\*\*\*`)
	markdownTokens := markdownToken.FindAllStringSubmatch(string(markdown), -1)

	// Range over all the markdown tokens
	for i := range markdownTokens {
		ttype := strings.ToLower(markdownTokens[i][1])
		token := strings.ToLower(markdownTokens[i][2])
		tokenContent := markdownTokens[i][3]
		var htmlContent string

		// Process Markdown content ready for inclusion
		if ttype == "text" {
			// This should be output in raw form and not processed by markdown conversion
			htmlContent = string(tokenContent)
		} else {
			// Process the markdown
			htmlContent = string(blackfriday.MarkdownCommon([]byte(tokenContent)))
		}
		htmlContent = string(strings.Trim(htmlContent, "\n\t "))
		if token == strings.ToLower(templateTokens[i][2]) {
			// We have a match
			description := templateTokens[i][3]
			// What to replace
			replace := "[[element type=\"" + ttype + "\" name=\"" + token + "\" description=\"" + description + "\"]]"
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

	partialsPath := relPath + projectDir + string(filepath.Separator) + "partials"
	if dirExist(partialsPath) {
		err := filepath.Walk(partialsPath, func(path string, f os.FileInfo, _ error) error {
			if !f.IsDir() {
				//go func() {
				filename := strings.ToLower(strings.Split(f.Name(), ".")[0])
				extension := strings.ToLower(strings.Split(f.Name(), ".")[1])
				if extension == "md" {
					// Get the markdown file
					md, err := ioutil.ReadFile(partialsPath + string(filepath.Separator) + f.Name())
					if err != nil {
						log.Fatal("Error reading a partial markdown file")
					}

					tmp, err := ioutil.ReadFile(relPath + projectDir + string(filepath.Separator) + "theme" + string(filepath.Separator) + conf.Theme +
						string(filepath.Separator) + "partials" + string(filepath.Separator) + filename + ".html")

					if err != nil {
						log.Fatal("Error reading a partial template file")
					}

					// Store to our maps
					partialsMarkdown[filename] = string(md)
					//partialsTemplate[filename] = string(strings.Trim(tmp, "\t\n "))
					partialsTemplate[filename] = strings.Trim(string(tmp), "\t\n ")
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
}

func copyThemeAssets() {
	err := copyDir(relPath+projectDir+string(filepath.Separator)+"theme"+string(filepath.Separator)+conf.Theme, relPath+projectDir+string(filepath.Separator)+"compiled")
	if err != nil {
		log.Fatal("Error could not build theme assets")
	}

	// Remove partials
	partialsPath := relPath + projectDir + string(filepath.Separator) + "compiled" + string(filepath.Separator) + "partials"
	if dirExist(partialsPath) {

		err = filepath.Walk(partialsPath, func(path string, f os.FileInfo, _ error) error {
			// Remove html templates
			if !f.IsDir() {
				if strings.ToLower(strings.Split(f.Name(), ".")[1]) == "html" {
					_ = os.Remove(partialsPath + string(filepath.Separator) + f.Name())
				}
			}
			return nil
		})
	}

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

	prevLevel = 1
	for i := range navElements {
		// Need to reorder based on order struct properties
		text := navElements[i].Text
		link := strings.Replace(navElements[i].Link, string(filepath.Separator), "/", -1)
		naturalLink := strings.Replace(navElements[i].NaturalLink, string(filepath.Separator), "/", -1)

		linkElements := strings.Split(naturalLink, "/")
		elementsCount := len(linkElements)

		if elementsCount == 3 {
			curLevel = 2
		} else {
			curLevel = 1
		}

		if curLevel == 1 && prevLevel == 2 {
			html += "\t\t</ul>\n\t</li>\n"
			html += "\t<li><a href=\"" + link + "\">" + text + "</a></li>\n"
		}

		if curLevel == 2 && prevLevel == 1 {

			html += "\t<li><a href=\"" + link + "\">" + text + "</a>\n"
			html += "\t\t<ul>\n"

		}
		if curLevel == 1 && prevLevel == 1 {
			html += "\t<li><a href=\"" + link + "\">" + text + "</a></li>\n"
		}
		if curLevel == 2 && prevLevel == 2 {
			html += "\t\t\t<li><a href=\"" + link + "\">" + text + "</a></li>\n"
		}
		prevLevel = curLevel
	}

	if curLevel == 2 && prevLevel == 2 {
		html += "\t\t</ul>\n\t</li>\n"
	}

	html += "</ul>"

	return html
}

func setRelPathProjDir() {
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
}

func buildProject() {
	// Establish target directory based on project, check to ensure 'config.toml' exists
	setRelPathProjDir()

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

	// Build Pages
	processDir(relPath+projectDir+string(filepath.Separator)+"pages", relPath+projectDir+string(filepath.Separator)+"compiled", "page")

	// Make Navigation
	nav := makeNav()

	fmt.Println(nav)

	// Write pages, replacing navigation token
	writePages(nav)

	// Write a sitemap.xml.gz
	err = createSitemap()
	if err != nil {
		log.Fatal("Error sitemap.xml could not be written")
	}

}

func createSitemap() error {
	// Create a sitemap
	compiledFolder := relPath + projectDir + string(filepath.Separator) + "compiled"
	var siteMapElements []*sitemap.Item

	// Create our sitemap from our pages map
	for i := range siteMap {
		element := new(sitemap.Item)
		filename := siteMap[i]

		// HTTP or HTTPS?
		prefix := "http://"
		if conf.Https == "on" {
			prefix = "https://"
		}
		// Replace filepath separator
		element.Loc = strings.Replace(prefix+filename, string(filepath.Separator), "/", -1)
		element.LastMod = time.Now()
		element.Changefreq = "weekly"
		if pages[i].Path == compiledFolder+string(filepath.Separator)+"index.html" {
			element.Priority = 0.8
		} else {
			element.Priority = 0.3
		}
		siteMapElements = append(siteMapElements, element)
	}

	// Write sitemap.xml.gz
	err := sitemap.SiteMap(compiledFolder+string(filepath.Separator)+"sitemap.xml.gz", siteMapElements)
	if err != nil {
		return err
	}
	return nil
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
