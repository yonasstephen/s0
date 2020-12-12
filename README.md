# S0
S0 is a resilient, distributed file storage. This is a self-learning project of building a distributed system that is able to scale horizontally and is fault tolerant. 

# Getting Started
Running a single instance S0
```
make all
```

# API
S0 has 2 simple API:
1. POST /upload
This API handles static file upload by the key `static_file`. Ensure that `Content-Type` is set to `application/x-www-form-urlencoded`
2. GET /[filename.extension]
This API serves the uploaded static file
