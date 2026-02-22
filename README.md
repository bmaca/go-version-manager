# Go Version Manager
A simple, single binary CLI tool to install and manage Go version on machines(for my case linux servers). Zero dependenices included.

### Why?
Sometimes on my servers I dont want docker just for language tooling. some existing tooling like *GVM* and *asdf* seem a bit more complicated for somethign as simple as what I need it for and thats just install this version of GO in this location. I also just wanted something I can SCP over to machines via ansible or just manual and it just works.

### What it does currently does(this can change)
* Downloads official Go releases from https://go.dev/dl
* extracts to your desired location(version isolated)
* Verifies SHA256 checksums
* No shell dependencies, or interpeters or config files.

### Usage
```
# Install a specific version to your desired directory
./goversion install --version 1.26.0 --dir $HOME/go_versions
```
### You can then aliases them if you want
```
~/bashrc
alias go1.26='/home/root/go_versions/1.26.0.bin/go'
```

### Future Ideas
* `list` and `uninstall` commands would be nice to have.
* potentially a symlink manager.