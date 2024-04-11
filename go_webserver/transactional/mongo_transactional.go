package transactional

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type mongoTransactional struct {
	client  *mongo.Client
	session mongo.Session
}

// NewMongoTransactional creates a new MongoDB transaction manager.
func NewMongoTransactional(client *mongo.Client) Transactional {
	return &mongoTransactional{
		client: client,
	}
}

func (m *mongoTransactional) BeginTransaction(ctx context.Context) (TransactionContext, error) {
	session, err := m.client.StartSession()
	if err != nil {
		return nil, err
	}
	m.session = session

	txnOpts := options.Transaction().SetReadConcern(readconcern.Local()).SetWriteConcern(writeconcern.Majority())
	err = session.StartTransaction(txnOpts)
	if err != nil {
		return nil, err
	}

	// Use the session as part of your transaction context, ensuring operations use this session.
	txnCtx := mongo.NewSessionContext(ctx, session)
	return txnCtx, nil
}

func (m *mongoTransactional) Commit(ctx context.Context) error {
	if m.session != nil {
		err := m.session.CommitTransaction(ctx)
		m.session.EndSession(ctx)
		return err
	}
	return errors.New("no session found, please start a transaction before committing")
}

func (m *mongoTransactional) Rollback(ctx context.Context) error {
	if m.session != nil {
		err := m.session.AbortTransaction(ctx)
		m.session.EndSession(ctx)
		return err
	}
	return errors.New("no session found, please start a transaction before rolling back")
}
