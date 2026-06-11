package workers

import (
	"context"
	"log"
	"sync"
	"time"
	authService "web-hosting/internal/modules/auth/service"
	presensiService "web-hosting/internal/modules/presensi/service"

	"github.com/robfig/cron/v3"
	"github.com/samber/do/v2"
)

type (
	Schedule interface {
		StartAll()
		InitCleaningTokens()
		InitPresensiReminder()
		StartPresensi()
		StopPresensi()
		IsPresensiRunning() bool
	}

	schedule struct {
		authService     authService.AuthService
		presensiService presensiService.PresensiService

		cleanupCron  *cron.Cron
		presensiCron *cron.Cron
		presensiMu   sync.RWMutex

		presensiRunning bool
	}
)

func NewSchedule(injector do.Injector, authService authService.AuthService, presensiService presensiService.PresensiService) Schedule {
	return &schedule{
		authService:     authService,
		presensiService: presensiService,
		cleanupCron:     cron.New(cron.WithLocation(time.Local)),
		presensiCron:    cron.New(cron.WithLocation(time.Local))}
}

func (s *schedule) StartAll() {
	s.InitCleaningTokens()
	s.InitPresensiReminder()

	s.cleanupCron.Start()
	s.StartPresensi()

	log.Println("Workers: Scheduled cleanup and presensi reminder jobs")
}

func (s *schedule) InitCleaningTokens() {
	_, err := s.cleanupCron.AddFunc("@midnight", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.authService.CleanupExpiredTokens(ctx)
	})
	if err != nil {
		log.Fatal("Workers: Failed to schedule cleanup job:", err)
	}
}

func (s *schedule) InitPresensiReminder() {
	_, err := s.presensiCron.AddFunc("0 6 * * 1-5", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		s.presensiService.CreatePresensi(ctx, nil)
		log.Println("Workers: Open Presensi Pegawai")
	})
	if err != nil {
		log.Fatal("Workers: Failed to schedule presensi reminder:", err)
	}
}

func (s *schedule) StartPresensi() {
	s.presensiMu.Lock()
	defer s.presensiMu.Unlock()

	if s.presensiRunning {
		log.Println("Workers: Presensi worker is already running.")
		return
	}

	s.presensiCron.Start()
	s.presensiRunning = true
	log.Println("Workers: Presensi Scheduler STARTED.")
}

func (s *schedule) StopPresensi() {
	s.presensiMu.Lock()
	defer s.presensiMu.Unlock()

	if !s.presensiRunning {
		log.Println("Workers: Presensi worker is already stopped.")
		return
	}

	s.presensiCron.Stop()
	s.presensiRunning = false
	log.Println("Workers: Presensi Scheduler STOPPED.")
}

func (s *schedule) IsPresensiRunning() bool {
	s.presensiMu.RLock()
	defer s.presensiMu.RUnlock()
	return s.presensiRunning
}
