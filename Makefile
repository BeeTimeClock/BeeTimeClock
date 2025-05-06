.PHONY: ui server
ui: ui-modules
	cd ui; yarn quasar build

ui-modules:
	cd ui; yarn install

server:
	go build -o beetimeclock

develop-frontend:
	cd ui; \
	VUE_APP_BACKEND_ADDRESS=http://localhost:8085 yarn quasar dev

develop-backend: ui
	PORT=8085 DB_HOST=localhost DB_PORT=5432 DB_PASSWORD=verysecretpassword DB_USER=postgres DATABASE=postgres air

all: ui-modules ui server

clean:
	rm -Rf ui/dist tmp/ beetimeclock
