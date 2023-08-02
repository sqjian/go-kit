package pool

type Pool interface {
	Get() (connection any, err error)
	Put(connection any) (err error)
	Release()
}
