package seeder

import (
	"context"
	"encoding/json"
	"github.com/gictorbit/sikabiz/db/userdb"
	"go.uber.org/zap"
	"io"
	"os"
	"sync"
)

type Seeder struct {
	userdb       userdb.UserDBPgConn
	logger       *zap.Logger
	filePath     string
	workerNumber int
}

func NewSeeder(userdb userdb.UserDBPgConn, logger *zap.Logger, filePath string) *Seeder {
	return &Seeder{
		userdb:       userdb,
		logger:       logger,
		filePath:     filePath,
		workerNumber: 10,
	}
}

// Worker function
func (s *Seeder) worker(ctx context.Context, jobs <-chan userdb.User, wg *sync.WaitGroup) {
	defer wg.Done()
	for user := range jobs {
		err := s.userdb.InsertUserData(ctx, user)
		if err != nil {
			s.logger.Error("Failed to insert user",
				zap.Error(err),
				zap.String("userID", user.ID),
			)
		}
	}
}

func (s *Seeder) RunSeeder(ctx context.Context) {
	file, err := os.Open(s.filePath)
	if err != nil {
		s.logger.Error("Failed to open JSON file", zap.Error(err))
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("Failed to read JSON file", zap.Error(err))
		return
	}

	// Parse JSON data
	var users []userdb.User
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		s.logger.Error("Failed to unmarshal JSON", zap.Error(err))
		return
	}

	jobs := make(chan userdb.User, s.workerNumber)
	var wg sync.WaitGroup

	for i := 0; i < s.workerNumber; i++ {
		wg.Add(1)
		go s.worker(ctx, jobs, &wg)
	}

	// Send users to the workers
	for _, user := range users {
		jobs <- user
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
}
