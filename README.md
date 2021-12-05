# File Encryption-Decryption Web-app

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity) ![License](https://img.shields.io/cran/l/devtools)

A file encryption-decryption HTTP server. The server encrypts or decrypts the file present on the server. File encryption and decryption is performed using [DES](https://www.cs.uri.edu/cryptography/dessimplified.htm#:~:text=Simplified%20DES%20is%20an%20algorithm,on%20blocks%20of%2012%20bits.) algorithm.
The encryption and decryption operation are exposed using REST APIs(The usage is explained below). For efficient routing and implementation of REST APIs gorilla mux server is used.

## How To Use
1. Clone the project
2. Open the folder in command prompt
3. Run "go build"
4. The application is shipped as a CLI tool. User can perform encryption and decryption of files from terminal by simply running the commands.
<br />

### Encryption
* Run the below command to encrypt a file.
```
fes encryptfile --filepath=</path/to/your/file>
```
<br />

### Decryption
* Run the below command to encrypt a file.
```
fes encryptfile --filepath=</path/to/your/file>
```
<br />

5. The application provides a way to set up a HTTP server through the CLI commands.
<br />
Run the below command to setup the HTTP server.
```
fes startServer
```
The above command will spawn up a HTTP server with encryption and decryption endpoints.
<br />

=> For encrypting a file, send a HTTP GET request using the url: http://localhost:8080/encrypt?filename=YourFileName. The server creates a file named as YourFileName.enc at the same location where the original file is present.
<br />

=> For decrypting a file, send a HTTP GET request using the url: http://localhost:8080/decrypt?filename=YourFileName. In this case, only the name of the file and not the ".enc" should be sent through the url. The file is decrypted and stored at the same location of the encrypted file.

## Dependencies
1. [gorilla-mux](https://github.com/gorilla/mux): A powerful HTTP router and URL matcher for building Go web servers with.
2. [cobra-CLI](https://github.com/spf13/cobra): Cobra is both a library for creating powerful modern CLI applications as well as a program to generate applications and command files.

NOTE:- Indirect dependencies are not listed here.

<br>

## Configurations
The application uses two ```json``` configuration files and one ```.txt```.
1. des_input.txt - Contains the inputs or properties of the Simplified-DES algorithm. The file looks like this:
```
key:0010010111
P10:3 5 2 7 4 10 1 9 8 6
P8:6 3 7 4 8 5 10 9
P4:2 4 3 1
IP:2 6 3 1 4 8 5 7
IP-1:4 1 3 5 7 2 8 6
E/P:4 1 2 3 2 3 4 1
S0:1 0 3 2 3 2 1 0 0 2 1 3 3 1 3 2
S1:0 1 2 3 2 0 1 3 3 0 1 0 2 1 0 3
```
The file contains colon ```:``` separated key pair values. Following is the brief description of properties:
```
    i. P10 - ten-length permutation of the algorithm. Ten space separated integers denoting the permutation.
    ii. P8 - eight-length permutation of the algorithm. Eight space separated integers denoting the permutation.
    iii. P4 - four-length permutation of the algorithm. Four space separated integers denoting the permutation.
    iv. IP - initial permutation of the algorithm. Eight space separated integers denoting the permutation.
    v. IP-1 - inverse permutation of the algorithm. Eight space separated integers denoting the permutation.
    vi. E/P - expansion permutation of the algorithm. Eight space separated integers denoting the permutation.
    vii. S0 - s0 matrix of the algorithm. Sixteen space separated integers denoting the permutation.
    viii. S1 - s1 matrix of the algorithm. Sixteen space separated integers denoting the permutation.
    ix. key - cipher key of the algorithm. Ten-length binary string.
```

2. des_config.json: Input values are provided as json attribute.
```json
{
	"fastmode": true,
	"threadcount": 8,
    "buffersize": 1048576
}
```
* **fastmode** - boolean attribute. If set to true, then encryption of a file will be done multithreaded manner.
* **threadcount** - Integer attribute. Defines how many buffers will be encrypted concurrently at a time.
* **buffersize** - Integer attribute. Defines the size of the buffer in bytes which will be read from file, encrypted and written back to disk. Default value is ```1048576``` i.e 1 MB.

3. config.json - Defines the configurations of the application server.
```
{
	"name": "FES App Server",
	"port": 8080,
    "protocol": "HTTP"
}
```
* **name** - Name of the application
* **port** - port on which app server listens for request.
* **protocol** - name of the protocol used by the server.

## Want to contribute
1. Explain the issue that must be fixed or a feature that must to be added.
2. Fork the repository to your github account.
3. Make your changes.
4. Create a pull request from your forked repository to master branch of this repository

## License
This project is licensed under GNU GPL v3 license.
Note:- The license is subject to changes in future.