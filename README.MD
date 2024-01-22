Example web app in golang

## App overview

Example web service app exposing REST endpoints.

* Web server
* DB persistence
* Tests

This app is an example web service app:

- it will accept test payload consisting of questions and answers, and persist them as test-result
  - input: questions, answers, user jwt-token (containing user uuid)
  - output: error or test-result-uuid
- it will accept test-result-uuid and user-uuid to verify test belongs to user
  - input: test-result-uuid, user-uuid
  - output: test payload (questions, answers), timestamp of completing the test

## Requirements
- Golang
- Available HTTP port 8000

## Limitations
- Persistence layer uses `nutsdb/nutsdb` which persist data between runs, but limitation is only one running server

### Build & Run & Test on host

`go build -race -ldflags="-s -w" -o bin/signer ./cmd/signer`

`go run -race ./cmd/signer`

`go test -race -v ./...`

### CI/CD

```
# build web server binary container
docker build . -t signer -f cmd/signer/Dockerfile
docker run -it signer

# build worker binary container
docker build . -t signer-worker -f cmd/worker/Dockerfile
docker run -it signer-worker
```

### Cleanup DB [nutsdb](https://github.com/nutsdb/nutsdb) is fs-backed

```
rm -rf /tmp/signerdb
```

### Example web app HTTP API endpoints

`POST /users/:userID/homework` Saves test assignment results from user, returns ID of the saved assignment.

```
curl -H 'Content-Type: application/json' -H "Authorization: Bearer $TOKEN" -d '{"questions": "1", "answers": "2"}' http://0.0.0.0:8000/users/$USER1/homework
```

`GET /users/:userID/homework/:homeworkID` Fetches previously saved assignment, returns questions & answers.

```
#check value of previous POST request and put ID here
HOMEWORKID=replaceme
curl -H "Authorization: Bearer $TOKEN" http://0.0.0.0:8000/users/$USER1/homework/$HOMEWORKID
```

### Example execution
```
➜  go-signer git:(main) ✗ curl -v -H 'Content-Type: application/json' -H "Authorization: Bearer $TOKEN" -d '{"questions": "1", "answers": "2"}' http://0.0.0.0:8000/users/$USER1/homework
*   Trying 0.0.0.0:8000...
* Connected to 0.0.0.0 (127.0.0.1) port 8000 (#0)
> POST /users/1055bb1c-860a-4c13-853a-ca56073511d0/homework HTTP/1.1
> Host: 0.0.0.0:8000
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Type: application/json
> Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDU1YmIxYy04NjBhLTRjMTMtODUzYS1jYTU2MDczNTExZDAifQ.hRZ8rYGUq_mPMwSt7XcyBGxEQOu3YYO9NxhtlujICm0
> Content-Length: 34
>
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=UTF-8
< X-Request-Id: DavECTQtvYzUxFoYSqVFNARIKpkdnBQR
< Date: Wed, 17 Jan 2024 01:21:31 GMT
< Content-Length: 51
<
{
  "id": "c0aa3998-50c3-4092-acad-51c8ad553bcf"
}
* Connection #0 to host 0.0.0.0 left intact

➜  go-signer git:(main) ✗ curl -v -H "Authorization: Bearer $TOKEN" http://0.0.0.0:8000/users/$USER1/homework/$HOMEWORKID
*   Trying 0.0.0.0:8000...
* Connected to 0.0.0.0 (127.0.0.1) port 8000 (#0)
> GET /users/1055bb1c-860a-4c13-853a-ca56073511d0/homework/c0aa3998-50c3-4092-acad-51c8ad553bcf HTTP/1.1
> Host: 0.0.0.0:8000
> User-Agent: curl/8.1.2
> Accept: */*
> Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDU1YmIxYy04NjBhLTRjMTMtODUzYS1jYTU2MDczNTExZDAifQ.hRZ8rYGUq_mPMwSt7XcyBGxEQOu3YYO9NxhtlujICm0
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< X-Request-Id: pxITkQXwLJpNnBPivApMalGbYquukfYR
< Date: Wed, 17 Jan 2024 01:22:10 GMT
< Content-Length: 88
<
{
  "questions": "1",
  "answers": "2",
  "created_at": "2024-01-17T01:21:31.316578Z"
}
* Connection #0 to host 0.0.0.0 left intact

```
