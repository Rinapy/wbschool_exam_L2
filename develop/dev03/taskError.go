package main

type ErrCloseFile struct{}

func (e *ErrCloseFile) Error() string {
	return "there was an error closing the transferred file"
}

type ErrOpenFile struct{}

func (e *ErrOpenFile) Error() string {
	return "there was an error opening the transferred file"
}

type ErrReadFile struct{}

func (e *ErrReadFile) Error() string {
	return "there was an error reading the transferred file"
}

type ErrWriteFile struct{}

func (e *ErrWriteFile) Error() string {
	return "there was an error writing the transferred file"
}

type ErrIndexFile struct{}

func (e *ErrIndexFile) Error() string {
	return "index out of range"
}
