# netmeg : CLI for sending HTTP requests implemented in Go

## Overview [![GoDoc](https://godoc.org/github.com/brodyzera/netmeg?status.svg)](https://godoc.org/github.com/brodyzera/netmeg)

## Install
```
go get github.com/brodyzera/netmeg
```

## Examples
`netmeg request -u https://google.com -m get -n 2`

`netmeg request -u https://jsonplaceholder.typicode.com/posts -m POST -n 10 -b "{\"title\": \"This is a title.\"}" -H "Content-Type:application/json, Test:123" -o post-response.log --mode file`

**Note the double quotes around the in-line -b and -H values, as well as the backslash-escaped double quotes within the JSON body.  If you are passing in a Request Body or Headers, use pre-made files instead.**

`netmeg request -u https://jsonplaceholder.typicode.com/posts -m POST -n 10 --bfile ./resources/body.json --hbody ./resources/headers.txt -o ./output/post-response.log --mode both`

You can also load pre-built requests from JSON files.  For example, create a JSON Array of JSON objects, like so;
### requests.json
```json
[
    {
        "description": "Simple Google GET x 5",
        "url": "https://google.com",
        "method": "get",
        "amount": 5,
        "body": "",
        "headers": "",
        "bfile": "",
        "hfile": "",
        "output": "",
        "mode": "file"
    },
    {
        "description": "JSON Placeholder test site POST x 10",
        "url": "https://jsonplaceholder.typicode.com/posts",
        "method": "POST",
        "amount": 10,
        "body": "{\"title\": \"This is a title.\"}",
        "headers": "Content-Type:application/json",
        "bfile": "",
        "hfile": "",
        "output": "post-response.log",
        "mode": "both"
    }
]
```
With the above JSON saved in requests.json, we can now load it using

`netmeg load -f requests.json`

The requests saved in requests.json will execute sequentially  (although the specified amount per request will still run asynchronously).

### Sample Workspace Setup
![Workspace in VS Code](https://user-images.githubusercontent.com/4110514/69850289-2ff2cb00-123c-11ea-8e31-44287b3d2fa2.png)

## License
Apache 2.0.