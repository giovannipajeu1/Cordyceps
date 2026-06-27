CC      = clang
CFLAGS  = -O2 -Wall -Wextra -mmacosx-version-min=12.0
LDFLAGS = -framework CoreFoundation
TARGET  = agent
SRCS    = main.c beacon.c shell.c persist.c evasion.c
OBJS    = $(SRCS:.c=.o)

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

clean:
	rm -f $(OBJS) $(TARGET)
