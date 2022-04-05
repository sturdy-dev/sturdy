package file

type Type string

const (
	UnknownType Type = ""
	TextType    Type = "text"
	BinaryType  Type = "binary"
	ImageType   Type = "image"
)
