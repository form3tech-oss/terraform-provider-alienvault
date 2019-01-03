package alienvault

import (
	"fmt"

	"github.com/form3tech-oss/alienvault"
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
	if v != string(alienvault.JobSourceFormatRaw) && v != string(alienvault.JobSourceFormatSyslog) {
		errs = append(errs, fmt.Errorf("%q must be either %q or %q, got: %s", key, alienvault.JobSourceFormatRaw, alienvault.JobSourceFormatSyslog, v))
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
