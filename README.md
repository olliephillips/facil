# Facil. A static website generator

## Overview 
Most static website generators focus on blog copy - and why not, blog publishing is hard to do easily. 

Facil is my (experimental) static website generator which focusses on normal websites. You know, the ones that are not blogs.

Facil is very early stage, but all the concepts I seem to need are proven and roughly working. The code is currently quite ugly however.

## Installing

You need to get and install using the Go tool at the moment. Maybe binaries will follow if/when project is better tested & stable. 

## Usage

Facil provides a command line interface (CLI) for all functionality courtesy of @SPF13's cobra package.

Currently Facil exposes the following command set:-

- ```facil setup``` : creates sites and themes directory in current directory. Adds a (very) basic default theme to the themes folder
       
- ```facil start --theme theme yourwebsite.domain``` : Scaffolds the directory and file structure for yourwebsite.domain into the sites directory. --theme is optional, omitting means site is scaffolded to use the default theme installed with 'facil setup' 
    
- ```facil build yourwebsite.domain``` : Builds site, parses TOML and markdown using the the specified theme template and writes built output to 'compiled' subdirectory.

- ```facil page --template page-name``` : (NOT YET IMPLEMENTED) The intent is to scaffold a new TOML/markdown page based on the chosen theme template.

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
[[meta name="whatever"]]
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
Elements can provide copy for any HTML element in a template. To add a token something like this to your template.


```
[[element name="title" description="Set the title"]]
```

### Partial tokens

Partial tokens allow the inclusion of content that is used in multiple places in the site. For example if you have three templates, default.html, left-sidebar.html and right-sidebar, they may share some elements such as a footer. In this scenario it is sensible to use a partial to create this content once but include it in all three templates

```
[[partial name="footer"]]
```

#### Example template html

The below demonstrates how the above tokens are uses in a template html file

```
<html>
    <head>
        <title>[[meta name="title"]]</title>
        <meta name="description" content="[[meta name="description"]]">
        <meta name="anything" content="[[meta name="random"]]"
    </head>
    <body>
        <div id="nav">
            [[navigation]]
        </div>
        <div id="body">
            <div id="title">
                [[element name="title" description="Set the title"]]
            </div>
            <div id="intro">
                [[element name="introduction" description="Add an introductory paragraph"]]
            </div>
            <div id="footer>"
                [[partial name="footer"]]
            </div>
        </div>
    </body>
</html>
```

## TOML/Markdown files



Work in progress.. 

## Roadmap

- Finish it
- Refactor it, with testing in mind
- Add tests
- Use it myself
- Maybe extend to include blog type copy

## Did you know

Did you know, 'facil' means 'easy' in Catalan?
















