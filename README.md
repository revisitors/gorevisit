# gorevisit
--
    import "github.com/revisitors/go.revisit.link"


## Usage

```go
var (
	ErrNotImplemented = errors.New("not implemented yet")
)
```

#### type ApiMsg

```go
type ApiMsg struct {
	Content ContentInfo
	Meta    struct{ Audio ContentInfo }
}
```


#### func  NewApiMsgFromJson

```go
func NewApiMsgFromJson(b []byte) (*ApiMsg, error)
```

#### func (*ApiMsg) Json

```go
func (a *ApiMsg) Json() ([]byte, error)
```

#### type ContentInfo

```go
type ContentInfo struct {
	Type string
	Data []byte
}
```
# gorevisit
--
    import "github.com/revisitors/go.revisit.link"


## Usage

```go
var (
	ErrNotImplemented = errors.New("not implemented yet")
)
```

#### type ApiMsg

```go
type ApiMsg struct {
	Content *Content
	Meta    *MetaContent
}
```

ApiMsg is a message containing Content, and MetaContent. the MetaContent should
be audio.

#### func  NewApiMsgFromJson

```go
func NewApiMsgFromJson(b []byte) (*ApiMsg, error)
```
NewApiMsgFromJson returns an ApiMsg struct pointer from a json byte array.

#### func (*ApiMsg) Json

```go
func (a *ApiMsg) Json() ([]byte, error)
```
Json serializes a gorevisit.ApiMsg back to JSON bytes

#### type Content

```go
type Content struct {
	Type string
	Data string
}
```

Content contains a type and a byte array and should be an image

#### type MetaContent

```go
type MetaContent struct {
	Audio *Content
}
```

AudioContent contains a Content pointer
