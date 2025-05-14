## AWDS Execution Guide
### 1. Building the Executable
- Create the executable using Makefile
```
    make
```
- A `/bin` folder will be created in the `/awds` directory with the executable file (awds) inside.

### 2. Running AWDS
```
    ./bin/awds
```
- Execute AWDS from the awds/ directory.

--- 
## Main Features
Since there are many items to include in the Body, it's recommended to use environments like Postman rather than curl. \
Send the Body in raw JSON format. Testing was conducted in a Postman environment.

### Device Registration (POST)
- Request path
```
    http://{HOST}:{PORT}/devices
```

- Body content
```
    {
        "end_point": {device endpoint},
        "description": {description, optional}
    }
```

### Pod Registration (POST)
- Request path
```
    http://{HOST}:{PORT}/pods
```

- Body content
```
    {
        "end_point": {pod endpoint},
        "description": {description, optional}
    }
```

### Job Registration (POST)
- Request path
```
    http://{HOST}:{PORT}/jobs
```

- Body content
```
    {
        "device_id": {device endpoint},
        "pod_id": {description, optional},
        "input_size": {input size, integer, optional},
    }
```

### Schedule (POST - Not GET!)
- Request path
```
    http://{HOST}:{PORT}/schedules/{job_id}
```

- No Body content

### Pod Retrieval (GET)
- Request path
```
    http://{HOST}:{PORT}/pods
```

- Returns pod list

### Device Retrieval (GET)
- Request path
```
    http://{HOST}:{PORT}/devices
```

- Returns device list

### Job Retrieval (GET)
- Request path
```
    http://{HOST}:{PORT}/jobs
```

- Returns job list

## Additional Convenience Features (PATCH, DELETE)
### Device Update (PATCH)
- Request path
```
    http://{HOST}:{PORT}/devices/{device_id}
```

- Body content
```
    {
        "end_point": {device endpoint},
        "description": {description, optional}
    }
```

### Device Deletion (DELETE)
- Request path
```
    http://{HOST}:{PORT}/devices/{device_id}
```

- No Body content

### Pod Update (PATCH)
- Request path
```
    http://{HOST}:{PORT}/pods/{pod_id}
```

- Body content
```
    {
        "end_point": {device endpoint},
        "description": {description, optional}
    }
```

### Pod Deletion (DELETE)
- Request path
```
    http://{HOST}:{PORT}/pods/{pod_id}
```

- No Body content

### Job Update (PATCH)
- Request path
```
    http://{HOST}:{PORT}/jobs/{job_id}
```

- Commonly used Body
```
    {
        "device_id": {device endpoint},
        "pod_id": {description, optional},
        "input_size": {input size, integer, optional},
        "completed": {completion status, boolean, optional},
    }
```

- Possible Body content
```
    {
        "device_id": {device endpoint},
        "pod_id": {description, optional},
        "input_size": {input size, integer, optional},
        "partition_rate": {distribution ratio, float, optional},
        "completed": {completion status, boolean, optional},
        "DeviceStartIndex": {device start index, integer, optional},
        "DeviceEndIndex": {device end index, integer, optional},
        "PodStartIndex": {pod start index, integer, optional},
        "PodEndIndex": {pod end index, integer, optional},
    }
```

### Job Deletion (DELETE)
- Request path
```
    http://{HOST}:{PORT}/jobs/{job_id}
```

- No Body content