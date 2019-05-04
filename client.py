import socket
import multiprocessing
import time


def client1():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('127.0.0.1', 8080))
    s.send(b'set-name::carlos sumare')

    while True:
        try:
            data = s.recv(1024)
            print('c1', str(data))

            if data == b'game-master::carlos sumare':
                s.send(b'set-response::Pao::alimento')

        except:
            print('error')
            return


def other_clients(name):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('127.0.0.1', 8080))
    s.send(('set-name::' + name).encode())

    while True:
        try:
            data = s.recv(1024)
            print(name, str(data))

        except:
            print('error')
            return



if __name__ == '__main__':
    c1 = multiprocessing.Process(target=client1)
    c1.start()

    time.sleep(1)

    for i in range(5):
        c = multiprocessing.Process(target=other_clients, args=('client' + str(i),))
        c.start()
