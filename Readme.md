### Album

*********************************
Work In Progress, not usable yet.
*********************************

A simple standalone photo gallery server & image processor.

Features:
  * Lightweight, based on flat files and in memory. No database or complex setup needed (for now).
  * Responsive : Adapts to mobile devices and such and provides scaled images & thumbnails automatically.
  * Albums & images with meta data and ordering
  * Image Scaling / Rotating / Padding service
  - Built-in admin API to manage content
  - REST API to use image server and content from external systems.

### Done:
* index albums
* index pics
* test indexer
* list albums
* list album pics
* serve album pics
* basic web frontend to serve responsive catalog.
* default image for albums with no / missing highlight
* generate and serve thumbnails (say 200px?)
* Move thumails to [root]/thumbs/ rather than spread all over the albums
* Pad the thubmnails
* Test Indexer basics
* Admin login / Very basic API Auth
* Make index store swappable (interface)
* Impement index store using KV store.

### TODO:
- 1) "Breacrumbs" / Album navigation
- 2) Header / footer or "embedding" (spit out html ??)
- 2) Make it obvious what's an album vs what's an image ?
- 2) generate and server scaled images (full/1440 - 1000 - 600) + 200 for thumbnail
  use https://github.com/scottjehl/picturefill ?
- 3) ability to scale images (original)
- Make file store an interface too so could substitute file system for any other io impl.
- Do image scaling in a go routine ? (would display not yet ready albums though)
- Use an interface for all index "storage" ops, so could easily replace with some "real" db later
- Sync JSON ops using channels to be concurent safe.

- Auth : Have a password in a config file or such
- Auth : Secure cookie store id

### TESTING :
- Test indexer changes / updates
- Test mage processing ? hum might be tricky excecpt maybe test against pre-made test images ?
- Test admin serices
- Test API's

# Admin features
- 2) ability to rotate images
- 2) ability to delete images
- 2) Update index (whole or indivual album/pics)
- 3) ability to select highlight
- 3) ability to update album meta (name, description, hidden etc..)
- 3) ability to reorder items
- implement hiiden albums & pics -> still can be accesed if no url ?? -> can be seen in json file !!
- ability to upload image or zip of images

#API features
- API's to retrieve albums & pics -> + to embed in other site ?
- API's to modify data ?
- Ability to do image ops in memory ad then stream them straight ? (No disk)

#Other features
- Stats (# of views) for images ?

