# Port scanner

## Usage
### 1. From the command line in the PortScanner directory, run the go build command to compile the code into an executable.
  
  ```go build```

### 2. Run the executable

  ```./PortScanner -p <port range from>:<port range to> <host>```
  
 #### Examples:
  - ```./PortScanner -p 100:1500 127.0.0.1```
  - ```./PortScanner localhost``` (default port range 0-65535)
