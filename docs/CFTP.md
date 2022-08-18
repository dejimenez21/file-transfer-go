# Custom File Transfer Protocol

The Custom File Transfer Protocol (**CFTP**) is a TCP-based application layer protocol for file transfer between multiple clients. This protocol is designed for a client - server architecture. It works as support for a publisher - subscriber pattern, where clients can send and receive files through channels.

## Messages

There are two types of messages: Requests and File Chunks. Request are used for the majority of the cases, while File Chunks are specifically used by the server for delivering files content to clients.

### Requests

Requests are structured as follows:  

1. **[Method](#method)** that indicates the action to perform.
2. **Meta** section, that is a json string containing data about the request.
3. **Channels** section, a comma-separated list of the channels involved in the request.
4. **File Information** a json string containing data about the file being transfer.

A request message should always end with a EOT character.

The only fully required section is the **Method**, all the others may or not be specified depending on the operation to perform. For example, a client who wants to suscribe to a channel doesn't need to include the **File Information** section on the request.

Every section described above should be separated by a line break (`"/n"`). Even if a section is not to be included the line break should.
 For example:

```ps1

SEND

chn1,chn2
[FILE_INFORMATION]
 ```

In the above example, the **Meta** section is not included, but the empty line where it wolud be is included.

>(Each section is described more explicitly in following parts of this documentation.)

### File Chunks

File Chunks are messages that contain parts of the content of the file being delivered to a client. A File Chunk has a **header** and a **body**. The structure of the header is the following:

1. The **chunk** keyword that identifies this message as a file chunk.
2. The **Delivery ID** that is used to indentify all related chunks belonging to the same file.
3. A **Sequence number** that indicates the right position that the current content chunk should take when storing the file being received.
4. The **Size** of the chunk being delivered.

On the other hand, the **body** is nothing more than the content chunk itself, in the format of a array of bytes.

> Every part of the header described above is required and should be separated by a line break (`"/n"`). The header and the body should be separated by a EOT character.

## Requests Sections

### Method

The Method represents the action to perform and is the first section of the request. Methods are upppercase words that represent verbs of the english language. Clients and servers use different methods.

Client's methods:

* `send`: tells the server that the client will begin streaming a file.
* `suscribe`: request to suscribe to a list of channels to start receiving every file sent through them.

Server's methods:

* `deliver`: tells the client that a file will begin to be delivered in the form of file chunk messages.

### Meta

The Meta section is defined as a json string to give the developers the ability to add as many information of the current request as they believe necessary.
