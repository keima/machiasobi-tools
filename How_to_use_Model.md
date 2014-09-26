# How to use Model
`github.com/knightso/base/gae/model` を組み込むと、
痒い所に手が届く`datastore`Entityを作成することが出来ます。

この使い方をまとめておきます。

## structの定義方法
```go
type Greeting struct {
    Meta // ←
    Id      string `datastore:"-" json:"id"` // for output json
    Name    string `json:"name"`
    Message string `json:"message"`
}
```
`Id`は個人的に必要だったので付けましたが、要らないなら用意しなくても良い。

### Metaが提供するもの
```go
type Meta struct {
	Key       *datastore.Key `datastore:"-" json:"-"`
	Version   int            `json:"version"`
	Deleted   bool           `json:"deleted"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
```
これらの最低限必要そうなものは定義済みなので、わざわざ自分で作らなくても良い。

### Metaがやってくれること
- Key
 - SetするときにこのKeyを用いる
 - GetしたときにEntityに対応するKeyが入っている
- Version
 - Putが呼ばれるたびにincrementされる。
- Deleted
 - これを執筆している段階(2014/09/26)では何をしているのか特定できていない
 - 論理削除フラグっぽいけど・・・？
- CreatedAt
- UpdatedAt
 - 作成日時と更新日時が自動で割り当てられる

加えて、modelが`github.com/qedus/nds`を使用してmemcacheを使用したDBアクセスを提供する。

## データの格納
いつもどおりです。
```go
g1 := Greeting {
    Name: "Kouta Imanaka"
    Message: "Hello, Gopher!"
}
```

`Id`はデータ返すときに使う予定なので、今は定義しません。

## Keyの作成
いつもどおりです。
```go
key := datastore.NewIncompleteKey(c, "Greeting", nil)
// 以下でも良い
// key := datastore.NewKey(c, "Greeting", "hogefugahogehoge", 0, nil)
```

## PUT,GET
```go
// PUT
g1.SetKey(key)
err1 := model.Put(c, g1)
if err != nil {
   http.error(...)
   return
}

// GET
g2 := Greeting{}
err2 := model.Get(c, key, &g2)
if err != nil {
   http.error(...)
   return
}

// GET(version check)
g3 := Greeting{}
err3 := model.GetWithVersion(c, key, g2.GetVersion(), &g3)
if err3 != nil {
   err3 == OptimisticLockError {
       // 取得しようとしたEntityのバージョンが指定値より大きかった時、
       // どこからか書き込みが行われたということになる
   } else {
       http.error(...)
   }
   return
}
```

GETしたGreetingを、`go-json-rest`などで投げ返すときに、一手間加える。

```go
g1.Id = g1.GetKey().StringId() // or IntId()
w.WriteJson(&g1)
```

これで、Entityを特定するJSONを返すことが出来たりする。

## QUERY
```go
greetings := make([]Greeting, 0, 10)
q := datastore.NewQuery("Greeting").Order("-UpdatedAt").Limit(10)
model.ExecuteQuery(c, &q, &greetings)

// set id
for index,value := range greetings {
    greetings[index].Id = value.GetKey().StringId()
}

w.WriteJson(&greetings)
```
Orderに存在しない変数名（フィールド名）を入れてしまうと、結果が0件となるのは地味にハマるので注意。

あと、for文もvalueのフィールドに代入してもgreetingsには変更が反映されないのもハマる。

## GetMulti, PutMulti
TBD

## FindIndexedKey, PutKeyToIndex
TBD (というかIndex関連よく分かってない、、、)

