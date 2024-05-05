
PUBLIC_HTML := mike@www.meffie.org:/var/www/canton-tabletop.org/

.PHONY: help
help:
	@echo "make <target>"
	@echo "targets:"
	@echo "  preview  Start the hugo server and a browser"
	@echo "  deploy   Publish files on public site"
	@echo "  build    Build files for local check"
	@echo "  clean    Remove generated files"

.PHONY:
preview:
	(sleep 1 && xdg-open http://localhost:1313) &
	hugo server

.PHONY: clean deploy
deploy: build
	rsync -avz --delete public/ $(PUBLIC_HTML)

.PHONY: build
build:
	hugo

.PHONY: clean
clean:
	rm -rf public/ resources/
