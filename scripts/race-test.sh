#!/bin/bash

port=5000
host=localhost
competitors=8
laps=5
minimal_lap_time_sec=45

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

echo "LAP #${lap}"

for racer in `seq 1 ${competitors}`; 
do
antenna=`shuf -i 1-4 -n 1`
time=`date +"%Y/%m/%d %T.%3N"`
sleep_time=`shuf -i 1-15 -n 1`
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

cat ${spool}
cat ${spool} | nc -q 0 ${host} ${port}
read -p "Continuing in ${sleep_time} seconds...." -t ${sleep_time}

done
sleep ${minimal_lap_time_sec}
done
}

raceXML
