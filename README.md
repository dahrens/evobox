evobox
======

This is just a playground for golang - learning concurrent concepts by implementing a
world where everything participant runs in its own goroutine...

dependencies
============

The backend if driven by gin, the golang dependencies are covered by Godeps...

* github.com/gin-gonic/gin
* github.com/gin-gonic/contrib/static
* github.com/Pallinder/go-randomdata

The Frontend uses several js and css libs, which get loaded using cdn's

* pixi.js
* jquery
* jquery datatables
* bootstrap3

The map was copied and modified from https://github.com/castled/Tilemap

install
=======

:code:`go get github.com/dahrens/evobox`

usage
=====

Make sure that you run the binary from within a folder where the assets are available.

.. code-block:: golang

	go install github.com/dahrens/evobox
	cd $GOPATH/src/github/dahrens/evobox
	evobox
