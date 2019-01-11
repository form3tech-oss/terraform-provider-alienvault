# AlienVault Terraform Provider

[![Build Status](https://travis-ci.org/form3tech-oss/terraform-provider-alienvault.svg?branch=master)](https://travis-ci.org/form3tech-oss/terraform-provider-alienvault)

Terraform Provider for [AlienVault USM Anywhere](https://www.alienvault.com/products/usm-anywhere).

## Summary

This provider enables you to setup AlienVault sensors and optionally assign/schedule custom log collection jobs to them.

## Installation

Download the relevant binary from [releases](https://github.com/form3tech-oss/terraform-provider-alienvault/releases) and copy it to `$HOME/.terraform.d/plugins/`.

## Authentication

There are two methods for authenticating with this provider.

You can specify credentials and your AlienVault fully-qualified domain name in the `provider` block, as below:

```hcl
provider "alienvault" {
    fqdn     = "mycompany.alienvault.cloud"
    username = "user@email.com"
    password = "..."
}
```

Alternatively, you can use the following environment variables:

```bash
export ALIENVAULT_FQDN="mycompany.alienvault.cloud"
export ALIENVAULT_USERNAME="user@email.com"
export ALIENVAULT_PASSWORD="..."
```

You can mix both methods, but provider block variables will override environment variables, where provided.

Note that the provider does not take client/secret credentials. This is because the provider currently makes use of an internal API, as the public v2 API does not yet support sensors or jobs.

## Resources

### `alienvault_sensor`

There are two parts to a sensor: the *appliance*, which is a VM running the AV agent software; and the *registration* which is what this resource models.

In order to create this resource, you are required to already have an appliance running. The creation of this resource will ensure the appliance is in the correct state before creating a sensor registration for it.

#### Fields

- `name` The name of the sensor, such as "my-production-sensor".
- `description` A description of the sensor. If not provided, this will default to "Created by terraform".
- `ip` The public IP address of the associated sensor appliance.

### `alienvault_job_aws_bucket`

A job for retrieving log files from an AWS bucket.

#### Fields

- `name` The name of the job, such as "route53-log-collection".
- `description` (Optional) A description of the job.
- `sensor` The ID of the sensor where this job should run. This is actually known as `UUID` in the AV environment, though confusingly there is also an `ID` field in their API which is not used. You will probably want to use something like `${alienvault_job_aws_bucket.whatever.id}`
- `schedule` The schedule of when to run this job. Can currently be set to "daily" or "hourly".
- `disabled` (Optional) This can be set to "True" to temporarily prevent the job from running.
- `source_format` (Optional) This should be set to "raw" for raw logs or "syslog" for SysLog files. Defaults to "raw".
- `plugin` The plugin to use to parse the log files. A list of available plugins can be found further down this page.
- `bucket` The name of the bucket to use to look for log files.
- `path` (Optional) The path within the specified bucket where log files can be found.

### `alienvault_job_aws_cloudwatch`

A job for retrieving log files from AWS CloudWatch streams.

#### Fields

- `name` The name of the job, such as "route53-log-collection".
- `description` (Optional) A description of the job.
- `sensor` The ID of the sensor where this job should run. This is actually known as `UUID` in the AV environment, though confusingly there is also an `ID` field in their API which is not used. You will probably want to use something like `${alienvault_job_aws_bucket.whatever.id}`
- `schedule` The schedule of when to run this job. Can currently be set to "daily" or "hourly".
- `disabled` (Optional) This can be set to "True" to temporarily prevent the job from running.
- `source_format` (Optional) This should be set to "raw" for raw logs or "syslog" for SysLog files. Defaults to "raw".
- `plugin` The plugin to use to parse the log files. A list of available plugins can be found further down this page.
- `region` The AWS region where the CloudWatch data is available.
- `group` (Optional) The CloudWatch group name. Defaults to "*", meaning all.
- `stream` (Optional) The CloudWatch stream name. Defaults to "*", meaning all.

## Available Plugins

These must be specified in the `plugin` field exactly as they appear in the list below.

- AIX Audit
- AWS API Gateway
- AWS Application Load Balancer
- AWS Web Application Firewall
- AWSWindows
- AlienVault Agent
- AlienVault Agent - Windows EventLog
- AlienVault NIDS
- Amazon AWS CloudTrail
- Amazon GuardDuty
- Amazon Macie
- Amazon Redshift
- Amazon Redshift User Activity
- Apache
- Apple Airport Extreme
- Arbor Networks Pravail APS
- Arpwatch
- ArticaProxy
- Aruba
- Aruba ClearPass CEF
- Aruba Clearpass
- Asterisk VoIP
- Aunt Bertha Website Activity Plugin
- Auth0
- Avaya Media Gateway
- Avaya VSP Switches
- Avaya Wireless LAN
- Azure IIS
- Azure Insight
- Azure Multifactor Authentication
- Azure SQL Server
- Azure Security Center
- Azure Web App
- Azure Windows Events
- Barracuda NextGen Firewall
- Barracuda NextGen Firewall Traffic
- Barracuda Spam Firewall
- Barracuda Web Application Firewall
- Barracuda Web Application Firewall CEF
- Barracuda Web Filter
- Bitdefender GravityZone
- Bluecoat W3C
- Box Events
- Brocade
- Buffalo TeraStation
- Business Intelligence Analytics
- Centrify Server Suite
- CheckPoint FW1
- CheckPoint FW1 Generic
- CheckPoint FW1 Loggrabber
- CheckPoint FW1 R.80 CEF
- CheckPoint FW1 R77.3O
- Cisco ACE
- Cisco ACS
- Cisco AMP for Endpoints
- Cisco ASA
- Cisco ASR
- Cisco ESA
- Cisco Firepower NGIPS
- Cisco Firepower NGWF
- Cisco ISE
- Cisco Ironport
- Cisco Lancope StealthWatch
- Cisco Meraki
- Cisco Nexus
- Cisco Pix
- Cisco Router
- Cisco Umbrella
- Cisco VPN
- Cisco WLC
- Citrix NetScaler
- Citrix Netscaler Application Firewall CEF plugin
- Clavister Firewall
- CloudFront RTMP distribution W3C
- CloudFront Web distribution W3C
- CloudPassage CEF
- Cloudflare Enterprise Log Share
- Cloudflare Enterprise Log Share Received
- CrowdStrike Falcon
- CyberArk Enterprise Password Vault
- Cylance CylancePROTECT
- Cyphort CEF plugin
- D-Link UTM Firewall
- DELL Compellent SC
- DELL IDRAC
- Darktrace Cyber Intelligence Platform
- Darktrace Cyber Intelligence Platform JSON
- Dell Networking X-Series
- Dell SecureWorks
- Dell SonicWall UTM
- DenyAll WAF
- Digital Defense Incorporated Frontline Vulnerability Manager
- Docker
- DrayTek Vigor
- Dropbox
- Dtex
- Duo Two-Factor Authentication CEF
- ELBAccess
- Endpoint Protector
- Eset
- Extreme Networks SummitX and Black Diamond Switches
- F5 BIG-IP ASM
- F5 BIG-IP Access Policy Manager
- F5 Big-ip
- Fail2ban
- FireEye Central Management System
- FireEye Endpoint Security HX Series
- FireEye Malware Protection Systems
- Forcepoint Triton AP-Web
- ForeScout NAC
- Fortinet FortiClient
- Fortinet FortiNAC
- Fortinet FortiWAN
- Fortinet FortiWeb
- Fortinet Fortigate
- FreeRadius
- G Suite Audit
- G Suite Drive
- Google Cloud Platform - Compute Engine
- Google Cloud Platform Audit
- H3C Switch
- HAProxy
- HP Storage Area Network Switch
- HP Switch
- HPE Integrated Lights Out
- HPE StoreOnce
- Huawei NGFW
- IBM Maximo
- IBM QRadar Network Security
- IBM Tivoli Access Manager WebSEAL
- Imperva SecureSphere
- Imperva SecureSphere CEF
- Incapsula CEF plugin
- Infoblox DDI
- Ipswitch WS_FTP
- Jenkins
- JumpCloudAPI
- Juniper EX Series
- Juniper NetScreen ScreenOS
- Juniper NetScreen ScreenOS Traffic
- Juniper Network Security Manager
- Juniper SRX Junos
- Juniper Secure Access VPN
- Kaspersky Security Center
- Kaspersky Security Center CEF
- Kerio Connect
- Linux Auditd
- Linux BIND
- Linux CRON
- Linux ClamAV
- Linux DHCP client
- Linux DHCPD
- Linux DNSMASQ
- Linux IPTables
- Linux SSH
- Linux SUDO
- Linux Systemd
- Linux Useradd/Groupadd
- Malwarebytes Breach Remediation
- Malwarebytes Endpoint Protection
- Malwarebytes Endpoint Security
- Malwarebytes Management Console
- ManageEngine ADAudit Plus
- McAfee Database Security
- McAfee EPO
- McAfee Network Security Platform
- McAfee Web Gateway
- Microsoft Advanced Threat Analytics
- Microsoft IIS 8.0+ Plugin
- Microsoft IIS pre-8.0 Plugin
- MikroTik Router
- Mimecast
- MySQL Community Edition
- NetApp Hybrid-Flash Storage System
- Netgate
- Netgear Switch
- Nginx
- Nginx Error
- Nginx NAXSI
- Nimble Storage
- OSSEC JSON
- OSSEC v2.5
- ObserveIT
- Office 365 Audit
- Office 365 Azure AD
- Office 365 Exchange
- Office 365 SharePoint
- Okta
- OpenVPN Syslog
- Oracle Audit Syslog
- Osquery
- PA File Sight
- PacketFence
- Palo Alto Traps
- Palo Alto Traps Management Service
- Paloalto PAN-OS
- Paloalto PAN-OS CEF
- Passwordstate
- Percona Audit Log
- Plixer Scrutinizer
- Postfix
- PostgreSQL
- PowerDNS
- ProFTPD
- Proofpoint Targeted Attack Protection
- Pulse Connect Secure
- Pure-FTPd
- RSA Authentication Manager
- Radware Cloud Services
- Riverbed STM
- Route 53 DNS Queries
- Ruckus SmartCell Gateway
- Ruckus Virtual SmartZone
- Ruckus Wireless ZoneDirector
- STEALTHbits File Activity Monitor
- Salesforce Activity
- Samba
- Sangfor Next-Generation Firewall
- SecureAuth
- SendMail
- SentinelOne
- ServerAccess
- Shrubbery Tacacs
- Silver Peak WAN Optimization
- Smoothwall Express
- Snort Syslog
- SoftEther VPN
- Sophos Central
- Sophos Central JSON
- Sophos Cyberoam
- Sophos Enterprise Console
- Sophos UTM
- Sophos UTM WAF
- Sophos Web Security
- Sophos XG
- SourceFire IDS
- SpyCloud
- Squid
- StrongSwan VPN
- Symantec ATP
- Symantec DLP
- Symantec EPM
- Symantec Encryption
- Syncplify.me
- Synology NAS
- Tesserent Next Gen Firewall
- Trend Micro Control Manager
- Trend Micro Control Manager CEF
- Trend Micro Deep Discovery Inspector
- Trend Micro Deep Security
- Trend Micro TippingPoint
- Trend Micro Vulnerability Protection
- Trustwave Secure Web Gateway
- Trustwave Secure Web Gateway Traffic
- UFW
- Ubiquiti EdgeRouter
- Ubiquiti Unifi
- VMRay Analyzer
- VMware AirWatch
- VMware Esxi
- VMware NSX
- VMware SSO
- VMware vCenter
- VMware vRealize
- VMware vShield
- VMwareAPI
- VPC Flow Logs
- Varonis DatAdvantage
- Vectra
- Versa Director
- Versa FlexVNF
- Virtual LoadMaster
- Vormetric Data Security Manager
- Watchguard Firebox
- Watchguard XTM
- Wazuh
- Webmin
- Webroot FlowScape
- Websense Email Security Gateway
- Websense Web Security Gateway
- Windows DHCP NxLog
- Windows DNS Server
- Windows Exchange NxLog
- Windows Firewall NxLog
- Windows IIS NxLog
- Windows NxLog
- Windows SQL NxLog
- Windows Snare
- Windows Winlogbeat
- ZeroFOX
- ZingBox IoT Guardian
- ZyXEL ZyWALL
- cb Defense
- cb Protection
- cb Response
- cb Response JSON
- pfSense Filter
- pfSense System
- zScaler NSS

## Example Usage

```hcl

resource "alienvault_sensor" "main" {
  count   = "${var.provision_alienvault_sensor ? 1 : 0}"
  ip      = "${data.aws_instance.alienvault_sensor.public_ip}"
  name    = "${var.stack_name}-sensor"
}

resource "alienvault_job_aws_cloudwatch" "route53_siem" {
  count         = "${var.provision_alienvault_sensor ? 1 : 0}"
  name          = "Route53 log collection"
  sensor        = "${alienvault_sensor.main.id}"
  schedule      = "hourly"
  region        = "us-east-1"
  group         = "${aws_cloudwatch_log_group.aws_route53_stack.name}"
  source_format = "raw"
  plugin        = "Route 53 DNS Queries"
}
```