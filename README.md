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
