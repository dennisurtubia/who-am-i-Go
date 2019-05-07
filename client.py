import socket
import multiprocessing
import time


def client1():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('127.0.0.1', 8080))
    s.send(b'set-name::carlos sumare\n')

    while True:
        try:
            data = s.recv(10000)
            if len(data) == 0:
                return
            print(data)

            if data == b'game-master::carlos sumare':
                time.sleep(2)
                s.send(b'set-response::Pao::alimento\n')

         

        except:
            print('error')
            return


def other_clients(name):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('127.0.0.1', 8080))
    s.send(('set-name::' + name + "\n").encode())

    while True:
        try:
            data = s.recv(1024)
            #print(name, str(data))

            if data == b'game-master::client0':
                time.sleep(2)
                s.send(b'set-response::Pao::alimento\n')

            if data == b'round-player::client0':
                s.send(b'player-question::o tata é fota?')


        except:
            print('error')
            return



if __name__ == '__main__':
    c1 = multiprocessing.Process(target=client1)
    c1.start()

    time.sleep(1)

    clients = [c1]

    for i in range(5):
        c = multiprocessing.Process(target=other_clients, args=('client' + str(i),))
        c.start()
        clients.append(c)

    for c in clients:
        c.join()

    print('fimm')
