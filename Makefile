.PHONY: ui server
ui: ui-modules
	cd ui; yarn quasar build

ui-modules:
	cd ui; yarn install

server:
	go build -o beetimeclock

develop-frontend:
	cd ui; \
	VUE_APP_BACKEND_ADDRESS=http://localhost:8080 yarn quasar dev

develop-backend: ui
	air

all: ui-modules ui server

clean:
	rm -Rf ui/dist tmp/ beetimeclock
