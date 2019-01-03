package alienvault

// JobApplication is the application associated with the job. Currently we support alienvault.JobApplicationAWS, which is Amazon AWS
type JobApplication string

const (
	// JobApplicationAWS Amazon AWS
	JobApplicationAWS JobApplication = "amazon-aws"
)

// JobAction is the action to take when running this job, such as checking a bucket for log files (alienvault.JobActionMonitorBucket)
type JobAction string

const (
	// JobActionMonitorBucket is the action of monitoring an S3 bucket for log files
	JobActionMonitorBucket JobAction = "s3TrackFiles"
	// JobActionMonitorCloudWatch is the action of monitoring cloudwatch for log files
	JobActionMonitorCloudWatch JobAction = "cloudWatchTrackFiles"
)

// JobType is the type of job, such as alienvault.JobTypeCollection for collecting log files
type JobType string

const (
	// JobTypeCollection is a job type which collects log files from a given source
	JobTypeCollection JobType = "collection"
)

// JobSourceFormat is the format which the log files are in - alienvault.JobSourceFormatRaw or alienvault.JobSourceFormatSyslog
type JobSourceFormat string

const (
	// JobSourceFormatRaw describes raw log files
	JobSourceFormatRaw JobSourceFormat = "raw"
	// JobSourceFormatSyslog describes log files in syslog format
	JobSourceFormatSyslog JobSourceFormat = "syslog"
)

// JobSchedule is a cron-like syntax which describes when to run the scheduled job. Constants are available to simplify this, such as alienvault.JobScheduleHourly
type JobSchedule string

const (
	// JobScheduleHourly will run every hour at :02
	JobScheduleHourly JobSchedule = "0 2 0/1 1/1 * ? *"

	// JobScheduleDaily will run daily at 00:02
	JobScheduleDaily JobSchedule = "0 2 0 1/1 * ? *"
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
	Plugin       string          `json:"plugin,omitempty"` // Plugin describes the plugin used to parse the log files e.g. "PostgreSQL" for postgres logs
	SourceFormat JobSourceFormat `json:"source"`           // SourceFormat is essentially alienvault.JobSourceFormatRaw or alienvault.JobSourceFormatSyslog
}
