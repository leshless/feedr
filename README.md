# Feedr

## 🌐Get all the news in one app!
### Feedr is a compact RSS reader CLI app written on Go
 - **Grab the latest updates from the news channels straight into your linux terminal**
 - **Manage your feed by editing the list of allowed RSS hosts**
 - **Setup the Feedr config to always stay tuned**

## ⚙️Installation (Linux)
### 1. Clone this repo:
```bash
git clone https://github.com/leshless/feedr.git
```
### 2. Compile and install to your GOPATH folder:
```bash
cd ./feedr
go build
go install
```
### 3. Add your $GOPATH/bin directory to the $PATH variable:
```bash
echo export PATH=$PATH:$GOPATH/bin >> ~/.bashrc
```
### 4. Now feedr can be lauched through command line:
```bach
feedr
```