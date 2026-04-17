package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type PeriodicType string

const (
	NotPeriodic  PeriodicType = "not_periodic"
	Daily        PeriodicType = "daily"
	Monthly      PeriodicType = "monthly"
	ExactDates   PeriodicType = "exact_dates"
	EvenOddDates PeriodicType = "even_odd_dates"
)

type EvenOddType string

const (
	EvenDays EvenOddType = "even"
	OddDays  EvenOddType = "odd"
)

type Recurrence struct {
	Type        PeriodicType `json:"periodic_type"`
	EveryNDays  int          `json:"every_n_days,omitempty"`
	MonthlyDate int          `json:"monthly_date,omitempty"`
	Dates       []time.Time  `json:"dates,omitempty"`
	EvenOdd     EvenOddType  `json:"even_odd,omitempty"`
}

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Recurrence  Recurrence `json:"recurrence"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (p PeriodicType) Valid() bool {
	switch p {
	case NotPeriodic, Daily, Monthly, ExactDates, EvenOddDates:
		return true
	default:
		return false
	}
}

func (e EvenOddType) Valid() bool {
	switch e {
	case EvenDays, OddDays:
		return true
	default:
		return false
	}
}
