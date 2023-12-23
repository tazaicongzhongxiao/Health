package log

type OrmLoggerI struct {
}

func (o *OrmLoggerI) Print(v ...interface{}) {
	v = v[2:] //只取sql语句
	Info(" %v  ", v)
}

func NewGormLogger() *OrmLoggerI {
	return &OrmLoggerI{}
}
