package schedule

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// Schedule holds the on time and the offtime
type Schedule struct {
	OnTime  time.Time
	OffTime time.Time
}

// TimeUntilOff returns the time until the schedule should be off
func (s *Schedule) TimeUntilOff() *time.Duration {
	until := time.Until(s.OffTime)
	return &until
}

// TimeUntilOn returns the time until the schedule should be off
func (s *Schedule) TimeUntilOn() *time.Duration {
	until := time.Until(s.OnTime)
	return &until
}

// Init initialized the bucket inside the bolt database
func Init(db *bolt.DB) error {
	err := scheduleOnBucket(db)
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	return nil
}

func scheduleOnBucket(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Schedule"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return err
}

// StoreSchedule commits a schedule to the database
func StoreSchedule(db *bolt.DB, schedule *Schedule) error {
	// date := time.Date(schedule.OnTime)
	// db.Update(func(tx *bolt.Tx) error {
	// 	var err error
	// 	b := tx.Bucket([]byte("Schedule"))
	// 	err = b.Put([]byte(date.Format(time.Unix)), []byte(""))
	// 	return err
	// })

	return nil
}
