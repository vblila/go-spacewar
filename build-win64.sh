# sudo apt-get install gcc-multilib
# sudo apt-get install gcc-mingw-w64

GOOS=windows GOARCH=amd64 \
  CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc \
  go build -ldflags "-s -w" -o ./bin/win64/spacewar.exe