# gls
It’s `ls` + `du` + `tree` with interactive GUI on your terminal! `gls` is created to easily view, filer and search your files and folders with their size whenever you need to open up some storage space. It wouldn’t be wrong to say that `gls` is a minimal yet powerful file manager CLI tool.

##  Installation
Installing `gls` on your machine is pretty simple: just clone the repo and run `cmd/gls.go`:

```bash
$ git clone github.com/ozansz/gls
$ cd gls
$ go run cmd/gls.go
```

> An install script will come in next feature update and you will be able to just run `gls` on your terminal!

## Usage
There are two running modes of `gls`: GUI and text-based.

The GUI mode is interactive and you will be able to use all of the [features](#features) of `gls`, such as searching by text/regular expression, traversing on the file tree, creating/opening/deleting files and many other things,  until you close the program.

The text mode however, is fairly simple and is a literal combination of running `tree` and `du` altogether, with some additional features.

### Default usage (GUI)
The command below runs `gls` with GUI, which is the default mode. It parses the file tree under the specified path along with the file and folder sizes on disk, then shows the tree view of the parsed tree.

```bash
go run cmd/gls.go —-path ~/Downloads
```

![Screenshot of the GUI mode of gls](./img/gui-screenshot.png)

### Text mode
The command below does the same parsing process as the command above does. Except, this one just dumps the parsed tree as a the `tree` command does with the file/folder sizes and permissions, to the terminal.

```bash
go run cmd/gls.go —-nogui —-path ~/Documents
```

## Features
`gls` includes (and still continues to include more) several features that mimic a normal file manager:
* List the files and folders under the specified path, in tree view
* Show current file info: size on disk, permissions, path, MIME type and last modification
* Sort the tree by the size on disk
* Search files/folders by name, using both plaintext and regular expressions
* Ignore specific files/folders by using regular expressions, similar to `.gitignore` style
	* Default ignore file is `.glsignore`, but infinitely many other ignore files can be specified through the CLI [arguments](#command-line-arguments)
* Open files and folders by default programs or executables that you specify
* Copy/paste and move files and folders
* Remove files
* Create (similar to `touch`) and open files to edit
* Walk on the file tree, collapse and expand nodes easily

### GUI shortcuts

| Shortcut           | Command            | Description                                                                                                                                                                |
| ------------------ | ------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `q`, `ESC`, `ˆC`        | quit               | Exits the program                                                                                                                                                          |
| `c`                  | collapse           | Collapses all nodes in the file tree view                                                                                                                                  |
| `e`                  | expand             | Expands all nodes in the file tree view                                                                                                                                    |
| `s`                  | search             | Opens modal to search nodes (files and folders) by name                                                                                                                    |
| `r`                  | regex search       | Same as search, but you can search using regular expressions                                                                                                               |
| `x`                  | restore            | Loads the original file tree view, mostly used after `search` and `regex search`                                                                                               |
| `o`                  | open               | Opens the selected (on hover) file/folder with the default program                                                                                                         |
| `p`                  | open               | Opens modal to specify the executable path which will be used to open the selected (on hover) file/folder                                                                  |
| `BACKSPACE` , `DEL`    | remove             | Removes the selected (on hover) file. Folder removal is currently not supported                                                                                            |
| `m`                  | mark               | Marks/unmarks the selected (on hover) file or folder. Marked nodes can be used later for `duplicate` and `move`                                                                |
| `u`                  | unmark             | Unmarks all the marked files and folders                                                                                                                                   |
| `n`                  | new                | Create and (optionally) open file **(will be available in next update)**                                                                                                       |
| `d`                  | duplicate          | Copy/pastes the marked files and folders to a specified destination. The destination is specified by the text input of the opened modal **(will be available in next update)** |
| `v`                  | move               | Moves the marked files and folders to a specified destination. The destination is specified by the text input of the opened modal **(will be available in next update)**       |
| `TAB`, `SPACE`, `ENTER`  | toggle expand node | Expands the node if currently collapsed, and vice versa, the selected (on hover) file or folder                                                                            |
| `ARROW KEYS`, `SCROLL` | navigate           | Navigates between nodes in the file tree view                                                                                                                              |

### Command line arguments

```bash
--debug
    	Increase log verbosity
--fmt string
   		size formatter, one of bytes, pow10 or none (default "bytes")
--ignore string
    	Comma-separated ignore files that specify which files folders to exclude
--nogui
    	text-only mode
--path string
    	path to run on (required)
--sort
    	sort nodes by size (default true)
--thresh string
    	size filter threshold, e.g. 10M, 100K, etc.
```