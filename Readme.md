# REST vs gRPC

## Overview
A concise comparison between **REST** and **gRPC** for designing APIs and microservices.

---

## REST 
- Protocol: HTTP/1.1 (commonly) with JSON payloads.
- Style: Resource-oriented and human-readable.
- Strengths:
  - Broad ecosystem and easy to test with curl / browsers.
  - Good for public HTTP APIs and services targeting browsers.
  - Simple and debuggable (JSON, plain text).
- Trade-offs:
  - Larger, textual payloads (more bandwidth).
  - No built-in streaming; less efficient binary handling.

---

## gRPC 
- Protocol: gRPC uses HTTP/2 and Protocol Buffers (binary).
- Style: RPC (Remote Procedure Call)  contract-first via `.proto` files.
- Strengths:
  - High performance (binary, multiplexed over HTTP/2).
  - Strongly typed contracts and auto-generated client/server code.
  - First-class streaming: unary, server-streaming, client-streaming, bidirectional.
- Trade-offs:
  - Less human-readable without tooling (binary payloads / proto).
  - Requires code generation tooling and HTTP/2 support at proxies.

### What is RPC? 
RPC stands for **Remote Procedure Call**. It is an abstraction that lets a program call a function on a remote server as if it were a local function. gRPC implements RPC semantics by:
- Defining services and message types in `.proto` files (the contract).
- Generating **stubs** (clients) and **skeletons** (servers) from that contract.
- Handling serialization (Protocol Buffers), network transport (HTTP/2), and method dispatch behind the scenes.

Common RPC call types in gRPC:
- **Unary**: single request, single response.
- **Server streaming**: client request, server streams multiple responses.
- **Client streaming**: client streams multiple requests, server sends single response.
- **Bidirectional streaming**: both sides stream messages independently.

---
Run this command after updating [sso.proto]

```
protoc -I .\proto --go_out=.\gen\go --go_opt=paths=source_relative --go-grpc_out=.\gen\go --go-grpc_opt=paths=source_relative .\proto\sso\sso.proto 
```
