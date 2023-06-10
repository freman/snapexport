# Snapexport - Export your chat logs from [SnapEngage](https://snapengage.com/)

Whether leaving [SnapEngage](https://snapengage.com/) or simply archiving for prosperity this tool will let you
export all of your conversations from [SnapEngage](https://snapengage.com/) as JSON filess in a directory

## Requirements

Snapengage restricts API access to the account owner, they will need to log in to get your Organisation ID amd
get/create an API token

## Arguments

### -end YYYY-MM-DD [env: SNAPENGAGE_ENV]

Stop exporting at date YYYY-MM-DD - Defaults to TODAY if unspecified.

### -start YYYY-MM_DD [env: SNAPENGAGE_START]

Start exporting from what date YYYY-MM-DD.

### -org 0000000000000000 [env: SNAPENGAGE_ORG]

Snapengage Organisation ID.

## -token 0123456789abcdef01234567890bcdef [emv: SNAPENGAGE_TOKEN]

Snapengage API Token

## -widget 12345678-9012-3456-789a-bcdef012345 [env: SNAPENGAGE_WIDGET]

Snapengage Widget ID 

## -output /path/to/save/to [env: SNAPENGAGE_OUTPUT]

Directory to write cases to
