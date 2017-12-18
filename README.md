# Veracode Auto Mitigate

## Description
Command line app that mitigates flaws in Veracode based on CWE, scan type, and specific text in the description.

## Parameters
`-config`: path to JSON config file

## Configuration File
A sample config file is below:
```
{
    "auth": {
      "credsFile": "/Users/bpitta/.veracode/credentials"
    },
    "scope": {
        "allApps": false,
        "appListTextFile": "/Users/bpitta/go/src/github.com/brian1917/vcodeAutoMitigate/applist.txt"
    },
    "mode":{
      "logOnly": false,
      "proposeOnly": false,
      "proposeAndAccept": true
    },
    "targetFlaws": {
      "cweList": "80",
      "requireTextInDesc": true,
      "requiredText":"realName",
      "static": true,
      "dynamic": true
    },
    "mitigationInfo":{
        "mitigationType": "netenv",
        "proposalComment": "Proposal text",
        "approvalComment": "Approval text"
    }
}
 ```
 **_Configuration Notes_**:
 1. The `appsListTextFile` parameter should be used when `allApps` is set to `false`. It should a text file with target app IDs on separate lines.
 2. The `cweList` parameters should be a comma separated list of CWEs that should be targeted.
 3. The `requiredTextInDesc` parameter will search for text in the flaw description. The text to search for should be placed in the `requiredText` parameter. For example, you can use this to target flaws on a specific cookie from a dynamic scan by including the cookie name.
 4. The `mitigationType` must be `comment`, `fp`, `appdesign`, `osenv`, or `netenv`.

## Credentials File
The credentials file should be set up as follows:
```
veracode_api_key_id = ID HERE
veracode_api_key_secret = SECRET HERE
```

## Executables
The executable is available in the release section of the repository: https://github.com/brian1917/vcodeAutoMitigate/releases

## Third-party Packages
github.com/brian1917/vcodeapi (https://godoc.org/github.com/brian1917/vcodeapi)