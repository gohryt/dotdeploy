## Basics
Make ```.deploy``` file in your working folder. Deploy file should looks like below and filled with commands:  
```
folder: "update"

Remote:
  - type: agent
    host: 70.34.202.107
    username: root

Do:
  - type: run
    name: one
    path: echo
    
    Query:
      - one

  - follow: one
    type: run
    name: two
    path: echo
    
    Query:
      - two

  - follow: two
    type: run
    path: echo
    
    Query:
      - three
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
#### Run
```
  - type: run

    path: go
    timeout: 8 #optional
    
    Environment:
      - CGO_ENABLED=0
    Query:
      - build
```
