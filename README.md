evobox
======

This is just a playground for golang and phaser.

dependencies
------------

The backend if driven by gin, the golang dependencies are covered by Godeps...

* https://github.com/gin-gonic/gin
* https://github.com/gin-gonic/contrib/static
* https://github.com/Pallinder/go-randomdata

The Frontend uses several js and css libs, which get loaded using cdn's

* http://phaser.io/
* http://jquery.com/
* http://www.datatables.net/
* http://getbootstrap.com/

artwork
-------

Since i am not good at graphic design i rely on open source published content
made by others. Thank for your sharing your Work.

Artists

  * http://www.lostgarden.com Daniel Cook
  * https://openclipart.org/user-detail/Keistutis
  * https://openclipart.org/user-detail/lemmling

Tools

* http://www.gimp.org
* http://inkscape.org
* http://www.synfig.org/
* http://www.codeandweb.com/texturepacker

install
-------

```
go get github.com/dahrens/evobox
```

usage
-----

Make sure that you run the binary from within a folder where the assets are available.

```
go install github.com/dahrens/evobox
cd $GOPATH/src/github/dahrens/evobox
evobox
```

