FROM --platform=${BUILDPLATFORM} golang:1.15-alpine AS build
WORKDIR /
ENV CGO_ENABLED=0
COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG OUTPUT

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /release/couchsport.back .
FROM scratch AS release
COPY --from=build /release/couchsport.back /