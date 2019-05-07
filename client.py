import socket
import multiprocessing
import time

def master():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('127.0.0.1', 8080))
    s.send(b'set-name::master\n')

    while True:
        try:
            data = s.recv(4096)
            print('[master] ', data)

            if data == b'game-master::master':
                time.sleep(1)
                s.send(b'set-response::resposta::dica\n')
            elif data ==  b'player-question::perguntaa':
                time.sleep(1)
                print('mestre está respondendo')
                s.send(b'master-response::true\n')
        except:
            return

def player1():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('127.0.0.1', 8080))
    s.send(b'set-name::player1\n')

    while True:
        try:
            data = s.recv(4096)
            print('[player1] ', data)

            if data == b'game-master::player1':
                time.sleep(1)
                s.send(b'set-response::resposta::dica\n')

            elif data == b'round_player::player1':
                time.sleep(1)
                print('mandando pergunta')
                s.send(b'player-question::perguntaa\n')

            elif data == b'master-response::true':
                print('mestre respondeu... mandando resposta')
                time.sleep(1)
                s.send(b'player-response::resposta\n')
        except:
            return

def player2():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('127.0.0.1', 8080))
    s.send(b'set-name::player2\n')

    while True:
        try:
            data = s.recv(4096)
            print('[player2] ', data)

            if data == b'round_player::player2':
                time.sleep(1)
                print('mandando pergunta do jogador 2')
                s.send(b'player-question::perguntaa\n')
            elif data == b'master-response::true':
                print('mestre respondeu... mandando resposta')
                time.sleep(3)
                s.send(b'player-response::resposta\n')
        except:
            return

if __name__ == '__main__':
    c1 = multiprocessing.Process(target=master)
    c1.start()
    time.sleep(1)
    c2 = multiprocessing.Process(target=player1)
    c2.start()
    time.sleep(1)
    c3 = multiprocessing.Process(target=player2)
    c3.start()

    c1.join()
    c2.join()
    c3.join()