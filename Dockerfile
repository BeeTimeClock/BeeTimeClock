FROM alpine:3.18 AS BUILD

RUN apk add --no-cache go yarn make

WORKDIR /build
COPY . /build

RUN make all

FROM alpine:3.18
ENV PORT=8080
COPY --from=BUILD /build/beetimeclock /beetimeclock

ENTRYPOINT /beetimeclock
