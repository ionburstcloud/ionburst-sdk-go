![Gitlab pipeline status](https://img.shields.io/gitlab/pipeline/ionburst/ionburst-sdk-net/main?color=fb6a26&style=flat-square)
[![slack](https://img.shields.io/badge/Slack-4A154B?style=flat-square&logo=slack&logoColor=white)](https://join.slack.com/t/ionburst-cloud/shared_invite/zt-panjkslf-Z5DOpU1OOeNPkXgklD~Cpg)

# Ionburst SDK for Go

The **Ionburst SDK for Go** enables Golang developers to easily work with [Ionburst Cloud][ionburst] and build ultra-secure and private storage into their applications.

* [API Docs][docs-api]
* [SDK Docs][sdk-website]
* [Issues][sdk-issues]
* [SDK Samples](https://docs.ionburst.io/#/sdk?id=usage)

## Getting Started

### Installation

```sh
$ go get gitlab.com/ionburst/ionburst-sdk-go
```

### Configuration

The Ionburst SDK can get its configuration (ionburst_id, ionburst_key, ionburst_uri) from the following three files.

If `ionburst_id` and `ionburst_key` are not specified by environment variable, they are obtained from the credentials file.

If `ionburst_uri` is not specified in the Ionburst constructor, it'll check the credentials file.

#### Environment Variables

```sh
IONBURST_ID=IB******************
IONBURST_KEY=eW91aGF2ZXRvb211Y2h0aW1lb255b3VyaGFuZHMh
IONBURST_URI=https://api.example.ionburst.cloud/
```

#### Credentials file

`~/.ionburst/credentials` on Mac, Linux and BSD, and `C:\Users\%USERNAME%\.ionburst\credentials` on Windows

```sh
[example]
ionburst_id=IB******************
ionburst_key=eW91aGF2ZXRvb211Y2h0aW1lb255b3VyaGFuZHMh
ionburst_uri=https://api.example.ionburst.cloud/
```

### ioncli

[ioncli](ioncli) is a command line utility written using this SDK to perform basic operations.

Please [click here](ioncli) to learn more.

### Usage

#### Initialise

```go
package client

import (
    ionburst "gitlab.com/ionburst/ionburst-sdk-go"
}

func main() {
    //Create a new client using the default config
    client, err := ionburst.NewClient() 

    //Create a client with implicit path and credentials profile name
    client, err := ionburst.NewClientPathAndProfile(configFilePath, credentialsProfileName, setDebugMode [true/false])

}
```

#### Upload Data

```go
//get a readable stream
ioReader, _ := os.Open(FilePath)

client.Put(FileID, ioReader, classification)
//if classification is an empty string ("") it wont be passed

//Upload from a filepath instead
client.PutFromFile(FileID, FilePath, classification)
```

#### Download Data

```go
ioReader, err := client.Get(FileID)

//Download to a filepath instead
err := client.GetToFile(FileID, OutputFilePath)

//Download and output the size of the downloaded content
ioReader, sizeOfContent, err := client.GetWithLen(FileID)
```

#### Delete Data

```go
err := client.Delete(FileID)
```

#### Get Classifcations

```go
classifications, _ := client.GetClassifications()
```

### Usage in Deferred Mode

#### Upload Data Deferred

```go
token, err := cli.PutDeferred(name, r, "")
if err != nil {
    t.Error(err)
    return
}
```

#### Download Data Deferred

```go
token, err := cli.GetDeferred(name)
if err != nil {
    t.Error(err)
    return
}
```

#### Check Data Deferred

```go
res, err = cli.CheckDeferred(tk)
if err != nil {
    t.Error(err)
    return
} else if !res.Success {
    t.Error(fmt.Sprintf("ERR: %s - %d", res.Message, res.Status))
    return
} else {
    ...
}
```

#### Fetch Data Deferred

```go
ioReader, err := cli.FetchDeferred(tk)
if err != nil {
    t.Error(err)
    return
}
```

## Getting Help

Please use the following community resources to get help. We use [Gitlab issues][sdk-issues] to track bugs and feature requests.

* Join the Ionburst Cloud community on [Slack](https://join.slack.com/t/ionburst-cloud/shared_invite/zt-panjkslf-Z5DOpU1OOeNPkXgklD~Cpg)
* Get in touch with [Ionburst Support](https://docs.ionburst.io/#/introduction?id=contact-amp-support)
* If you have found a bug, please open an [issue][sdk-issues]

### Opening Issues

If you find a bug, or have an issue with the Ionburst SDK for Go we would like to hear about it. Check the existing [issues][sdk-issues] and try to make sure your problem doesn’t already exist before opening a new issue. It’s helpful if you include the version of `ionburst-sdk-go` and the OS you’re using. Please include a stack trace and steps to reproduce the issue.

The [Gitlab issues][sdk-issues] are intended for bug reports and feature requests. For help and questions with using the Ionburst SDK for Go please make use of the resources listed in the Getting Help section. There are limited resources available for handling issues and by keeping the list of open issues clean we can respond in a timely manner.

## SDK Change Log

The changelog for the SDK can be found in the [CHANGELOG file.](CHANGELOG.md)

## Contributors

A massive thanks to [Craig Smith](https://github.com/spuddleziz) for developing this SDK.

## Dependencies

* [github.com/antihax/optional](https://github.com/antihax/optional) - a utility library for the handling of default and optional values
* [golang.org/x/oauth2](https://golang.org/x/oauth2) - Google's OAuth library

[ionburst]: https://ionburst.cloud
[sdk-website]: https://ionburst.cloud/docs/sdk
[sdk-source]: https://gitlab.com/ionburst/ionburst-sdk-go
[sdk-issues]: https://gitlab.com/ionburst/ionburst-sdk-go/issues
[sdk-license]: https://gitlab.com/ionburst/ionburst-sdk-go/-/blob/master/LICENSE
[docs-api]: https://ionburst.cloud/docs/api
[ioncli]: https://gitlab.com/ionburst/ionburst-sdk-go/-/tree/master/ioncli