FROM python:3.12-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY agents/requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy agent code
COPY agents/ ./agents/

# Create directories for mounted secrets
RUN mkdir -p /etc/secrets/llm-keys

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import sys; sys.exit(0)" || exit 1

# Entrypoint: receives JOB_ID and TARGET_URLS from Argo
ENTRYPOINT ["python", "-m", "agents.entrypoint"]
