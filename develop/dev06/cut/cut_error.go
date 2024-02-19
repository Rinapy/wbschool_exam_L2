package cut

type IndexValueError struct{}

func (e *IndexValueError) Error() string {
	return "index error"
}

type DataNotFound struct{}

func (e *DataNotFound) Error() string {
	return "data not found, please check your passed arguments"
}
