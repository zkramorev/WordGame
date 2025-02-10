import redis

redis_client = redis.StrictRedis(host='localhost', port=6379, db=0)

with open("russian_nouns.txt", "r", encoding="utf-8") as file:
    for line in file:
        word = line.strip()
        redis_client.sadd("russian_words", word)

print("Слова загружены в Redis.")
