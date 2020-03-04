# Zipalign APK

Runs the zipalign Android build tool on a signed APK. Allows for page aligning all shared object files in an APK to ensure libraries are properly extracted after install. Fixes "Failed to extract native libraries" error caused by "lib-xxxx.so is not page-aligned" install error.


## How to use this Step

Can be run directly with the [bitrise CLI](https://github.com/bitrise-io/bitrise),
just `git clone` this repository, `cd` into it's folder in your Terminal/Command Line
and call `bitrise run test`.