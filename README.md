# netmeg : CLI for Sending Concurrent HTTP/S Requests
![](https://github.com/brodyzera/netmeg/workflows/go-test/badge.svg)

## Installation
```
go get github.com/brodyzera/netmeg
```
Running the command above will install netmeg to your GOPATH (assuming that you have Go installed).

If desired, you can also download the compiled binaries from the [releases page](https://github.com/Brodyzera/netmeg/releases).
Compiled binaries are supplied for both Linux and Windows.

## Usage
At this point in time, there are two sub-commands for netmeg; **request** and **load**.
### netmeg request
```
Flags:
  -n, --amount int       Amount of requests to send (default 1)
      --bfile string     File containing Request body (overrides --body and -b flags)
  -b, --body string      Request body
  -H, --headers string   Header list formated as {key}:{value}, separated by commas
  -h, --help             help for request
      --hfile string     File containing Headers (overrides --headers and -H flags)
  -m, --method string    HTTP method to use for the request (default "get")
      --mode string      Output mode for result (console, file, or both) (default "console")
  -o, --output string    Path to file for results (default "results-{timestamp}.json")
  -u, --url string       URL to send the request to
```
The only "required" flag is `--url`, meaning that if you simply want to send a GET request to "https://localhost", your command would look like this;

```
netmeg request --url https://localhost
```
This works since `--amount` defaults to **1**, and `--method` defaults to **GET**.  Note that `--mode` defaults to **console**,
meaning that the command's results will be written to your standard console.  Valid values for `--mode` are **console**,
**file**, or **both**.

To send multiple requests (in parallel), use the `--amount` flag.

```
netmeg request -u https://localhost --amount 5
```
The command above will send 5 **GET** requests to **https://localhost** in parallel.

Lets try a more complex command;

```
netmeg request -u https://jsonplaceholder.typicode.com/posts -m POST -n 10 -b "{\"title\": \"This is a title.\"}" -H "Content-Type:application/json, Test:123" -o post-response.log --mode file
```
This command will send 10 **POST** requests to **https://jsonplaceholder.typicode.com/posts** with the specified **Headers**
and **Post Body**.  The `--mode` and `-o` flags ensure that the results are written to the file **post-response.log** in
our current working directory, rather than printing the results to our standard console.  If you don't supply a filename
when using **file** mode, a file named **results-{timestamp}.json** will be created in your current working directory.

If you include **{timestamp}** in the filename for the `--output` flag, like this;

```
netmeg request -u https://localhost --amount 5 --mode file -o get-results-{timestamp}.json
```
your system's current date and time (year|month|day|hour|minute|second) will be injected in to request result filename.

Note the double quotes around the in-line `-b` and `-H` values, as well as the backslash-escaped double quotes within
the **JSON Post Body**.  If you are passing in a Request Body or Headers, use pre-made files instead.

```
netmeg request -u https://jsonplaceholder.typicode.com/posts -m POST -n 10 --bfile ./resources/body.json --hfile ./resources/headers.txt -o ./output/post-response.log --mode both
```
In our current working directory, we have a folder named **resources** which contains the files **body.json** and **headers.txt**.
These files are referenced using the `--hfile` and `--bfile` flags, rather than the `-b` and `-H` flags.  Also note that
the `--mode` is set to **both**, meaning that the request results will be written to both a file as well as the console.

### netmeg load
You can also load pre-built requests from JSON files.  For example, create a JSON Array of JSON objects, like so;
#### requests.json
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

### Results and Output
Output from both the **file** and **console** modes will be identical.

#### Example Output
```
new_request
Status_Code: 200
Response_Time: 0.339s
Body: This is a test response body

new_request
Status_Code: 200
Response_Time: 0.98s
Body: This is another test response body
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
Apache 2.0.
