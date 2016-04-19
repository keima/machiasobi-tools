How to use qbg
==========

[favclip/qbg](https://github.com/favclip/qbg)の使い方。

0. go get

```
$ go get -u github.com/favclip/qbg/cmd/qbg
```

1. structにコメントタグを付ける

```go

// +qbg
type Hoge struct {
    Fuga string
}
```

2. コメントをつけたモデルがあるディレクトリで`qbg`コマンドを実行

```
$ cd path/to/model
$ qbg
```

3. おわり

```
$ ll
total 32
-rw-r--r--  1 keima  staff   312 Apr 19 20:37 menu.go
-rw-r--r--  1 keima  staff  4816 Apr 19 20:38 menuitem_query.go
-rw-r--r--  1 keima  staff    13 Apr 19 19:28 routing.go
```
