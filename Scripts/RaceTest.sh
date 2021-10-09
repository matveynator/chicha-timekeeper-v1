#!/bin/bash
# configuration options:
# tested under Linux and Mac OS X.
# TagStreamCustomFormat (Custom): %k, ${MSEC1}, %a
##################################################
port=4000
host=localhost
competitors=10 #riders
results=10 #results from one rider
laps=5 #laps
minimal_lap_time_sec=60
xml=1  #0 -> csv (%k, ${MSEC1}, %a), 1 -> xml
random=1 #0 = 1 2 3 4 5; #1 = 4 1 2 3 5
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

			#random race id:
      if [ "$random" == "1" ]
      then
        racer=$((RANDOM % ${competitors}))
      fi

			#random ammount of data from antenna
			if [ "$random" == "1" ]
			then
				iterations=$((RANDOM % ${results}))
			else
				iterations=${results}
			fi

			for result in `seq 1 ${iterations}`;
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
				else
					time=`date +"%Y/%m/%d %T.%3N"`
					unixtime=`date +%s%3N`
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

				#[ "${racer}" != "1" ] 
				cat ${spool}
				cat ${spool} | nc ${netcat_args} ${host} ${port}
			done

			[ "${racer}" != "${competitors}" ] && read -p "Next rider in ${sleep_time} seconds...." -t ${sleep_time}

		done

		[ "${lap}" != "${laps}" ] && echo ""; echo ""; read -p  "Next lap in ${minimal_lap_time_sec} seconds..." -t ${minimal_lap_time_sec}

	done
	echo ""
	echo "${cmdname} finished."
}

raceXML
