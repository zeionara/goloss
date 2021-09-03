# goloss

A small tool which performs automatic adjustments of the sound level on your device to keep in down even when something loud happens

## Install dependencies

Install pavumeterc which is used for live volume level detection using pulse-audio library

```sh
sudo apt-get install autoconf libglibmm-2.4-dev libpulse-mainloop-glib0 libpulse-dev lynx
git clone git@github.com:CoolDuke/pavumeterc.git
cd pavumeterc/
./bootstrap.sh
./autogen.sh
./configure 
make
sudo make install
```

Install go language dependencies

```sh
go mod download
```

## Usage

To run the tool execute the following command in which the only argument is minimum time delay between two consequent volume adjustments made by the script

```sh
go run main.go 10
```
