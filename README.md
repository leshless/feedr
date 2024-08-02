### feedr
Small RSS reader CLI app

## Installation (Linux)
# 1. Clone this repo:
```bash
    git clone https://github.com/leshless/feedr.git
```
# 2. Compile and install to your GOPATH folder:
```bash
    cd ./feedr
    go build
    go install
```
# 3. Add your $GOPATH/bin directory to the $PATH variable:
```bash
    echo export PATH=$PATH:$GOPATH/bin >> ~/.bachrc
```
# 4. Now feedr can be lauched through command line:
```bach
    feedr
``