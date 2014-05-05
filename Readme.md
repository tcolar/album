### Album

*********************************
Work In Progress, not usable yet.
*********************************

A lightweight standalone photo gallery server.

Features:
  - Lightweight, based flat files. No database or complex setup needed.
  - Responsive : Adapts to mobile devices and such and provides scaled images automatically.
  - .... TODO ....

### TODO:
* index albums
* index pics
* test indexer
* list albums
* list album pics
* serve album pics
* basic web frontend to serve responsive catalog.
- "Breacrumbs" / Album navifation
- Header / footer or "embedding" (spit out html ??)
- Make it obvious what's an album vs what's an image ?
- generate and serve thumbnails (say 200px?)
- generate and server scaled images (full, 500, 1000 ?)

- Use an interface for all index storage ops, so could easily replace with some "real" db later
- Sync JSON ops using chnannels to be concurent safe.

# admin
- admin login
- ability to update album meta (name, description, hidden etc..)
- Update index (whole or indivual album/pics)
- implement hiiden albums & pics -> still can be accesed if no url ?? -> can be seen in json file !!
- ability to select highlight
- ability to reorder items
- ability to scale images (original)
- ability to rotate images
- ability to upload image or zip of images

#API
- API's to retrieve albums & pics -> + to embed in ther site ?
- API's to modify data ?

#Other
- Stats (# of views) for images ?

