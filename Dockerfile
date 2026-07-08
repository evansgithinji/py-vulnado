FROM python:3.11-slim

RUN apt-get update && apt-get install -y \
    iputils-ping \
    libxml2-utils \
    sqlite3 \
    curl \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app/ ./app/

RUN mkdir -p /app/data /app/uploads /app/files /app/images /app/static /app/exports /app/backups

EXPOSE 5001

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:5001/health || exit 1

CMD ["python", "-m", "app.main"]
