FROM golang:1.20-buster as agentbuild
WORKDIR /app
ADD . ./
RUN go get -d -v ./...
RUN cd agent && go build -o /go/bin/agent

FROM python:3.8-alpine AS pyspybuild
# only works with pip version 20.2.4
# https://github.com/benfred/py-spy/issues/353
RUN pip install pip==20.2.4
RUN echo 'manylinux1_compatible = True' > /usr/local/lib/python3.8/site-packages/_manylinux.py
RUN pip3 install py-spy==0.3.14

FROM bitnami/minideb as asyncprofiler
RUN install_packages curl tar ca-certificates
RUN curl -o async-profiler-2.9-linux-x64.tar.gz -L \
    https://github.com/jvm-profiling-tools/async-profiler/releases/download/v2.9/async-profiler-2.9-linux-x64.tar.gz
RUN tar -xvf async-profiler-2.9-linux-x64.tar.gz && mv async-profiler-2.9-linux-x64 async-profiler


FROM bitnami/minideb as nodejsbuild
RUN install_packages linux-perf git perl ca-certificates
RUN git clone https://github.com/brendangregg/FlameGraph

FROM bitnami/minideb:stretch
WORKDIR  /app
COPY --from=agentbuild /go/bin/agent /app/agent
#Copy for python profiler
COPY --from=pyspybuild /usr/local/bin/py-spy /app/py-spy
#Copy for jvm profiler
COPY --from=asyncprofiler /async-profiler /app/async-profiler
#Copy for perf profiler
COPY --from=nodejsbuild /FlameGraph /app/FlameGraph
COPY --from=nodejsbuild /usr/bin/perf /app/perf
COPY --from=nodejsbuild /usr/bin/perl /usr/bin/perl

CMD [ "/app/agent" ]