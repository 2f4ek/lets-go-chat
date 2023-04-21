# Lets GO Chat
[![Go Reference](https://pkg.go.dev/badge/github.com/2f4ek/lets-go-chat.svg)](https://pkg.go.dev/github.com/2f4ek/lets-go-chat)
## Description
This is a chat application that allows users to communicate with each other in real-time. 
## Features
- User authentication (under development)
- Real-time messaging (under development)
## Technologies 
   - GoLang
## Installation 
   1. Clone the repository ```bash git clone https://github.com/2f4ek/lets-go-chat``` 
   2. Go to the main app directory ```cd cmd/app``` 
   3. Install dependencies ```bash go get``` 
   4. Run the server ```bash go run main.go``` 
## Usage 
Please note that currently only the features described below are available. All other functionality is under development.
## Functions: 
### [HashPassword](https://github.com/2f4ek/lets-go-chat/blob/main/pkg/hasher/hasher.go#L10)
```go 
func HashPassword(password string) string
``` 
This function is used to hash a user's password for security purposes. 
### Input 
- `password` - A string representing the user's password. 
### Output
- A string representing the hashed password. 
### Usage
```go 
hashedPassword := HashPassword("password") 
```
###
### [CheckPasswordHash](https://github.com/2f4ek/lets-go-chat/blob/main/pkg/hasher/hasher.go#L19)
```go
CheckPasswordHash(password, hash string) bool
```
This function is used to check if a given password matches a given hash. 
### Input
- `password` - A string representing the user's password. 
- `hash` - A string representing the hashed password. 
### Output 
- A boolean value indicating whether the password matches the hash or not. 
### Usage 
```go
isMatch := CheckPasswordHash("password", hashedPassword)
```

### License This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.