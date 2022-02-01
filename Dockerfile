FROM golang:1.17 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG SERVICE_NAME=oneonone

WORKDIR /${SERVICE_NAME}/bin
COPY ./bin ./

RUN if [ "$BUILDPLATFORM" = "linux/amd64" ]; then mv ${SERVICE_NAME}.linux.amd64 ${SERVICE_NAME}; fi
RUN if [ "$BUILDPLATFORM" = "linux/arm64" ]; then mv ${SERVICE_NAME}.linux.arm64 ${SERVICE_NAME}; fi


FROM --platform=$BUILDPLATFORM gcr.io/distroless/base

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG SERVICE_NAME=oneonone

COPY --from=build /${SERVICE_NAME}/bin/${SERVICE_NAME} /app/server

EXPOSE 18081

ENTRYPOINT ["app/server"]
