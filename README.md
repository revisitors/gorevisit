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
	Content *Content     `json:"content"`
	Meta    *MetaContent `json:"meta"`
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
	Type string `json:"type"`
	Data string `json:"data"`
}
```

Content contains a type and a byte array and should be an image

#### type MetaContent

```go
type MetaContent struct {
	Audio *Content `json:"audio"`
}
```

AudioContent contains a Content pointer
