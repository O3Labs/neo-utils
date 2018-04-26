# neo-utils
Useful functions for NEO blockchain written in Go. 

This library can be fully compiled to native iOS and Android framework using `gomobile bind`

note: `gomobile` does not support slice parameter yet so some functions are optimized to take a comma separated string as a param instead of a array of string.  
For methods specfically designed to be used on mobile see ```mobile.go```  

### Installation 
`go get github.com/o3labs/neo-utils/neoutils`


## Compile this library to native mobile frameworks.

### Install gomobile
`go get golang.org/x/mobile/cmd/gomobile`  
`gomobile init`  


### Compile to iOS framework
XCode is required.  
`gomobile bind -target=ios -o=output/ios/neoutils.framework github.com/o3labs/neo-utils/neoutils`

### Compile to Android framework
Android NDK is required. https://developer.android.com/ndk/guides/index.html  
`gomobile init -ndk ~/Library/Android/sdk/ndk-bundle/`  
`ANDROID_HOME=/Users/apisit/Library/Android/sdk gomobile bind -target=android -o=output/android/neoutils.aar github.com/o3labs/neo-utils/neoutils`
