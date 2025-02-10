import redis
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI()
redis_client = redis.StrictRedis(host='localhost', port=6379, db=0)


class WordRequest(BaseModel):
    word: str


def is_correct_word(word: str) -> bool:
    """
    Возвращает True, если слово есть в Redis, иначе False.
    """
    is_present_in_database = redis_client.sismember("russian_words", word)
    if not is_present_in_database:
        print(f"Слово '{word}' не найдено в Redis.")
        return False
    return True


@app.post("/analyze/")
def analyze_word_endpoint(request: WordRequest):
    """
    Обрабатывает запрос и возвращает результат анализа.
    """
    word = request.word.strip().lower()
    if not word:
        raise HTTPException(status_code=400, detail="Слово не может быть пустым")

    return {"word": word, "is_correct": is_correct_word(word)}
