# AlienVault

[![Build Status](https://travis-ci.org/form3tech-oss/alienvault.svg?branch=master)](https://travis-ci.org/form3tech-oss/alienvault)

A basic Go package providing a client for the AlienVault API.

Whilst AV do provide a public API, this does not yet support operations on job scheduling and sensors. For this reason, this client utilises an unoffical internal API used by the AV web UI to get the job done. The plan is to move this to the public API as soon as support for the required data types is made available.

## Example Usage

```go
alienVaultClient := alienvault.New(
    os.GetEnv("AV_FQDN"),
    alienvault.Credentials{
        Username: os.GetEnv("AV_USERNAME"),
        Password: os.GetEnv("AV_PASSWORD"),
    })

if err := alienVaultClient.Authenticate(); err != nil {
    panic(err)
}

job, err := alienVaultClient.GetJob("...")
if err != nil {
    panic(err)
}

fmt.Printf("Job details: %#v\n", *job)
```