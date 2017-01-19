# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread
from threading import Lock

i = 0
mutex=Lock()

def ThreadFunction1():
    global i
    for num in range(0, 1000000):
        mutex.acquire()
        i = i + 1
        mutex.release()

# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")

def ThreadFunction2():
    global i
    for num in range(0, 1000000):
        mutex.acquire()
        i = i - 1
        mutex.release()

def main():
    Thread1 = Thread(target = ThreadFunction1, args = (),)
    Thread2 = Thread(target = ThreadFunction2, args = (),)
    Thread1.start()
    Thread2.start()
    Thread1.join()
    Thread2.join()
    print(i)


main()
