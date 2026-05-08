package workers

import (
	"context"
	"log"
	"time"
	"web-hosting/internal/modules/auth/service"

	"github.com/robfig/cron/v3"
	"github.com/samber/do/v2"
)

type (
	Schedule interface {
		StartSchedule()
	}

	schedule struct {
		authService service.AuthService
	}
)

func NewSchedule(injector do.Injector, authService service.AuthService) Schedule {
	return &schedule{authService: authService}
}

func (s *schedule) StartSchedule() {
	log.Println("Worker Starting...")
	c := cron.New(cron.WithLocation(time.Local))

	_, err := c.AddFunc("@midnight", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		s.authService.CleanupExpiredTokens(ctx)
	})

	if err != nil {
		log.Fatal("Failed to schedule cleanup expired tokens:", err)
	}

	c.Start()

	log.Println("Scheduler running...")
}
