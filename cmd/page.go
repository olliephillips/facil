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
	"strings"

	"github.com/spf13/cobra"
)

var (
	pageName string
	template string
)

func addPage() error {

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

		setBasePath()

		err := addPage()
		if err != nil {
			log.Fatal("Error unable to create new page " + domain)
		}
	},
}

func init() {
	RootCmd.AddCommand(pageCmd)
	startCmd.Flags().StringVarP(&template, "template", "", "default", "The template to use with new page")
}
