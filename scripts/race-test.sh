#!/bin/bash

port=5000
host=localhost
competitors=3
laps=2
minimal_lap_time_sec=15
xml=1

LANG=C
cmdname=`basename $0`
newtmpdir=`mktemp -d /tmp/${cmdname}.XXXXXX`
spool="$newtmpdir/spool"

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
antenna=`shuf -i 1-4 -n 1`
time=`date +"%Y/%m/%d %T.%3N"`
unixtime=`date +%s%3N`
sleep_time=`shuf -i 1-15 -n 1`

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
cat > ${spool}  <<EOF
10000802020000010000079${racer}, ${unixtime}, ${antenna}
EOF
fi

[ "${racer}" != "1" ] && echo
cat ${spool}
cat ${spool} | nc -q 0 ${host} ${port}


[ "${racer}" != "${competitors}" ] && read -p "Next rider in ${sleep_time} seconds...." -t ${sleep_time}

done

[ "${lap}" != "${laps}" ] && read -p  "Next lap in ${minimal_lap_time_sec} seconds..." -t ${minimal_lap_time_sec}

done
}

raceXML
