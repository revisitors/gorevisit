# gorevisit
--
    import "github.com/revisitors/gorevisit"


## Usage

```go
var (
	//ErrUnsupportedType is returned when a Transform does not support the type(s) passed to it
	ErrUnsupportedType = errors.New("unsupported type")
)
```

#### type APIMsg

```go
type APIMsg struct {
	Content *Content     `json:"content"`
	Meta    *MetaContent `json:"meta"`
}
```

APIMsg is a message containing Content, and MetaContent. the MetaContent should
be audio.

#### func  NewAPIMsgFromJSON

```go
func NewAPIMsgFromJSON(b []byte) (*APIMsg, error)
```
NewAPIMsgFromJSON returns an APIMsg struct pointer from a json byte array.

#### func (*APIMsg) JSON

```go
func (a *APIMsg) JSON() ([]byte, error)
```
JSON serializes a gorevisit.APIMsg back to JSON bytes

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

MetaContent contains a Content pointer

#### type Transformer

```go
type Transformer interface {
	Transform(*APIMsg) (*APIMsg, error)
}
```

Transformer interface transforms an APIMsg into another APIMsg
