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
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var domain, theme string

func createConfigToml() {
	var fileOutput string
	sitePath := basePath + string(filepath.Separator) + "sites" + string(filepath.Separator) + domain
	fileOutput += "domain = \"" + domain + "\"\n"
	fileOutput += "theme = \"" + theme + "\"\n"

	// Write file
	writeFile(sitePath+string(filepath.Separator)+"config.toml", fileOutput)
}

func createMarkdownTemplate(template string, filename string, isPage bool) {
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
			log.Fatal("Could not open template file")
		}

		// Is it a page or a partial template?
		if isPage {
			// Parse meta with regex
			var metaToken = regexp.MustCompile(`\[\[meta\sname\=\"([a-zA-Z]*)\"]]`)
			metaTokens := metaToken.FindAllStringSubmatch(string(temp), -1)

			// Compose meta output

			fileOutput += "+++\n\n"
			fileOutput += "[Meta]\n"
			for i := range metaTokens {
				fileOutput += metaTokens[i][1] + " = \"\"\n"
			}

			// Add design tokens
			fileOutput += "\n[Design]\n"
			fileOutput += "template = \"" + strings.Replace(template, ".html", "", -1) + "\"\n"
			fileOutput += "\n+++\n\n"
		}

		// Parse element with regex
		var elementToken = regexp.MustCompile(`\[\[element\sname\=\"([a-zA-Z]*)\"\sdescription\=\"(.*)"]]`)
		elementTokens := elementToken.FindAllStringSubmatch(string(temp), -1)

		// Compose element output
		for i := range elementTokens {
			fileOutput += "*** " + strings.Title(elementTokens[i][1]) + " (" + elementTokens[i][2] + ") ***\n"
			fileOutput += "#Your markdown syntax here\n\n"
		}

		// Write to file
		writeFile(sitePath+string(filepath.Separator)+filename, fileOutput)
	}
}

func scaffoldSite() {
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
				log.Fatal("Error creating site directories")
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
			log.Fatal("Selected theme could not be copied to site")
		}
	}

	// Create config.toml file - domain, theme
	createConfigToml()

	// Add default markdown file and any partials - enough to render homepage
	createMarkdownTemplate("default.html", "pages/index.md", true)

	// List partials, run createMarkdownTemplate() for each
	partialsDir := dest + string(filepath.Separator) + "partials"
	if dirExist(partialsDir) {
		// Create a site partials directory
		if !dirExist(basePath + "sites" + string(filepath.Separator) + domain + string(filepath.Separator) + "partials") {
			err := os.Mkdir(basePath+"sites"+string(filepath.Separator)+domain+string(filepath.Separator)+"partials", 0755)
			if err != nil {
				log.Fatal("Error creating 'partials' directory")
			}
		}
		err := filepath.Walk(partialsDir, func(path string, f os.FileInfo, _ error) error {
			if !f.IsDir() {
				if strings.ToLower(strings.Split(f.Name(), ".")[1]) == "html" {
					tempFile := "partials" + string(filepath.Separator) + f.Name()
					mdFile := "partials" + string(filepath.Separator) + strings.Replace(f.Name(), ".html", "", -1) + ".md"
					createMarkdownTemplate(tempFile, mdFile, false)
				}
			}
			return nil
		})

		if err != nil {
			log.Fatal("Could not parse partials templates")
		}
	}

	// Are there blog templates, if so handle those?
	blogTemp := dest + string(filepath.Separator) + "blog.html"
	blogPostTemp := dest + string(filepath.Separator) + "blog-post.html"
	if dirExist(blogTemp) && dirExist(blogPostTemp) {
		// Need to create a 'blog' directory
		blogDir := basePath + "sites" + string(filepath.Separator) + domain + string(filepath.Separator) + "blog"
		if !dirExist(blogDir) {
			err := os.Mkdir(blogDir, 0755)
			if err != nil {
				log.Fatal("Error creating 'blog' directory")
			}
		}

		// With blog.html we create a blog.md in pages
		createMarkdownTemplate("blog.html", "pages/blog.md", true)

		// With blog-post.html we create welcome-to-my-blog.md in site blog directory
		createMarkdownTemplate("blog-post.html", "blog/welcome-to-my-blog.md", true)
	}

}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Scaffolds a new website",
	Long: `Scaffolds a new website and creates markdown files based on theme chosen.
    
    Uses --theme flag to specify theme to setup site with, or the included 'default' theme
    `,
	Run: func(cmd *cobra.Command, args []string) {
		domain = strings.Join(args, " ")
		setBasePath()
		scaffoldSite()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&theme, "theme", "", "default", "The theme to use with new site")
}
