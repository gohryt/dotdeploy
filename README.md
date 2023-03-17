## Installation
Run  
```go install github.com/gohryt/dotdeploy/command/dotdeploy```
on the pc with go installed.  
Then run ```dotdeploy --help``` to check it was installed correctly.  
## Basics
Make ```.deploy``` file in your working folder. Deploy file should looks like below and filled with commands:  
```
folder: "update"

Remote:
  - type: agent
    host: 70.34.202.107
    username: root

Do:
  - type: execute
    name: two
    Path:
      connection: agent
      path: echo
    Query:
      - two

  - follow: two   
    type: execute
    Path: 
      path: bash
    Query:
      - -c 'echo $PATH'
```
#### Options
```folder``` is folder wich will be created on start and deleted on end of processing.  
```keep``` is a flag which means that program should't delete ```folder```.  
## Remote
#### Key
```
  - type: password
    host: 1.1.1.1
    file: /home/example/.ssh/id_ed25519
    username: root
    password: example
```
#### Password
```
  - type: password
    host: 1.1.1.1
    username: root
    password: example
```
#### Agent
```
  - type: agent
    host: 1.1.1.1
    username: root
```
## Do
#### Copy
```
  - type: copy
    From:
      connection: agent #optional
      path: main
    To:
      connection: agent #optional
      path: update/main
```
#### Move
```
  - type: move
    From:
      connection: agent #optional
      path: main
    to: update/main     #optional
```
#### Execute
```
  - type: execute
    timeout: 8    #optional
    Path:
      connection: #optional
      path: go
    Environment:
      - CGO_ENABLED=0
    Query:
      - build
```
