
## MFRep
MFRep (multi-file replace) is a cli tool for quick automated editing of file contents written on Go.

## Installation (Linux and MacOS)
### 1. Clone this repo:
```bash
git clone https://github.com/leshless/mfrep.git
```
### 2. Compile and install to your GOPATH folder:
```bash
cd ./mfrep
go build
go install
```
### 3. Add your $GOPATH/bin directory to the $PATH variable:
```bash
echo export PATH=$PATH:$GOPATH/bin >> ~/.bashrc
```
### 4. Now mfrep can be lauched through command line:
```bach
mfrep
```


## Usage

```bash
USAGE:
  mfrep [options] <search_regexp> <replace_string>

OPTIONS:
  --path <path_regexp>, -p <path_regexp>  Regular expression to specify which files should be affected.

  --details, -d           Whenever the output should list affected files. (default: false)
  --silent, -s            Whenever the output should provide no summary. (default: false)
  --recursive, -r         Whenever the files in subdirectories should be affected too. (default: false)

  --help, -h              Show help
```

## Description

Iterates through the current working directory contents and marks the files which will be affected by the replace. If the `--path` option is specified, only the files with names satisfying `<path_regexp>` will be marked.

In each file finds all substring, that satisfy `<search_regexp>` and replaces them with `<replace>`.

Notice that the `<replace>` string can contain default Sprintf placeholders (`%v` or `%s`) for submatches of `<search_regexp>`. The number of capturing groups in regular expression should be equal to the number of placeholders.

For recursive iteration through subdirectories use `--recursive`

For seeing full action summary use `--details`

For seeing no summary use `--silent` 

For seeing help use `--help`
