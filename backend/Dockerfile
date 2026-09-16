FROM python:3.12.8-slim-bookworm

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1
WORKDIR /app
COPY requirements.txt ./
RUN pip install --no-cache-dir --requirement requirements.txt
COPY app ./app
EXPOSE 80
HEALTHCHECK --interval=10s --timeout=2s --retries=6 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://127.0.0.1/healthz', timeout=2)" || exit 1
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "80", "--no-access-log"]
