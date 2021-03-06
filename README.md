#Description
This is terminal application using input and output from terminal to simulate a 3D environment. The rendering is done using space characters using color background.

#Development
To debug Using VSCode, as Delve cannot be used in interqctive mode, the headless mode is used:
```dlv debug --headless --listen=:2345 --log --api-version=2```
Then the following remote-launch configuration is used:
```
{
    "name": "3dEngine Remote debug",
    "type": "go",
    "request": "attach",
    "mode": "remote",
    "remotePath": "${workspaceFolder}",
    "port": 2345,
    "host": "127.0.0.1",
}
``` 

To debug Using VSCode, as Delve cannot be used in interactive mode, the headless mode is used:
```dlv debug --headless --listen=:2345 --log --api-version=2```
Then the following remote-launch configuration is used:
```
{
    "name": "3dEngine Remote debug",
    "type": "go",
    "request": "attach",
    "mode": "remote",
    "remotePath": "${workspaceFolder}",
    "port": 2345,
    "host": "127.0.0.1",
}
``` 

#Build/Execution
Inside the source folder:
```go build && go install```
and then
```$GOPATH/bin/3dGame````

#client/server
* launch server
```go build && ./3dGame --mode remoteServer```
* launch client
```go build && ./3dGame --mode remoteClient```
* debug client headless (using config file above)
```dlv debug --headless --listen=:2345 --log --api-version=2 -- --mode remoteClient```

