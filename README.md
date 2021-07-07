# Noxss

Noxss is a very basic implementation of the fundamental principle explained in the [Noxes paper](https://sites.cs.ucsb.edu/~vigna/publications/2006_kirda_kruegel_vigna_jovanovic_SAC.pdf) by Engin Kirda, Christopher Kruegel, Giovanni Vigna, and Nenad Jovanovic.
I built it as a proof of concept for a homework in IT Security to demonstrate the basic concepts. 

## What this can do
- act as a http proxy between the browser and the rest of the internet
- inspect incoming responses and build a list of one-time allowed requests (since static external links are considered safe by the authors of the paper)
    - img sources only, but extension to _url_ and _href_ is trivial
- automatically allow normal opening of websites (empty referer)
- automatically allow local links (compares referer and destination)
- asking a user whether a request matching none if the rules above should be allowed or blocked

## Limitations
- http only (no SSL support)
- loading multiple sites in the same session will lead to weird behavior because there is no isolation, but it should be good enough for simple demos

## Dependencies

- tested on linux with firefox
- [goproxy](https://github.com/elazarl/goproxy) is used for the proxy stuff
- [Zenity](https://help.gnome.org/users/zenity/stable/) is used to display user confirmation and must be installed (change the zenityBin constant if your executable is not in /usr/bin/zenity)
- Go for building

## Usage

```
go run main.go
```

Set the http proxy of your browser to localhost:8080