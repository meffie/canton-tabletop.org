
BGG_USERS=Gidcumb73 JTCrawfo m3ff13 rsebrell tedstriker

.PHONY: help
help:
	@echo "make <target>"
	@echo "targets:"
	@echo "  post T='title string'  Create a new markdown file"
	@echo "  preview  Start the hugo server and a browser"
	@echo "  deploy   Publish files on public site"
	@echo "  build    Build files for local check"
	@echo "  games    Update game collection"
	@echo "  clean    Remove generated files"

.PHONY: post
post:
	@test ! -z "$(T)" || (echo "usage: make post T='<title string>'" >&2; exit 1)
	year=`date +'%Y'`; mm_dd=`date +'%m-%d'`; title=`echo $(T) | tr ' ' '-'`; \
	hugo new content posts/$${year}/$${mm_dd}__$${title}.md

.PHONY: preview
preview: clean
	(sleep 1 && xdg-open http://localhost:1313) &
	hugo --baseURL=http://localhost:1313 server

.PHONY: deploy
deploy: clean build
	rsync -avz --delete public/ $(PUBLIC_HTML)

.PHONY: build
build: data/games.json
	hugo

.PHONY: games
games: bin/bgg-export
	bin/bgg-export $(BGG_USERS) >data/games.json

data/games.json: bin/bgg-export Makefile
	bin/bgg-export $(BGG_USERS) >data/games.json

bin/bgg-export: src/bgg-export.go
	mkdir -p bin
	go build -o bin/bgg-export src/bgg-export.go

.PHONY: clean
clean:
	rm -rf public/ resources/
