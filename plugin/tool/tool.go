package tool

import (
	"bytes"
	"sync"
)

type PluginToolsImpl struct {
	OptsImpl
	BufferImpl
}

type BufferImpl struct {
	buf bytes.Buffer
}

func (t *BufferImpl) Write(in []byte) error {
	_, err := t.buf.Write(in)
	return err
}

func (t *BufferImpl) Read() []byte {
	return t.buf.Bytes()
}

type OptsImpl struct {
	kvs sync.Map
}

func (o *OptsImpl) Set(k interface{}, v interface{}) {
	o.kvs.Store(k, v)
}

func (o *OptsImpl) Get(key interface{}) (interface{}, bool) {
	return o.kvs.Load(key)
}
