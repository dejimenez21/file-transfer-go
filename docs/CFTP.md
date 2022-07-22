# Custom File Transfer Protocol

The Custom File Transfer Protocol (**CFTP**) is a TCP-based application layer protocol for file transfer between multiple clients. This protocol is designed for a client - server architecture. It works as support for a publisher - subscriber pattern, where clients can send and receive files through channels.

## Messages

Messaging consists of Requests and Responses, both been used by both the client and the server. Each request should get a response indicating if the request was successful. Timeouts periods are determined by each system in its configuration.

### Requests

Requests are structured as follows:  

1. **[Method](#method)** that indicates the action to perform.
2. **Meta** section, that is a json string containing data about the request.
3. **Channels** section, a comma-separated list of the channels involved in the request.
4. **File** bytes.

The only fully required section is the **Method**, all the others may or not be specified depending on the operation to perform. For example, a client who wants to suscribe to a channel doesn't need to include the **File** section on the request.

Every section described above should be separated by a line break (`"/n"`). Even if a section is not to be included the line break should.
 For example:

```ps1

SEND

chn1,chn2
[FILE_BYTES]
 ```

In the above example, the **Meta** section is not included, but the empty line where it wolud be is included.

>(Each section is described more explicitly in following parts of this documentation.)

### Responses

Every request should be responded indicating if said request was successful. The response should be a short plain text string lowercased and without spaces. Responses should always notify an expected behavior. Exceptions will be indetified by the absence of responses within the timeout period. The available responses are:

* `ok`: it indicates success.

## Requests Sections

### Method

The Method represents the action to perform and is the first section of the request. Methods are upppercase words that represent verbs of the english language. Clients and servers use different methods.

Client's methods:

* `SEND`
* `SUSCRIBE`

Server's methods:

* `DELIVER`