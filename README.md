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
    name: one
    path: echo
    
    Query:
      - one

  - follow: one
    type: execute
    name: two
    path: echo
    
    Query:
      - two

  - follow: two
    type: execute
    path: echo
    
    Query:
      - three

  - follow: two   
    type: execute
    connection: agent
    path: bash

    Query:
      - -c 'echo $PATH'

  - type: upload
    from: example
    connection: agent
    to: /home/admin/example

  - follow: upload
    type: download
    connection: agent
    from: /home/admin/example
    to: example
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
    from: main
    to: update/main #optional
```
#### Move
```
  - type: move
    from: main
    to: update/main #optional
```
#### Upload
```
  - type: upload
    from: main
    connection: agent
    to: /home/admin/main #optional
```
#### Move
```
  - type: download
    connection: agent
    from: /home/admin/main
    to: main #optional
```
#### Execute
```
  - type: execute
    connection: optional
    path: go
    timeout: 8 #optional
    
    Environment:
      - CGO_ENABLED=0
    Query:
      - build
```
