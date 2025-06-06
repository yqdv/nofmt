# nofmt
A gofmt alternative that gets out of your way.

Nofmt attempts to achieve the following goals:

 - Correct indentation and spacing without introducing unnecessary newlines. \
   For example, the following is a formatted result:
    ```go
    results := []string{}
    for msg := range rx1chan { results = append(results, msg) }
    for msg := range rx2chan { results = append(results, msg) }
    for msg := range rx3chan { results = append(results, msg) }
    data, err := Process(results); if err != nil { return err }
    for !acknowledged { acknowledged = tx.Send(data) }
    ```

 - Code already formatted with `gofmt` should not be modified by `nofmt`. \
   For example, the following code should be valid and remain unchanged:
    ```go
    for !acknowledged {
        acknowledged = tx.Send(data)
    }
    ```

 - Idempotence. Formatting a file twice should produce the same results as formatting it once.

 - Full `goimports` functionality. Imports can be added and removed automatically.

## Installation

Install `go`
```bash
#--- Installation ---
VER=1.24.4
cd $HOME
rm -rf $HOME/goroot
mkdir -p $HOME/goroot $HOME/go/bin
wget "https://dl.google.com/go/go${VER}.linux-amd64.tar.gz"
tar -xf "go${VER}.linux-amd64.tar.gz" -C $HOME/goroot --strip-components 1 go
hash -r

echo 'export PATH="$HOME/bin:$HOME/go/bin:$HOME/goroot/bin:$PATH"' >> $HOME/.bashrc
```
Log out and back in again to ensure the PATH is updated.

Build modified `gofmt`, `goimports`, and `gopls`:
```bash
cd $HOME/src
git clone https://github.com/yqdv/nofmt.git
git clone https://github.com/golang/tools.git

# Backup the go/printer/nodes.go file to go/printer/nodes.go.orig
md5sum $HOME/goroot/src/go/printer/nodes.go
cp $HOME/goroot/src/go/printer/{nodes.go,nodes.go.orig}

# Patch go/printer/nodes.go to update linebreak()
cd $HOME/src/nofmt
patch $HOME/goroot/src/go/printer/nodes.go < inject/go/printer/nodes.go.patch

# Build "gofmt", "goimports", and "gopls" using the modified source
cd $HOME/goroot/src/cmd/gofmt
go build
cp gofmt $HOME/go/bin/gofmt-nofmt

# Build "goimports" using the modified source
cd $HOME/src/tools/cmd/goimports
go build
mv goimports $HOME/go/bin/goimports-nofmt

# Build "gopls" using the modified source
cd $HOME/src/tools/gopls
go build
mv gopls $HOME/go/bin/gopls-nofmt

# Backup existing binaries and create links to the new ones
if [[ -f "$HOME/go/bin/gofmt"         && ! -L "$HOME/go/bin/gofmt"         ]]; then cp $HOME/go/bin/{gofmt,gofmt.orig};             fi
if [[ -f "$HOME/go/bin/goimports"     && ! -L "$HOME/go/bin/goimports"     ]]; then cp $HOME/go/bin/{goimports,goimports.orig};     fi
if [[ -f "$HOME/go/bin/gopls"         && ! -L "$HOME/go/bin/gopls"         ]]; then cp $HOME/go/bin/{gopls,gopls.orig};             fi
if [[ -f "$HOME/goroot/bin/gofmt"     && ! -L "$HOME/goroot/bin/gofmt"     ]]; then cp $HOME/goroot/bin/{gofmt,gofmt.orig};         fi
if [[ -f "$HOME/goroot/bin/goimports" && ! -L "$HOME/goroot/bin/goimports" ]]; then cp $HOME/goroot/bin/{goimports,goimports.orig}; fi
if [[ -f "$HOME/goroot/bin/gopls"     && ! -L "$HOME/goroot/bin/gopls"     ]]; then cp $HOME/goroot/bin/{gopls,gopls.orig};         fi

ln -sf $HOME/go/bin/gofmt-nofmt     $HOME/go/bin/gofmt
ln -sf $HOME/go/bin/goimports-nofmt $HOME/go/bin/goimports
ln -sf $HOME/go/bin/gopls-nofmt     $HOME/go/bin/gopls
ln -sf $HOME/go/bin/gofmt-nofmt     $HOME/goroot/bin/gofmt
ln -sf $HOME/go/bin/goimports-nofmt $HOME/goroot/bin/goimports
ln -sf $HOME/go/bin/gopls-nofmt     $HOME/goroot/bin/gopls

hash -r

# Revert the source file patch
cp $HOME/goroot/src/go/printer/{nodes.go.orig,nodes.go}
```
