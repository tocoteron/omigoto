FROM postgres:17

# Install dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends wget git ca-certificates build-essential postgresql-server-dev-17

# Install sqldef
RUN wget -O - https://github.com/sqldef/sqldef/releases/latest/download/psqldef_linux_amd64.tar.gz | tar xvz && \
    mv psqldef /usr/local/bin/

# Install pgvector extension
RUN cd /tmp && \
    git clone --branch v0.8.0 https://github.com/pgvector/pgvector.git && \
    cd pgvector && \
    make && \
    make install && \
    rm -rf /tmp/pgvector

# Set locale
RUN localedef -i ja_JP -c -f UTF-8 -A /usr/share/locale/locale.alias ja_JP.UTF-8
ENV LANG ja_JP.UTF-8
