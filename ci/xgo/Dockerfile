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

FROM crazymax/xgo:1.16.6

# Install latest 7-zip for newer DMG/HFS support.
RUN wget https://www.7-zip.org/a/7z2102-linux-x64.tar.xz && \
    tar xvf 7z2102-linux-x64.tar.xz -C /usr/bin/

ARG macfuse_version=4.1.2
ARG winfsp_version=1.10B2
RUN wget -q -O macfuse.dmg "https://github.com/osxfuse/osxfuse/releases/download/macfuse-${macfuse_version}/macfuse-${macfuse_version}.dmg" && \
    7zz e macfuse.dmg "macFUSE/Extras/macFUSE ${macfuse_version}.pkg" && \
    7zz e "macFUSE ${macfuse_version}.pkg" Core.pkg/Payload && \
    7zz e Payload && \
    (7zz x Payload~ -o/tmp || true) && \
    mkdir /usr/local/include/osxfuse && \
    cp -R /tmp/usr/local/include/fuse /usr/local/include/osxfuse && \
    ls -lah /tmp/usr/local/lib/ && \
    cp /tmp/usr/local/lib/libfuse.2.dylib /usr/local/lib/libosxfuse.dylib

RUN dpkg --add-architecture i386 && \
    dpkg --add-architecture amd64 && \
    apt-get update && \
    apt-get install -y --no-install-recommends libfuse-dev:i386 && \
    apt-get install -y --no-install-recommends libfuse-dev:amd64 && \
    apt-get download libfuse-dev:i386 && \
    dpkg -x libfuse-dev*i386*.deb /

RUN wget -q -O winfsp.zip "https://github.com/billziss-gh/winfsp/archive/v${winfsp_version}.zip" && \
    7zz e winfsp.zip "winfsp-${winfsp_version}/inc/fuse/*" -o/usr/local/include/winfsp

ENV OSXCROSS_NO_INCLUDE_PATH_WARNINGS 1