########################################################################
#### Example 9009:
#### PORT=9009 ./better.sh build
#### PORT=9009 ./better.sh start
#### PORT=9009 ./better.sh restart
#### PORT=9009 ./better.sh stop
########################################################################

PORT="${PORT:-9008}"
REDIS_HOST="${REDIS_HOST:-127.0.0.1:6387}"
MONGO_HOST="${MONGO_HOST:-mongodb://poker:siwS*Hies234@112.126.57.221:87}"
MONITOR_PORT="${MONITOR_PORT:-9109}"
ENV="${ENV:-develop}"

NAME="game${PORT}"
CONSOLE_LOG="./logger/console${PORT}.log"
LOG_DIR=./logger
PID_FILE="./game${PORT}.pid"
TOPLEVEL=$(git rev-parse --show-toplevel 2> /dev/null)

prestart()
{
    mkdir -pv ${TOPLEVEL}/src/script/logger ${TOPLEVEL}/src/pre-result ${TOPLEVEL}/src/html-result
}

build()
{
    GOPATH=/data/MissPoker go build -ldflags="-s -w" -o ./$NAME ../game/main.go
}

buildd()
{
    GOPATH=/data/MissPoker go build -o ./$NAME ../game/main.go
}

start()
{
    nohup ./$NAME --logtostderr=true -env=$ENV -redis=$REDIS_HOST -db=$MONGO_HOST --bindHost=:$PORT --monitorPort=:$MONITOR_PORT -log_dir=$LOG_DIR > $CONSOLE_LOG 2>&1 &
    echo $! > "$PID_FILE"
}

startd()
{
    GODEBUG=gctrace=1 nohup ./$NAME --logtostderr=true -env=$ENV -redis=$REDIS_HOST -db=$MONGO_HOST --bindHost=:$PORT --monitorPort=:$MONITOR_PORT -log_dir=$LOG_DIR > $CONSOLE_LOG 2>&1 &
    echo $! > "$PID_FILE"
}

stop()
{
    if [ -f "$PID_FILE" ]; then
        kill -15 $(cat "$PID_FILE") &> /dev/null
        if [ $? != 0 ]; then
            echo -e "\033[1;31mFail to kill\033[0m"
            exit
        fi
        sleep 4
        if ps aux | grep "$NAME" | grep -q -v grep; then
            echo -e "\033[1;31m${NAME}Î´Õý³£¹Ø±Õ\033[0m"
            exit
        fi
        rm -f "$PID_FILE"
    else
        echo -e "\033[1;31m${PID_FILE} Not Found\033[0m"
    fi
}

case $1 in
    build)
        build
        ;;
    buildd)
        buildd
        ;;
    start)
        prestart
        start
        ;;
    startd)
        prestart
        startd
        ;;
    stop)
        stop
        ;;
    restart)
        prestart
        stop
        start
        ;;
    restartd)
        prestart
        stop
        startd
        ;;
    *)
        ;;
esac
