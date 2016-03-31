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

## Roadmap

- Finish it
- Refactor it, with testing in mind
- Add tests
- Use it myself

## Did you know

Did you know, 'facil' means 'easy' in Catalan?
















