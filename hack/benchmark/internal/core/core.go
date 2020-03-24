package core

// TODO: delete this interface
type Core interface {
	Search(vec, size, epsilon, radius interface{}) (interface{}, error)
	Insert(vec interface{}) (interface{}, error)
	InsertCommit(vec, poolSize interface{}) (interface{}, error)
	BulkInsert(vecs interface{}) (interface{}, []error)
	BulkInsertCommit(vecs, poolSize interface{}) (interface{}, []error)
	CreateAndSaveIndex(poolSize interface{}) error
	CreateIndex(poolSize interface{}) error
	SaveIndex() error
	Remove(id interface{}) error
	BulkRemove(ids interface{}) error
	GetVector(id interface{}) (interface{}, error)
	Close()
}

type Core32 interface {
	Search(vec []float32, size int, epsilon, radius float32) (interface{}, error)
	Insert(vec []float32) (uint, error)
	InsertCommit(vec []float32, poolSize uint32) (uint, error)
	BulkInsert(vecs [][]float32) ([]uint, []error)
	BulkInsertCommit(vecs [][]float32, poolSize uint32) ([]uint, []error)
	CreateAndSaveIndex(poolSize uint32) error
	CreateIndex(poolSize uint32) error
	SaveIndex() error
	Remove(id uint) error
	BulkRemove(ids ...uint) error
	GetVector(id uint) ([]float32, error)
	Close()
}

type Core64 interface {
	Search(vec []float64, size int, epsilon, radius float64) (interface{}, error)
	Insert(vec []float64) (uint, error)
	InsertCommit(vec []float64, poolSize uint32) (int, error)
	BulkInsert(vecs [][]float64) ([]int, []error)
	BulkInsertCommit(vecs [][]float64, poolSize int) ([]uint, []error)
	CreateAndSaveIndex(poolSize int) error
	CreateIndex(poolSize int) error
	SaveIndex() error
	Remove(id uint) error
	BulkRemove(ids ...uint) error
	GetVector(id int) ([]float64, error)
	Close()
}
