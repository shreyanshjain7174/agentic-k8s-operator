# Stage 1: Builder
FROM python:3.12-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base

COPY agents/requirements.txt .

RUN pip install --no-cache-dir --target=/app/deps -r requirements.txt

# Stage 2: Runtime
FROM python:3.12-alpine AS runtime

WORKDIR /app

RUN apk add --no-cache curl \
    && rm -rf /var/cache/apk/*

RUN addgroup --gid 1000 appgroup \
    && adduser --uid 1000 --ingroup appgroup --disabled-password --no-create-home appuser

COPY --from=builder --chown=1000:1000 /app/deps /app/deps

COPY --chown=1000:1000 agents/ ./agents/

RUN mkdir -p /etc/secrets/llm-keys \
    && chown -R 1000:1000 /app /etc/secrets/llm-keys \
    && chmod -R 750 /app \
    && chmod 700 /etc/secrets/llm-keys

ENV PYTHONPATH=/app/deps:$PYTHONPATH \
    PATH=/app/deps/bin:$PATH

USER 1000:1000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import sys; sys.exit(0)" || exit 1

ENTRYPOINT ["python", "-m", "agents.entrypoint"]
