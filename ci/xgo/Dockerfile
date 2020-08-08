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

FROM crazymax/xgo:1.14.x

RUN dpkg --add-architecture i386 && \
    dpkg --add-architecture amd64 && \
    apt-get update && \
    apt-get install -y --no-install-recommends p7zip-full

ARG osxfuse_version=3.8.3
ARG winfsp_version=1.8B1

RUN wget -q -O osxfuse.dmg "https://github.com/osxfuse/osxfuse/releases/download/osxfuse-${osxfuse_version}/osxfuse-${osxfuse_version}.dmg" && \
    7z e osxfuse.dmg "FUSE for macOS/Extras/FUSE for macOS ${osxfuse_version}.pkg" && \
    7z e "FUSE for macOS ${osxfuse_version}.pkg" Core.pkg/Payload && \
    7z e Payload && \
    7z x Payload~ -o/tmp && \
    cp -R /tmp/usr/local/include/osxfuse /usr/local/include && \
    ls -lah /tmp/usr/local/lib/ && \
    cp /tmp/usr/local/lib/libosxfuse.2.dylib /usr/local/lib/libosxfuse.dylib

RUN apt-get update && \
    apt-get install -y --no-install-recommends libfuse-dev:i386 && \
    apt-get install -y --no-install-recommends libfuse-dev:amd64 && \
    apt-get download libfuse-dev:i386 && \
    dpkg -x libfuse-dev*i386*.deb /

RUN wget -q -O winfsp.zip "https://github.com/billziss-gh/winfsp/archive/v${winfsp_version}.zip" && \
    7z e winfsp.zip "winfsp-${winfsp_version}/inc/fuse/*" -o/usr/local/include/winfsp

ENV OSXCROSS_NO_INCLUDE_PATH_WARNINGS 1