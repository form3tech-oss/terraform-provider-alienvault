package alienvault

type JobApplication string

const (
	JobApplicationAWS JobApplication = "amazon-aws"
)

type JobAction string

const (
	JobActionMonitorBucket     JobAction = "s3TrackFiles"
	JobActionMonitorCloudWatch JobAction = "cloudWatchTrackFiles"
)

type JobType string

const (
	JobTypeCollection JobType = "collection"
)

type JobSourceFormat string

const (
	JobSourceFormatRaw    JobSourceFormat = "raw"
	JobSourceFormatSyslog JobSourceFormat = "syslog"
)

type JobSchedule string

const (
	JobScheduleHourly JobSchedule = "0 0 0/1 1/1 * ? *"
)

type job struct {
	UUID        string         `json:"uuid,omitempty"` // UUID is a unique ID for the job. Read-only.
	SensorID    string         `json:"sensor"`         // SensorID is the ID of the sensor to use to run this job.
	Schedule    JobSchedule    `json:"schedule"`       // Schedule is a slightly obscure cron format, such as "0 0 0/1 1/1 * ? *" meaning hourly
	Name        string         `json:"name"`           // Name is a human-readable name for the job
	Description string         `json:"description"`    // Description is a human-readable description of the job
	Disabled    bool           `json:"disabled"`       // Disabled describes whether the job should run or not. You can set this if you wish to temporarily disable the job.
	App         JobApplication `json:"app"`            // App describes the app associated with this job e.g. "amazon-aws". You do not usually need to populate this, it will be filled by default.
	Action      JobAction      `json:"action"`         // Action describes the action associated with this job e.g. "s3TrackFiles". You do not usually need to populate this, it will be filled by default.
	Type        JobType        `json:"type"`           // Type describes the type of job e.g. "collection" for log collection jobs. You do not usually need to populate this, it will be filled by default.
	Custom      bool           `json:"custom"`         // Custom describes whether the job was built in or a custom job created by the user. Read-only.
}

type jobParams struct {
	Plugin       string          `json:"plugin,omitempty"` // Plugin describes the plugin used to aprse the log files e.g. "PostgreSQL" for postgres logs
	SourceFormat JobSourceFormat `json:"source"`           // SourceFormat is essentially alienvault.JobSourceFormatRaw or alienvault.JobSourceFormatSyslog
}
