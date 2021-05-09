## Description
A small application representing CRUD operations on a university database.

## Paths table
<table>
<tr>
<td>Path</td>
<td>Method</td>
<td>Description</td>
<td>Body example</td>
</tr>
<tr>
<td>/courses</td>
<td>GET</td>
<td>Get all courses</td>
<td>
  
```json
{
  "courses":
    [
      {
        "Code":207,
        "Title":"Mobile Application Development",
        "DepartmentCode":5,
        "Description":"Mobile Application Development course description..."
      },
      {
        "Code":208,
        "Title":"Java Web Development",
        "DepartmentCode":5,
        "Description":"Java Web Development course description..."
      },
      {
        "Code":209,
        "Title":"Architecture Operating Systems",
        "DepartmentCode":5,
        "Description":"Architecture Operating Systems course description..."
      }
    ]
}
```

</td>
</tr>
<tr>
<td>/courses/{code}</td>
<td>GET</td>
<td>Get course by code</td>
<td>
  
```json
{
  "course":
      {
        "Code":207,
        "Title":"Mobile Application Development",
        "DepartmentCode":5,
        "Description":"Mobile Application Development course description..."
      }
}
```

</td>
</tr>
<tr>
<td>/courses</td>
<td>POST</td>
<td>Create new course</td>
<td>
  
  
```json
{
  "code":101
}
```
</td>
</tr>
<tr>
<td>/courses/{code}</td>
<td>PATCH</td>
<td>Update course description</td>
<td>
  
  
```json
{
  "code":101
}
```
</td>
</tr>
<tr>
<td>/courses/{code}</td>
<td>DELETE</td>
<td>Delete course by code</td>
<td>
  
  
```json
{}
```
</td>
</tr>
</table>

## How to run
1. Run server with command: `go run ./server/cmd/main.go` (port `:8080` by default)
2. Run gateway with command: `go run ./client/cmd/main.go` (port `:8383` by default)

The application will be aviable on the following address: `localhost:8383/{path}`.

## Unit tests
```
go test -race
```

## Protocol Buffers
To generate go models you should be in `/proto` folder. Run the following command:
```
protoc -I . course.proto --grpc-gateway_out . --go_out=plugins=grpc:.
```

## Docker
`/docker` folder contains all config files. To run the application in containers enter the following command:
```
docker-compose up
```
