from concurrent import futures

import grpc

import redis
import word_pb2
import word_pb2_grpc

REDIS_HOST = "redis"
REDIS_PORT = 6379
REDIS_DB = 0


class WordService(word_pb2_grpc.WordServiceServicer):
    def __init__(self):
        self.redis_client = redis.StrictRedis(host=REDIS_HOST, port=REDIS_PORT, db=REDIS_DB)

    def CheckWord(self, request, context):
        word = request.word.strip().lower()
        print(f"Получен запрос: {word}", flush=True)
        if not word:
            return word_pb2.WordResponse(word=word, is_correct=False)

        is_present = self.redis_client.sismember("russian_words", word)
        return word_pb2.WordResponse(word=word, is_correct=is_present)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    word_pb2_grpc.add_WordServiceServicer_to_server(WordService(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    print("gRPC сервер запущен на порту 50051...", flush=True)
    server.wait_for_termination()


serve()
