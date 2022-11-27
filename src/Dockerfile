#FROM python:3.6-alpine@sha256:3a6163bba1b54c1eb60e94551be93a8cf4c6dc09dccfbfbb4bd245d90bd3c281
FROM python:3.9

RUN pip3 install --no-cache-dir cryptography==37.0.4 && \
    pip3 install --no-cache-dir dvc[ssh]==2.19.0 && \
    pip3 install --no-cache-dir dvc[s3]==2.19.0

WORKDIR /app
ADD validator /app/validator

ENTRYPOINT ./validator
