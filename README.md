# gorevisit
--
    import "github.com/revisitors/gorevisit"


## Usage

#### func  BytesToDataURI

```go
func BytesToDataURI(data []byte, contentType string) string
```
BytesToDataURI, given a byte array and a content type, creates a Data URI of the
content

#### type AudioData

```go
type AudioData struct {
	Data string `json:data"`
}
```

AudioData holds a reference to the data URI of sound data in a Revisit.link
message See: https://developer.mozilla.org/en-US/docs/Web/HTTP/data_URIs

#### type ImageData

```go
type ImageData struct {
	Data string `json:"data"`
}
```

ImageData holds a reference the data URI of image data in a Revisit.link message
See: https://developer.mozilla.org/en-US/docs/Web/HTTP/data_URIs

#### func (ImageData) ByteReader

```go
func (i ImageData) ByteReader() io.Reader
```
ByteReader returns an io.Reader for the image data in a Revisit message

#### type MetaData

```go
type MetaData struct {
	Audio AudioData `json:"audio"`
}
```

MetaData wraps the Audio data of a Revisit.link message as per the specification
See: http://revisit.link/spec.html

#### type RevisitMsg

```go
type RevisitMsg struct {
	Content ImageData `json:"content"`
	Meta    MetaData  `json:"meta"`
}
```

RevisitMsg holds a decoded Revisit.link message See:
http://revisit.link/spec.html

#### func  ImageRevisitor

```go
func ImageRevisitor(m *RevisitMsg, t func(src image.Image, dst image.RGBA) error) (*RevisitMsg, error)
```
ImageRevisitor, given a RevisitMsg and an image transformation function, runs
the image data through the transformation and returns a new RevisitMsg with the
transformed image

#### func  NewRevisitMsgFromFiles

```go
func NewRevisitMsgFromFiles(mediaPath ...string) (*RevisitMsg, error)
```
NewRevisitMsgFromFiles, given the path to an image file and optional path to an
audio file, creates a JSON encoded Revisit.link message

#### type RevisitService

```go
type RevisitService struct {
}
```

RevisitService holds the necessary context for a Revisit.link service.
Currently, this consists of an imageTransformer

#### func  NewRevisitService

```go
func NewRevisitService(it func(image.Image, image.RGBA) error) *RevisitService
```
NewRevisitService, given an image transformation function, returns a new
Revisit.link service

#### func (*RevisitService) PostHandler

```go
func (rs *RevisitService) PostHandler(w http.ResponseWriter, r *http.Request)
```
PostHandler accepts POSTed revisit messages from a Revisit.link hub, transforms
the message, and returns the transformed message to the hub

#### func (*RevisitService) Run

```go
func (rs *RevisitService) Run()
```
Run starts the Revisit.link service

#### func (*RevisitService) ServiceCheckHandler

```go
func (rs *RevisitService) ServiceCheckHandler(w http.ResponseWriter, r *http.Request)
```
ServiceCheckHandler responts to availability requests from a Revisit.link hub

#### func (*RevisitService) ServiceHandler

```go
func (rs *RevisitService) ServiceHandler(w http.ResponseWriter, r *http.Request)
```
ServiceHandler appropriately routes ervice requests from a Revisit.link hub
