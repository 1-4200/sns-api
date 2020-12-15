# sns-api
## Layered architecture

    [handler] 
    ↓↓↓↓
    [usecase]
    ↓↓↓↓
    [domain]
    ↑↑↑↑
    [infrastructure]

## Directory structure
    .
    ├── api # Configure the api server
    ├── config
    ├── domain # Implementation related to business logic
    ├── handler # Implementations related to Request and Response
    │   └── tweet
    ├── infrastructure # Implementation related to technology
    │   ├── elastic
    │   ├── mysql
    ├── usecase # Implementation of usecase involved
    ├── log # Folder for log output
    └── logger # Logging process (global)
