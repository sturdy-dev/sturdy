package blobs

type ID string

type Blob struct {
	ID   ID     `db:"id"`
	Data []byte `db:"data"`
}
