machitools (まち☆つーるず)
==============
API Backend System for Machi ★ Appli (マチ★アプリ)


```
// setup backend
$ direnv allow
  (or `export GOPATH=$GOPATH:/path/to/machitools/repo` )
$ cd src
$ goapp get
$ cd -

// setup frontend
$ cd static
$ npm install
$ bower install
$ cd -

$ goapp serve ./src
$ cd static
$ gulp default (or gulp watch)
```
