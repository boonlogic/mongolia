package mongolia

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
)

type ConnectOptions struct {
	options.ClientOptions

	Timeout *time.Duration
}

func (o *ConnectOptions) ApplyURI(uri string) *ConnectOptions {
	o.ClientOptions = *o.ClientOptions.ApplyURI(uri)
	return o
}

func (o *ConnectOptions) SetTimeout(timeout time.Duration) *ConnectOptions {
	o.Timeout = &timeout
	return o
}

func Connect(uri string, dbname string, timeout time.Duration) error {
	opts := options.Client().ApplyURI(uri)
	conf := &mgm.Config{
		CtxTimeout: timeout,
	}
	if err := mgm.SetDefaultConfig(conf, dbname, opts); err != nil {
		return err
	}
	return nil
}
