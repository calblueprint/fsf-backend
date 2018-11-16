#
# Makefile for TCLink C library API.
# Copyright (c) 2001 TrustCommerce.
#

#### Tools

CC = gcc
AR = ar sr

#### Options

MYFLAGS = -Wall -O2 -fPIC
CFLAGS =  -I/usr/local/ssl/include $(MYFLAGS)
LDFLAGS =  -L/usr/local/ssl/lib -lssl -lcrypto -lnsl 

OBJS = tclink.o tctest.o mem.o validate.o
OBJS2 = tclink.o tcexample.o mem.o validate.o
OBJS3 = tclink.o mem.o validate.o

BIN = tctest
BIN2 = tcexample

OUTPUTLIBS = libtclink.a libtclink.so

#### Make targets

all: $(BIN) $(BIN2) dynamic static

tctest: Makefile $(OBJS)
	$(CC) -o $(BIN) $(OBJS) $(LDFLAGS)

tcexample: Makefile $(OBJS2)
	$(CC) -o $(BIN2) $(OBJS2) $(LDFLAGS)

dynamic: Makefile $(OBJS3)
	$(CC) -shared -o libtclink.so $(OBJS3) $(LDFLAGS)

static: Makefile $(OBJS3)
	$(AR) libtclink.a $(OBJS3)

install:
	cp libtclink.so /usr/lib
	cp tclink.h /usr/include

clean:
	rm -f $(OUTPUTLIBS) $(OBJS) $(OBJS2) $(BIN)

