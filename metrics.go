/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"time"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Metric struct {
	Rate  MetricRate
	Type  MetricType
	Name  string
	Mean  float64 // Mean value per hour
	Total uint    // Total over the past hour
}

type (
	MetricRate uint
	MetricType uint
)

/////////////////////////////////////////////////////////////////////
// INTERFACE

// Metrics returns various metrics for host and
// custom metrics
type Metrics interface {
	Driver

	// Uptimes for host and for application
	UptimeHost() time.Duration
	UptimeApp() time.Duration

	// Load Average (1, 5 and 15 minutes)
	LoadAverage() (float64, float64, float64)

	// Return counter channel, which when you send a value on
	// it will increment a counter
	//NewCounter(MetricType, MetricRate, string) (chan<- uint, error)

	// Return all metrics of a particular type, or METRIC_TYPE_NONE
	// for all metrics
	//Metrics(MetricType) []*Metric
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	METRIC_RATE_NONE MetricRate = iota
	METRIC_RATE_SECOND
	METRIC_RATE_MINUTE
	METRIC_RATE_HOUR
	METRIC_RATE_DAY
)

const (
	METRIC_TYPE_NONE MetricType = iota
)
