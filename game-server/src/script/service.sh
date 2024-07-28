gamePid=game.pid

usage="Usage: service.sh [start|stop]"

mkLogFolders()
{
  mkdir -p logger
}

#export GODEBUG=gctrace=1
export PORT=55555

start()
{
	nohup ./game --logtostderr=true -redis=127.0.0.1:6387 -db=127.0.0.1:27221 -log_dir=logger > logger/console.log 2>&1 &
	echo $! > game.pid
}

stop()
{
  if [ -f $gamePid ]; then
    if kill `cat $gamePid` > /dev/null 2>&1; then
      echo "Killed the game."
    fi
    rm -f $gamePid
  else
    echo 'game.pid not found'
  fi
}

case $1 in
  (start)
    mkLogFolders
    stop
    echo StartService.
    start
    ;;
  (stop)
    echo StopService.
    stop
    ;;
  (*)
    echo $usage
    exit 1
    ;;
esac

