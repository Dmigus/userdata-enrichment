package usecases

import "context"

type UpdateRequest struct {
	age, sex, nat bool
	newAge        int
	newSex        string
	newNat        string
}

//(name, surname, patronymic, age, sex, nationality)

type KVStorage interface {
	Update(ctx context.Context, key Key, ur UpdateRequest) error
	Create(ctx context.Context, key Key) error
	Delete(ctx context.Context, key Key) error
}

type PushQueue interface {
	Push(ctx context.Context, key Key) error
}

type Service struct {
	db    KVStorage
	queue PushQueue
}

func (s *Service) Update(ctx context.Context, key Key, ur UpdateRequest) error {
	return s.db.Update(ctx, key, ur)
}

func (s *Service) Create(ctx context.Context, key Key) error {
	err := s.queue.Push(ctx, key)
	if err != nil {
		return err
	}
	return s.db.Create(ctx, key)
}

func (s *Service) Delete(ctx context.Context, key Key) error {
	return s.db.Delete(ctx, key)
}
