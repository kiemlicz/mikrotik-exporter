FROM python:3.11-slim-bookworm

WORKDIR /mikrotik-exporter
COPY
RUN pip install ./
