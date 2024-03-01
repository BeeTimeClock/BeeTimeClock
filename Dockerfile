FROM alpine:3.19 AS BUILD

RUN apk add --no-cache go yarn make

WORKDIR /build
COPY . /build

ARG VUE_APP_COMMIT=dirty
ENV VUE_APP_COMMIT=${VUE_APP_COMMIT}

RUN make all

FROM alpine:3.19
ENV PORT=8080
COPY --from=BUILD /build/beetimeclock /beetimeclock

ENTRYPOINT /beetimeclock
