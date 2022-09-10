# nofmt
A gofmt alternative that gets out of your way.

## Installation

Install `go`
```bash
#--- Installation ---
VER=1.19.1
cd $HOME
rm -rf $HOME/goroot
mkdir -p $HOME/goroot $HOME/go/bin
wget "https://dl.google.com/go/go${VER}.linux-amd64.tar.gz"
tar -xf "go${VER}.linux-amd64.tar.gz" -C $HOME/goroot --strip-components 1 go
hash -r

echo 'export PATH="$HOME/bin:$HOME/go/bin:$HOME/goroot/bin:$PATH"' >> $HOME/.bashrc
```
Log out and back in again to ensure the PATH is updated.

Build `nofmt`:
```bash
git clone https://github.com/yqdv/nofmt.git
cd nofmt

# Backup the go/printer/nodes.go file to go/printer/nodes.go.orig
cp $HOME/goroot/src/go/printer/{nodes.go,nodes.go.orig}

# Patch go/printer/nodes.go to update linebreak()
patch $HOME/goroot/src/go/printer/nodes.go < inject/go/printer/nodes.go.patch

# Build "gofmt" using the modified source and save it as "nofmt"
cd $HOME/goroot/src/cmd/gofmt
go build
mv gofmt $HOME/go/bin/nofmt
hash -r

# Revert the source file patch
mv -f $HOME/goroot/src/go/printer/{nodes.go.orig,nodes.go}
```

Install vim-go:
```bash
mkdir -p $HOME/.vim/pack/plugins/start
git clone 'https://github.com/fatih/vim-go' $HOME/.vim/pack/plugins/start/vim-go

mkdir -p $HOME/.vim/after/plugin
cat << 'EOF' > $HOME/.vim/after/plugin/vim-go.vim
let g:go_fmt_command = 'nofmt'
let g:go_imports_mode = 'gopls'

let g:go_fmt_autosave = 1
let g:go_imports_autosave = 1
EOF
```

Inside vim, run:
```vim
:GoInstallBinaries
```

The `nofmt` binary works identically to `gofmt` except for having some different formatting rules.
Now `nofmt` (and `goimports`) will be run automatically on file save.

Note that you can still manually run `gofmt` from within vim when needed:
```bash
#--- Run command on whole file ---
:%!gofmt

#--- Run command on selection (selected function) ---
Ctrl-v and select a whole function (selection must be parseable)
:!gofmt

#--- Run GoImports (via gopls) ---
:GoImports
```

Wait! I want spaces instead of tabs!

No problem. Install Go and vim-go (including GoInstallBinaries).
Then replace `nofmt`, `goimports`, and `gopls` (make sure to stop the `gopls` server first):
```bash
mkdir -p $HOME/src
cd $HOME/src

# Clone golang/tools and nofmt
git clone https://github.com/golang/tools.git
git clone https://github.com/yqdv/nofmt.git
cd nofmt

# Patch go/printer/nodes.go (for nofmt)
cp $HOME/goroot/src/go/printer/{nodes.go,nodes.go.orig}
patch $HOME/goroot/src/go/printer/nodes.go < inject/go/printer/nodes.go.patch

# Patch go/format/format.go (for spaces)
cp $HOME/goroot/src/go/format/{format.go,format.go.orig}
patch $HOME/goroot/src/go/format/format.go < inject/go/format/format.go.patch

# Patch go/doc/comment/print.go (for spaces)
cp $HOME/goroot/src/go/doc/comment/{print.go,print.go.orig}
patch $HOME/goroot/src/go/doc/comment/print.go < inject/go/doc/comment/print.go.patch

# Patch cmd/gofmt/gofmt.go (for spaces)
cp $HOME/goroot/src/cmd/gofmt/{gofmt.go,gofmt.go.orig}
patch $HOME/goroot/src/cmd/gofmt/gofmt.go < inject/cmd/gofmt/gofmt.go.patch

# Patch tools/cmd/goimports/goimports.go (for spaces)
cp $HOME/src/tools/cmd/goimports/{goimports.go,goimports.go.orig}
patch $HOME/src/tools/cmd/goimports/goimports.go < inject/tools/cmd/goimports/goimports.go.patch

# Build new nofmt
cd $HOME/goroot/src/cmd/gofmt
go build
mv gofmt $HOME/go/bin/nofmt

# Revert patched nodes.go (to not get rolled imports)
mv -f $HOME/goroot/src/go/printer/{nodes.go.orig,nodes.go}

# Build new goimports and gopls
cd $HOME/src/tools/cmd/goimports
go build
mv goimports $HOME/go/bin/goimports

cd $HOME/src/tools/gopls
go build
mv gopls $HOME/go/bin/gopls

hash -r

# Revert patched files
mv -f $HOME/goroot/src/go/format/{format.go.orig,format.go}
mv -f $HOME/goroot/src/go/doc/comment/{print.go.orig,print.go}
mv -f $HOME/goroot/src/cmd/gofmt/{gofmt.go.orig,gofmt.go}
mv -f $HOME/src/tools/cmd/goimports/{goimports.go.orig,goimports.go}

# Set expandtab for spaces
mkdir -p $HOME/.vim/after/ftplugin
cat << 'EOF' > $HOME/.vim/after/ftplugin/go.vim
setlocal expandtab
EOF
```
Note: You will get rolled imports if nodes.go is patched when goimports and gopls are built.
