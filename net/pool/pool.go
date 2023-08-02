package pool

type Pool interface {
	Get() (connection interface{}, err error)
	Put(connection interface{}) (err error)
	Release()
}
