main.o: main.cpp plus.hpp
	g++ -c main.cpp -o main.o -I.
plus.o: plus.cpp plus.hpp
	g++ -c plus.ccp -o plus.o -I.
all: main.o plus.o
	g++ main.o plus.o -o main
clear:
	rm *.o 