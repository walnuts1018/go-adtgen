# go-adtgen

`go-adtgen` は、Go言語向けの直和型 (Sum Types / Algebraic Data Types) や直積型 (Product Types)、およびそれらを安全・快適に扱うためのヘルパー関数を自動生成するツールです。

## 特徴

- **網羅的な型安全パターンマッチング**: 関数の引数による `Match` だけでなく、構造体形式で `Default` フォールバック可能な `Match...Cases` を提供
- **複数の直和型への所属**: 1つの構造体が複数の直和型インターフェース（例: エラー型、イベント型）に同時に所属可能
- **Discriminator JSONシリアライズ**: `discriminator=<key>` 指定によるタグ付きJSONの相互変換（`Marshal` / `Unmarshal`）を自動生成
- **直積型の双方向変換**: 複数の構造体をフラットに合成し、コンストラクタ `New` や射影 `To...` を自動生成

## インストール

```bash
go get -tool github.com/walnuts1018/go-adtgen
```

## 基本的な使い方

1. **生成用ファイルの作成**
   ビルドタグ `//go:build adtgen_generate` を指定したファイル（例: `generate_types.go`）を作成し、生成したい型の定義を記述します。

   ```go
   //go:build adtgen_generate

   package mypkg

   // +adtgen:sum=VariantA,VariantB;discriminator=type
   type MySumType interface{}

   // +adtgen:product=BaseInfo,DetailInfo
   type MyProductType struct{}
   ```

2. **go:generate の追加**
   パッケージ内の通常のGoファイルに以下を追記します。

   ```go
   //go:generate go tool go-adtgen
   ```

3. **コードの生成**

```bash
go generate ./...
```

各 `//go:build adtgen_generate` ファイルから、同一ディレクトリに `<source>_adtgen.go` が生成されます。

---

## 直和型 (Sum Types) の使い方

`// +adtgen:sum=<Variant1>,<Variant2>...` を指定すると、指定された構造体群をバリアントとするインターフェースおよびヘルパー関数が生成されます。

### オプション指定

- `discriminator=<field_name>`: JSONの判別フィールド（例: `type` や `kind`）を指定して `Marshal<Sum>` / `Unmarshal<Sum>` を生成します。
- `options=no-setter`: 共通フィールドの Setter メソッドの生成を無効化します。

```go
// +adtgen:sum=UserCreated,UserUpdated;discriminator=event_type
type UserEvent interface{}
```

### 生成される主な機能

#### 1. 構造体 Matcher (`Match...Cases` & `Default`)

すべてのバリアントを網羅するハンドラ、または一部のハンドラと `Default` フォールバックを指定できます。

```go
result := MatchUserEventCases(event, UserEventCases[string]{
    UserCreated: func(e UserCreated) string {
        return "作成: " + e.UserID
    },
    Default: func(e UserEvent) string {
        return "その他のイベント"
    },
})
```

戻り値が2つの場合は `MatchUserEventCases2` と `UserEventCases2` を利用できます。

#### 2. 関数引数 Matcher (`Match...` / `Visit...`)

すべてのバリアントの処理を引数として渡す網羅的パターンマッチングです。

```go
result := MatchUserEvent(event,
    func(created UserCreated) string { return created.UserID },
    func(updated UserUpdated) string { return updated.UserID },
)

// 戻り値不要の走査
VisitUserEvent(event,
    func(created UserCreated) { log.Println("created", created) },
    func(updated UserUpdated) { log.Println("updated", updated) },
)
```

#### 3. 安全なキャスト (`As<Sum><Variant>`)

直和型インターフェースから特定の具象型を安全に取り出します。

```go
if created, ok := AsUserEventUserCreated(event); ok {
    fmt.Println("Created ID:", created.UserID)
}
```

#### 4. JSON シリアライズ / デシリアライズ

`discriminator` を指定した場合、判別フィールドを自動で付与・解析する `Marshal` / `Unmarshal` 関数が生成されます。

```go
var event UserEvent = &UserCreated{UserID: "123"}

// {"event_type":"UserCreated","UserID":"123"}
data, err := MarshalUserEvent(event)

decoded, err := UnmarshalUserEvent(data)
```

※ `discriminator` 未指定の場合は、各バリアントの構造に最も一致する型へ自動デコードする `Unmarshal` が生成されます。

#### 5. 共通フィールドへのアクセス

すべてのバリアントが共通のフィールドを持っている場合、インターフェースに `Get<Field>` および `Set<Field>`（`no-setter` 未指定時）が自動生成されます。

---

## 直積型 (Product Types) の使い方

`// +adtgen:product=<Struct1>,<Struct2>...` を使用すると、指定したすべての構造体のフィールドをフラットに結合した新しい構造体が生成されます。

### 定義例

```go
type Base struct { ID string }
type Detail struct { Name string }

// +adtgen:product=Base,Detail
type Combined struct{}
```

### 生成される主な機能

#### コンストラクタ (`New<Type>`)

元の構造体群から結合された構造体を生成します。

```go
b := Base{ID: "123"}
d := Detail{Name: "Item"}
c := NewCombined(b, d)
```

#### 抽出メソッド (`To<Struct>`)

結合された構造体から元の構造体へ射影・変換します。

```go
base := c.ToBase()
detail := c.ToDetail()
```
