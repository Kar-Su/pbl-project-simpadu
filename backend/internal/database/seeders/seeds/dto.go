package seeds

import "github.com/google/uuid"

type KelasMahasiswaSeedRequest struct {
	KelasName   string    `json:"kelas_name"`
	MahasiswaID uuid.UUID `json:"mahasiswa_id"`
}
