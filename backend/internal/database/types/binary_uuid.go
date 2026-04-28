package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type BinaryUUID uuid.UUID

func (u BinaryUUID) Value() (driver.Value, error) {
	return uuid.UUID(u).MarshalBinary()
}

func (u *BinaryUUID) Scan(value any) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan BinaryUUID: %v", value)
	}

	parseUuid, err := uuid.FromBytes(bytes)
	if err != nil {
		return err
	}

	*u = BinaryUUID(parseUuid)
	return nil
}

func (u BinaryUUID) String() string {
	return uuid.UUID(u).String()
}
