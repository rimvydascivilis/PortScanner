# Port scanner

## Usage
### 1. From the command line in the PortScanner directory, run the go build command to compile the code into an executable
  
  ```go build```

### 2. Run the executable
  - ```./PortScanner --help``` to get help
  
#### Examples:
  - ```./PortScanner -p 100:1500 -h <host>```
  - ```./PortScanner -n ``` (default port range 0-65535, default host localhost)
