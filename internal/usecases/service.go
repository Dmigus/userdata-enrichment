package usecases

import "context"

type UpdateRequest struct {
	age, sex, nat bool
	newAge        int
	newSex        string
	newNat        string
}

// (name, surname, patronymic, age, sex, nationality)
type Transaction interface {
	KVOps
	PushQueue
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type KVOps interface {
	Update(ctx context.Context, key Key, ur UpdateRequest) error
	Create(ctx context.Context, key Key) error
	Delete(ctx context.Context, key Key) error
}

type TransactionalStorage interface {
	BeginTx(ctx context.Context) (Transaction, error)
}

type PushQueue interface {
	Push(ctx context.Context, key Key) error
}

type Service struct {
	db    TransactionalStorage
	queue PushQueue
}

func New(db TransactionalStorage, queue PushQueue) *Service {
	return &Service{
		db:    db,
		queue: queue,
	}
}

func (s *Service) Update(ctx context.Context, key Key, ur UpdateRequest) error {
	tr, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tr.Rollback(ctx)
	if err = tr.Update(ctx, key, ur); err != nil {
		return err
	}
	return tr.Commit(ctx)
}

func (s *Service) Create(ctx context.Context, key Key) error {
	tr, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tr.Rollback(ctx)
	if err = tr.Push(ctx, key); err != nil {
		return err
	}
	if err = tr.Create(ctx, key); err != nil {
		return err
	}
	return tr.Commit(ctx)
}

func (s *Service) Delete(ctx context.Context, key Key) error {
	tr, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tr.Rollback(ctx)
	if err = tr.Delete(ctx, key); err != nil {
		return err
	}
	return tr.Commit(ctx)
}
