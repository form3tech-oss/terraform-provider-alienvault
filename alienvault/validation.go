package alienvault

import (
	"fmt"

	"github.com/google/uuid"
)

func validateJobPlugin(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	for _, plugin := range plugins {
		if plugin == v {
			return
		}
	}

	errs = append(errs, fmt.Errorf("%q must be a supported plugin, '%s' is not supported", key, v))
	return
}

func validateJobSourceFormat(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if v != JobSourceFormatRaw && v != JobSourceFormatSyslog {
		errs = append(errs, fmt.Errorf("%q must be either %q or %q, got: %s", key, JobSourceFormatRaw, JobSourceFormatSyslog, v))
	}
	return
}

func validateJobSensor(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if _, err := uuid.Parse(v); err != nil {
		errs = append(errs, fmt.Errorf("%q must be a valid UUID, got: %s", key, v))
	}
	return
}
