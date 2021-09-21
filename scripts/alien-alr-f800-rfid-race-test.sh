#!/bin/bash
# configuration options:
# tested under Linux and Mac OS X.
# TagStreamCustomFormat (Custom): %k, ${MSEC1}, %a
##################################################
port=4000
host=localhost
competitors=10
laps=10
minimal_lap_time_sec=60
xml=1   #0 -> csv (%k, ${MSEC1}, %a), 1 -> xml
##################################################

LANG=C
cmdname=`basename $0`
newtmpdir=`mktemp -d /tmp/${cmdname}.XXXXXX`
spool="$newtmpdir/spool"
os=`uname`
if [ "${os}" == "Linux" ]
then 
	netcat_args="-q 0"
elif [ "${os}" == "Darwin" ]
then
	netcat_args="-w 0"
fi

test
function cleanup () {
  rm -rf "${newtmpdir}"
}

trap 'cleanup' EXIT
trap 'cleanup' SIGTERM


function raceXML() {
for lap in `seq 1 ${laps}`;
do
[ "${lap}" != "1" ] && echo
echo "LAP #${lap}"

for racer in `seq 1 ${competitors}`; 
do
antenna=$((RANDOM % 4))
sleep_time=$((RANDOM % 10))

if [ "${os}" == "Linux" ]
then
	time=`date +"%Y/%m/%d %T.%3N"`
	unixtime=`date +%s%3N`
elif [ "${os}" == "Darwin" ]
then
	time=`date +"%Y/%m/%d %T.000"`
        unixtime=`date +%s000`
fi

if [ "${xml}" == "1" ]
then
cat > ${spool}  <<EOF
<Alien-RFID-Tag>
  <TagID>1000 0802 0200 0001 0000 079${racer}</TagID>
  <DiscoveryTime>${time}</DiscoveryTime>
  <LastSeenTime>${time}</LastSeenTime>
  <Antenna>${antenna}</Antenna>
  <ReadCount>1</ReadCount>
  <Protocol>2</Protocol>
</Alien-RFID-Tag>
EOF
else
echo "10000802020000010000079${racer}, ${unixtime}, ${antenna}" > ${spool}
fi

[ "${racer}" != "1" ] && echo
cat ${spool}
cat ${spool} | nc ${netcat_args} ${host} ${port}


[ "${racer}" != "${competitors}" ] && read -p "Next rider in ${sleep_time} seconds...." -t ${sleep_time}

done

[ "${lap}" != "${laps}" ] && read -p  "Next lap in ${minimal_lap_time_sec} seconds..." -t ${minimal_lap_time_sec}

done
}

raceXML
