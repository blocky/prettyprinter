package prettyprinter

type FieldError struct {
	Error error
}

func MakeFieldError(err error) FieldError {
	return FieldError{err}
}

func (be FieldError) MarshalJSON() ([]byte, error) {
	val := errorToString(be.Error)
	if val == "" {
		return prettyJSON(nil)
	}
	return prettyJSON(val)
}

type KeyValueError struct {
	Error FieldError `json:"err"`
}

func MakeKeyValueError(err error) KeyValueError {
	return KeyValueError{MakeFieldError(err)}
}

func (kve *KeyValueError) MarshalJSON() ([]byte, error) {
	type Alias KeyValueError
	return prettyJSON(&struct {
		*Alias
	}{
		Alias: (*Alias)(kve),
	})
}
