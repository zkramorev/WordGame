package main

import (
	pb "WordGame/wordCheckerPB"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type WordCheckerClient struct {
	conn   *grpc.ClientConn
	client pb.WordServiceClient
}

func NewWordCheckerClient(address string) (*WordCheckerClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	client := pb.NewWordServiceClient(conn)
	return &WordCheckerClient{conn: conn, client: client}, nil
}

func (wc *WordCheckerClient) Close() {
	if wc.conn != nil {
		wc.conn.Close()
	}
}

func (wc *WordCheckerClient) IsMoveCorrect(prevWord string, currWord string, wordsList map[string]int) (bool, error) {
	if _, exists := wordsList[currWord]; exists {
		return false, errors.New("Это слово уже использовалось в игре")
	}

	if len(prevWord) != 0 {
		runesPrevWord := []rune(prevWord)
		runesCurrWord := []rune(currWord)
		lastRunePrevWord := runesPrevWord[len(runesPrevWord)-1]
		if lastRunePrevWord == 'ъ' || lastRunePrevWord == 'ь' || lastRunePrevWord == 'ы' {
			lastRunePrevWord = runesPrevWord[len(runesPrevWord)-2]
		}
		if runesCurrWord[0] != lastRunePrevWord {
			fmt.Println(runesCurrWord[0], lastRunePrevWord)
			return false, errors.New("Первая буква вашего слова не совпадает с последней предыдущего слова :(")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	req := &pb.WordRequest{Word: currWord}
	res, err := wc.client.CheckWord(ctx, req)
	if err != nil {
		return false, errors.New("Ошибка при вызове сервиса проверки слов :(")
	}

	if !res.IsCorrect {
		return false, errors.New("Такого слова не знаю...")
	}
	return true, nil
}
