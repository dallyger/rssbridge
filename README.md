# RSS-Bridge

Consume sites through your RSS reader.

Supported websites are scraped and turned into an RSS/Atom feed each time a
HTTP request hits the server.

## Usage

Executing the program will start a web server on port 3000.
Feeds can be crawled and created by calling their endpoint.
Available endpoints can be seen at [example](./example/).

Those files are plain text files describing the HTTP request to send.
Take a look, it's really straight forward to understand.

Or run them directly in the terminal using [Hurl](https://hurl.dev). Example:

```sh
hurl --variable base_url=http://localhost:3000 ./example/store.shopware.com-plugin-changelog.hurl
```

## Setup

Requirements: `go`,`make`, `docker` (optional)

Clone the project:

```sh
git clone https://github.com/dallyger/rssbridge
cd rssbridge
```

Compile binaries and run it:

```
make build && bin/rssbridge
```

Or use Docker:

```
make do-build do-run
```

Tested on Linux. Windows is not supported. Use docker or WSL or something.

## Privacy concerns

The originating IP addresses and host header are forwarded to the scraped site.

This is by design as the service relies on those third-party sites for the
content. Therefore, it is only fair to be nice and forward as much information
as possible.

If a consumer spams a site through this, they may just block the originating
IPs instead of the whole bridge service.

