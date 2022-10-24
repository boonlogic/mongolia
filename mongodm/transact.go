package mongodm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// transact runs the passed-in function inside a Transaction.
func transact(ctx context.Context, fn func(context.Context) error) error {
	sess, err := options2.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(ctx)

	txn, err := makeTransaction(sess, fn)
	if err != nil {
		return err
	}
	if err = mongo.WithSession(ctx, sess, txn); err != nil {
		if aerr := sess.AbortTransaction(ctx); aerr != nil {
			return aerr
		}
		return err
	}
	return nil
}

// makeTransaction decorates a func so that it runs as a Transaction under the given mongo.Session.
func makeTransaction(sess mongo.Session, fn func(context.Context) error) (func(mongo.SessionContext) error, error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	opts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	txn := func(ctx mongo.SessionContext) error {
		if err := sess.StartTransaction(opts); err != nil {
			return err
		}
		if err := fn(ctx); err != nil {
			return err
		}
		if err := sess.CommitTransaction(ctx); err != nil {
			return err
		}
		return nil
	}
	return txn, nil
}
