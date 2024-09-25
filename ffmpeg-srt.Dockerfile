FROM ubuntu:latest
ENV LD_LIBRARY_PATH=/usr/local/lib:/usr/local/lib64
WORKDIR /
RUN apt-get update -qq && \
    apt-get upgrade -y && \
    apt-get -y install \
    tclsh \
    libssl-dev \
    build-essential \
    autoconf \
    automake \
    build-essential \
    cmake \
    git-core \
    libass-dev \
    libfreetype6-dev \
    libgnutls28-dev \
    libmp3lame-dev \
    libvorbis-dev \
    libtool \
    libx264-dev \
    libx265-dev \
    libnuma-dev \
    libvpx-dev \
    libfdk-aac-dev \
    libopus-dev \
    meson \
    ninja-build \
    pkg-config \
    texinfo \
    wget \
    yasm \
    nasm \
    zlib1g-dev \
    libunistring-dev \
    libaom-dev \
    libsvtav1-dev \
    libsvtav1enc-dev \
    libsvtav1dec-dev \
    libdav1d-dev
RUN mkdir /ffmpeg_sources
RUN cd /ffmpeg_sources && \
    git clone https://git.ffmpeg.org/ffmpeg.git ffmpeg
RUN cd /ffmpeg_sources && \
    git clone https://github.com/Haivision/srt.git srt
RUN cd /ffmpeg_sources/srt && \
    ./configure && \
    make && \
    make install
RUN cd /ffmpeg_sources/ffmpeg && \
    ./configure \
    --prefix=/usr/local \
    --enable-gpl \
    --enable-gnutls \
    --enable-libaom \
    --enable-libass \
    --enable-libfdk-aac \
    --enable-libfreetype \
    --enable-libmp3lame \
    --enable-libopus \
    --enable-libsvtav1 \
    --enable-libdav1d \
    --enable-libvorbis \
    --enable-libvpx \
    --enable-libx264 \
    --enable-libx265 \
    --enable-nonfree \
    --enable-libsrt \
    --enable-shared && \
    make && \
    make install && \
    hash -r
RUN which ffmpeg
RUN ffmpeg -version
ENTRYPOINT ["ffmpeg"]