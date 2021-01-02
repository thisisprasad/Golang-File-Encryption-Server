# File Encryption-Decryption Web-app

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity) ![License](https://img.shields.io/cran/l/devtools)

A file encryption-decryption HTTP server. The server encrypts or decrypts the file present on the server. File encryption and decryption is performed using [DES](https://www.cs.uri.edu/cryptography/dessimplified.htm#:~:text=Simplified%20DES%20is%20an%20algorithm,on%20blocks%20of%2012%20bits.) algorithm.
The encryption and decryption operation are exposed using REST APIs(The usage is explained below). For efficient routing and implementation of REST APIs gorilla mux server is used.

## How To Use
1. Clone the project
2. Open the folder in command prompt
3. Run "go build"
4. Run "fes"
5. For encrypting a file, send a HTTP GET request using the url: http://localhost:8080/encrypt?filename=YourFileName. The server creates a file named as YourFileName.enc at the same location where the original file is present.
6. For decrypting a file, send a HTTP GET request using the url: http://localhost:8080/decrypt?filename=YourFileName. In this case, only the name of the file and not the ".enc" should be sent through the url. The file is decrypted and stored at the same location of the encrypted file.

## Dependencies
1. [gorilla-mux](https://github.com/gorilla/mux): A powerful HTTP router and URL matcher for building Go web servers with 

## Want to contribute
1. Explain the issue that must be fixed or a feature that must to be added.
2. Fork the repository to your github account.
3. Make your changes.
4. Create a pull request Lm your forked repository to master branch of this repository

## License
This project is licensed under GNU GPL v3 license.
Note:- The license is subject to changes in future.