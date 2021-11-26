  
#  Mainland Node Server
  
This program implements the serverside of the Mainland Node, as used within our test stand's electronics system.
##  What's a Mainland Node?
  
Time for a quick overview on how Liquid Propulsion's electronics system will work:

![](docs/assets/2d056a8107bbf678b95046bb33bf6b510.png?0.7619751023061705)  
Our new electronics system is based around a CAN Bus, whereby there are a series of computerized "nodes" which can broadcast 8 byte messages onto the bus in order to communicate with each other. The Mainland Node effectively acts as the the only active node in the system, it broadcasts two types of packets: update ID and state. The Update ID packet forces a node to change it's ID, then the state packet updates the bus' state. For example, during cylinder pressurization a state packet with the stage ID of 1 may be used to open NPV (Nitrogen Pressure Valve) which will then backpressure the fuel tank. The state packet must be repeated every 20ms, with an error of no more than 20ms, otherwise all Island boards will go into a disabled state and will remain disabled for 5 seconds.
  
Thus, the Mainland Server implements the CAN bus interface and networking interfaces necessary to allow for Mission Control to safely control our test setup remotely.
##  Requirements
  
The Mainland Server simply requires Golang 1.16 or newer. It can be downloaded [here](https://go.dev )
##  Running the Server
  
First, clone the repository and all of it's submodules:
```sh
git clone https://github.com/Liquid-Propulsion/mainland-server.git
  
cd mainland-server
  
git submodule update --init --recursive
```
Then simply run the go file, or compile it into a binary:
```sh
# Run via Golang
go run server.go
  
# Compile to binary (Linux/Mac)
go build . -o main
./main
  
# Compile to binary (Windows)
go build . -o main.exe
./main.exe
```
  