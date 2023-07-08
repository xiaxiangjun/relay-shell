package common

type Stream interface {
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)
	Close() error
}
