package domain

import "time"

type Plan struct {
	ID       int
	Name     string
	Duration time.Duration
	Price int
	IsActive bool
}

type PlanRepository interface {
	
} 