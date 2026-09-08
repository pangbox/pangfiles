# Based on the cgofuse xgo Dockerfile by Bill Zissimopoulos.
# Modified to use crazymax/xgo as the base image.
# --
# Copyright (c) 2017-2020 Bill Zissimopoulos
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
# --
# https://github.com/billziss-gh/cgofuse/

FROM docker.io/alpine:3.24.1

RUN apk add --no-cache curl 7zip

ARG macfuse_version=4.1.2
RUN curl -L -o macfuse.dmg "https://github.com/osxfuse/osxfuse/releases/download/macfuse-${macfuse_version}/macfuse-${macfuse_version}.dmg" && \
    7zz e macfuse.dmg "macFUSE/Extras/macFUSE ${macfuse_version}.pkg" && \
    7zz e "macFUSE ${macfuse_version}.pkg" Core.pkg/Payload && \
    7zz e Payload && \
    (7zz x Payload~ -o/tmp || true) && \
    mkdir -p /usr/local/include/osxfuse && \
    cp -R /tmp/usr/local/include/fuse /usr/local/include/osxfuse && \
    ls -lah /tmp/usr/local/lib/ && \
    cp /tmp/usr/local/lib/libfuse.2.dylib /usr/local/lib/libosxfuse.dylib

ARG winfsp_version=2.0
RUN curl -L -o winfsp.zip "https://github.com/winfsp/winfsp/archive/v${winfsp_version}.zip" && \
    7zz e winfsp.zip "winfsp-${winfsp_version}/inc/fuse/*" -o/usr/local/include/winfsp

RUN apk add --no-cache 7zip

FROM docker.io/crazymax/xgo:1.25.9

RUN dpkg --add-architecture i386 && \
    dpkg --add-architecture amd64 && \
    apt-get update && \
    apt-get install -y --no-install-recommends libfuse-dev:i386 && \
    apt-get install -y --no-install-recommends libfuse-dev:amd64 && \
    apt-get download libfuse-dev:i386 && \
    dpkg -x libfuse-dev*i386*.deb / && \
    mkdir -p /usr/local/include && \
    mkdir -p /usr/local/lib

COPY --from=0 /usr/local/include/osxfuse /usr/local/include/osxfuse
COPY --from=0 /usr/local/lib/libosxfuse.dylib /usr/local/lib/libosxfuse.dylib
COPY --from=0 /usr/local/include/winfsp /usr/local/include/winfsp

ENV OSXCROSS_NO_INCLUDE_PATH_WARNINGS 1
