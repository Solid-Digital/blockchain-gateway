git diff --name-only --diff-filter=d HEAD origin/develop | cpio -pd diff && rm -r lib priv test && mv diff/lib ./lib && mv diff/priv ./priv && mv diff/test ./test
