# Facil. A static website generator

## Overview 
Most static website generators focus on blog copy - and why not, blog publishing is hard to do easily. 

Facil is my (experimental) static website generator which focusses on normal websites. You know the ones that are not blogs.

Facil is very early stage, but all the concepts I seem to need are proven and roughly working. The code is quite ugly however.


## Installing

You need to get and install using the Go tool at the moment. Maybe binaries will follow if/when project is better tested & stable. 

## Usage

Facil provides a command line interface (CLI) for all functionality courtesy of @SPF13's cobra package.

Currently Facil exposes the following command set:-

- ```facil setup``` : creates sites and themes directory in currenty directory. Adds a (very) basic default theme to the themes folder
       
- ```facil start --theme theme yourwebsite.domain``` : Scaffolds the directory and file structure for yourwebsite.domain into the sites directory. --theme is optional, omitting means site is scaffolded to use the default theme installed with 'facil setup' 
    
- ```facil build yourwebsite.domain``` : Builds site, parses toml and markdown using the the specified theme template and writes built output to 'compiled' subdirectory.

- ```facil page --template page-name``` : (NOT YET IMPLEMENTED) The intent is to scaffold a new toml/markdown page based on the chosen theme template.

## Themes
A theme is a collection of template files, js, and CSS assets.

## Template files

Template files belong to a theme. They are basically just HTML files within which are some special placeholders, often referred to as "tokens"

Facil reads template files and understands the following included placeholders/tokens (referred to as tokens from now):

### Meta tokens
Meta tokens store information about a page. They are formatted as follows:-

```
[[meta name="whatever"]]
```

Within a theme template you might add a Meta token to the ```<title/>``` element like this:-

```
<title>[[meta name ="title"]]</title>
```

### Toml/Markdown files

Work in progress.. 

## Roadmap

- Finish it
- Refactor it, with testing in mind
- Add tests
- Use it myself

## Did you know

Did you know, 'facil' means 'easy' in Catalan?
















