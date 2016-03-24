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
	"strings"

	"github.com/spf13/cobra"
)

var project string

func processMarkdownTemplate() {
	// Get the template from file

	// Read template from theme

	// Merge markdown file and template file into new output

	// Write new file to complied folder
}

func buildProject() {
	// Establish target directory based on project

	// Copy theme assets to compiled folder

	// Traverse projects pages folder

	// Parse each markdown file

	//
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
