FROM python:3.11

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY words_redis_storage.py .
COPY russian_nouns.txt .

CMD ["python", "words_redis_storage.py"]