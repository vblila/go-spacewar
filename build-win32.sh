# sudo apt-get install gcc-multilib
# sudo apt-get install gcc-mingw-w64

GOOS=windows GOARCH=386 \
  CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc \
  go build -ldflags "-s -w" -o ./bin/win32/spacewar.exe