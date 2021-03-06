# Facil. A static website generator

## Overview 
Most static website generators focus on blog copy - and why not, blog publishing is hard to do easily. 

Facil is my (experimental) static website generator which focusses on normal websites. You know, the ones that are not blogs.

Facil is very early stage, but all the concepts I seem to need are proven and roughly working. The code is currently quite ugly however.

## Example
The project's website can now be found [here](https://olliephillips.github.io/facil/). To be fair there were some hoops to jump through to get dependencies paths working on Github pages - see [this issue](https://github.com/olliephillips/facil/issues/6).

The website was built using a free responsive HTML theme which was converted to work with Facil. The content of that site was compiled to static HTML using Facil.

## Installing

You need to get and install using the Go tool at the moment. Maybe binaries will follow if/when project is better tested & stable. 

## Usage

Facil provides a command line interface (CLI) for all functionality courtesy of @SPF13's cobra package.

Currently Facil exposes the following command set:-

- ```facil setup``` : creates sites and themes directory in current directory. Adds a (very) basic default theme to the themes folder
       
- ```facil start --theme theme yourwebsite.domain``` : Scaffolds the directory and file structure for yourwebsite.domain into the sites directory. --theme is optional, omitting means site is scaffolded to use the default theme installed with 'facil setup' 
    
- ```facil build yourwebsite.domain``` : Builds site, parses TOML and markdown using the the specified theme template and writes built output to 'compiled' subdirectory.

- ```facil page --template template page-name``` :  The intent is to scaffold a new TOML/markdown page based on the chosen theme template.

## Themes
A theme is a collection of template files, JavaScript, CSS and image assets.

New themes should be added to the ```/themes``` directory. Only a ```default.html``` template is mandatory.
Partial templates should be placed in a `partials` folder. An example of simple structure which includes two page templates, one partial template, together with JavaScript, CSS and image assets, would be:

- theme_name/
  - partials/footer.html
  - css/
  - js/
  - images/
  - default.html
  - left_sidebar.html
  
## Template files

Template files belong to a theme. A theme must have a minimum of one tempate file, named ```default.html``` though it can have as many as desired.They are basically just HTML files within which are some special placeholders, often referred to as "tokens"

Facil reads template files and understands the following included placeholders/tokens (referred to as tokens from now):

### Meta tokens
Meta tokens store information about a page. They are formatted as follows:-

```
[[meta name="title"]]

```

Currently, these meta tokens are supported:-

```
[[meta name="title"]]
[[meta name="description"]]
[[meta name="keywords"]]
[[meta name="author"]]
```

Within a theme template you might add a Meta token to the ```<title/>``` element like this:-

```
<title>[[meta name="title"]]</title>
```

### Navigation token
This token, when found in a template, will be replaced with a HTML unordered list of navigation options, built using the page hierarchy.

```
[[navigation]]

```
### Element tokens
Elements can provide copy for any HTML element in a template. Copy can be HTML or simple text. The `type` attribute differentiates the two in the theme template. Text type elements are not processed as markdown, so add no extra html markup to the page.

To add a HTML element token, add something like this to your template. 

```
[[element type="html" name="title" description="Set the title"]]
```

To add a text element token, do this.

```
[[element type="text" name="title" description="Set the title"]]
```

Note the `description` attribute, think of this as a note or tip that provides a steer as to what content should be entered for the token.

### Partial tokens

Partial tokens allow the inclusion of content that is used in multiple places in the site. For example if you have three templates, default.html, left-sidebar.html and right-sidebar, they may share some elements such as a footer. In this scenario it is sensible to use a partial to create this content once but include it in all three templates

```
[[partial name="footer"]]
```

#### Example template html

The below demonstrates how the above tokens are used in a template html file

```
<html>
    <head>
        <title>[[meta name="title"]]</title>
        <meta name="description" content="[[meta name="description"]]">
        <meta name="keywords" content="[[meta name="keywords"]]">
        <meta name="author" content="[[meta name="author"]]">
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
            <div id="footer>"
                [[partial name="footer"]]
            </div>
        </div>
    </body>
</html>
```

## Creating a new site

A new site is created using the `facil start --theme theme yourwebsite.domain` command. As mentioned briefly this command scaffolds the directory and file structure for the new site. To create a new website, for the domain mywebsite.com, using the default theme, we'd use this command

```
facil start mywebsite.com
```

Which would create the following directory and file structure in the sites folder:

- mywebsite.com/
 - compiled/
 - pages/index.md
 - partials/footer.md
 - theme/default/
 - config.toml
 

Working through each:

The `compiled` directory is where the built site will be placed. The build process basically merges the TOML/markdown files with the chosen template, to create a new pure HTML file.

The `pages` directory is where all your TOML/markdown files go. Each is a page on your site. When a new site is scaffolded, the chosen themes `default.html` template is used to create `index.md` in this folder. This will become index.html, or the homepage of the website, when the site is built.  We look at the format of the TOML/markdown files below.

The `partials` directory will include one markdown file for each partial template in the theme. Here we have just a `footer.md` file. There is no TOML config information in a partial markdown file.

The `theme` directory contains the theme in use with this site. It is copied from the `themes` folder so that site level customizations can be made to the theme. In this example the theme is the `default` theme.
 
## config.toml

`config.toml` is a file containing information about the site. It is also automatically generated but can be manually edited.

```
domain = "mywebsite.com"
theme = "default"
https = "off" # Options are off, on
pretty = "off" # Options are off, on

```

Setting `https` to on will generate the sitemap with a `https` prefix instead of `http`
Setting `pretty` to on will generate pages, nav and sitemap in pretty url form. In this cases a folder takes on the page name, and the page file is named `index.html`. Both the navigation and sitemap omit the `index.html`


## TOML/Markdown files

Our TOML/Markdown files have the .md file extension. `index.md` is created when the site is first scaffolded. Additional files can be created manually, or better yet with the `facil page` command.

Assuming, the example template html shown above was the `default.html` template in the theme, the following content would be generated in `pages/index.md`:-

```
+++

[Meta]
title = ""
description = ""
keywords = ""
author = ""

[Navigation]
text = ""
order = "99"

[Design]
template = "default"

+++

***TEXT*** Title (Set the title)

Your Title syntax here

***

***HTML*** Introduction (Add an introductory paragraph)

# Your Introduction markdown/html syntax here

***

```

## Sitemap creation

Each time a site is built with the `build` command, a gzipped sitemap is created in the root (sitemap.xml.gz).
All URLs included in the sitemap will be prefixed with 'http://' by default. 
Sites using TLS should set their config.toml `https` property to "on" so that URLs will instead be prefixed with 'https://'.

## Roadmap

- Address issues log
- Refactor it, with testing in mind
- Add tests
- Use it myself

## Why static websites?

Amongst other things, speed, security and simple hosting requirements are key reasons for hosting static rather than dynamic pages.

## Did you know

Did you know, 'facil' means 'easy' in Catalan?
















