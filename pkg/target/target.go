package target

type Target interface {
	Write(interface{}) error
}
