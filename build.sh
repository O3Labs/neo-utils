#Android NDK is required. https://developer.android.com/ndk/guides/index.html
#to build for Android you need to run the command below
#gomobile init -ndk ~/Library/Android/sdk/ndk-bundle/
#if you are using go 1.9.4 and ran into the issue related to #cgo CFLAGS: -fmodules
#here is the workaround:  export CGO_CFLAGS_ALLOW='-fmodules|-fblocks'
rm -rf output/*
mkdir output/ios/
mkdir output/android/
echo "Building for iOS..."
CGO_CFLAGS_ALLOW='-fmodules|-fblocks' gomobile bind -target=ios -ldflags=-w -o=output/ios/neoutils.framework github.com/o3labs/neo-utils/neoutils
echo "Building for Android..."
ANDROID_HOME=/Users/$USER/Library/Android/sdk gomobile bind -target=android -o=output/android/neoutils.aar github.com/o3labs/neo-utils/neoutils
