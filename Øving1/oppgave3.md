## What is concurrency? What is parallelism? What's the difference?
- Concurrency is the decomposability of a program. Concurrency means that the result of a computing should be the same even if the program is executed out of order.
- Paralellism is performing several computations is parallel.
- The difference is that splitting up a computation is paralellism and performing several independent calculations at once is concurrency
    
## Why have machines become increasingly multicore in the past decade?
- We are reaching hardware limits in heatdissapation and powerconsumtion, the solution to this is making many small cores.
## What kinds of problems motivates the need for concurrent execution? (Or phrased differently: What problems do concurrency help in solving?)
- Concurrency helps solve/model problems that happens simultaneous 
## Does creating concurrent programs make the programmer's life easier? Harder? Maybe both? (Come back to this after you have worked on part 4 of this exercise)
- a bit of both. Timing and sharing variables becomes an issue, but problems that are inherently paralell is easier to model
## What are the differences between processes, threads, green threads, and coroutines?
- a process is a program with executable code and one ore more threads
- a thread is a executing unit, managed by OS
- green threads is threads that is not managed by OS, but scheduled by runtime
- corutines is executing units scheduled by programmer/user
## Which one of these do pthread_create() (C/POSIX), threading.Thread() (Python), go (Go) create?
- pthread_create() creates a OS managed thread
- threading.Thread() creates a OS manages thread and made avaliable by GIL
- go creates a co-rutine
## How does pythons Global Interpreter Lock (GIL) influence the way a python Thread behaves?
- The GIL handles threads and only make one thread avaliable for OS for each interpreter spawned
## With this in mind: What is the workaround for the GIL (Hint: it's another module)?
- spawning more interpreters
## What does func GOMAXPROCS(n int) int change?
- distributes gorutines over more OS threads
