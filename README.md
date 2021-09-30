### Chicha: the competition timekeeper (chronograph).

Free chronograf for motorcycle, cars, bycicle and other types of competitions. 
UHF-RFID compatible. 

<img src="https://raw.githubusercontent.com/matveynator/chicha/main/chicha.jpg" width="600">

Бесплатная программа хронометража для любого трека на базе UHF-RFID 860-960 MHz. 

- ## [↓ Download latest version of CHICHA.](http://files.matveynator.ru/chicha/) 
- ## [↓ Скачать последнюю версию CHICHA.](http://files.matveynator.ru/chicha/)

Supported OS: Mac, Linix, FreeBSD, DragonflyBSD, OpenBSD, NetBSD, Plan9, Windows
Supported architectures: x86-32, x86-64, ARM, ARM64. 
Supported databases: PostgreSQL (all platforms), SQLite (x86-32, x86-64 only).

###  supported readers / поддерживаемые считыватели: 

- [Alien ALR-F800](https://www.alientechnology.com/products/readers/alr-f800/)


### supported antennas / поддерживаемые антенны:

- [MT-263005/N](https://www.arcantenna.com/products/mt-263005-n-902-928-mhz-fcc-13-dbi-dbi-linear-v-h-polarity-directional-antenna-with-n-female-termination): 902-928 MHz FCC 13 dBi UHF-RFID antenna (MTI) 12++ meter range
- [CZS258](https://aliexpress.ru/item/32891562027.html) -  860-960 MHz 9dBi UHF-RFID antenna (aliexpress) 12 meter range


## Alien ALR-F800 RFID Reader settings:

```
http://IP/auth/tagstream.html
Login: alien
Pass: password
TagStreamAddress: IP:4000 
TagStreamCustomFormat (Custom): %k, ${MSEC1}, %a
```

### проверить приходит ли что нибудь на 4000 порт "apt-get install nc socat"
### Linux:
```
socat - TCP-LISTEN:4000,fork,reuseaddr 
nc -lvp 4000
```
### MAC OS X:
```
nc -l 0.0.0.0 4000
```


### TagStreamCustomFormat (Custom): %k, ${MSEC1}, %a
```

100008020200000100000189, 1622570553397, 3
100008020200000100000269, 1622570553478, 3
100008020200000100000269, 1622570553602, 3
100008020200000100000189, 1622570553611, 3
100008020200000100000268, 1622570553616, 3
100008020200000100000267, 1622570553635, 3
```

### ### TagStreamCustomFormat: XML
```
<Alien-RFID-Tag>
  <TagID>1000 0802 0200 0001 0000 0796</TagID>
  <DiscoveryTime>2021/05/16 12:00:34.730</DiscoveryTime>
  <LastSeenTime>2021/05/16 12:00:34.730</LastSeenTime>
  <Antenna>2</Antenna>
  <ReadCount>1</ReadCount>
  <Protocol>2</Protocol>
</Alien-RFID-Tag>
<Alien-RFID-Tag>
  <TagID>1000 0802 0200 0001 0000 0796</TagID>
  <DiscoveryTime>2021/05/16 12:00:34.823</DiscoveryTime>
  <LastSeenTime>2021/05/16 12:00:34.823</LastSeenTime>
  <Antenna>3</Antenna>
  <ReadCount>1</ReadCount>
  <Protocol>2</Protocol>
</Alien-RFID-Tag>
```

### chicha example log:
```
lap: 9, tag: 100008020200000100000793, position: 1, time: 1013000, gap: 0, best lap: 41000, start#: 1 
lap: 7, tag: 100008020200000100000794, position: 2, time: 978000, gap: 84000, best lap: 61000, start#: 3 
lap: 6, tag: 100008020200000100000797, position: 3, time: 804000, gap: 92000, best lap: 45000, start#: 2 
lap: 6, tag: 100008020200000100000791, position: 4, time: 883000, gap: 171000, best lap: 84000, start#: 6 
lap: 6, tag: 100008020200000100000798, position: 5, time: 983000, gap: 271000, best lap: 98000, start#: 7 
lap: 6, tag: 100008020200000100000795, position: 6, time: 1002000, gap: 290000, best lap: 104000, start#: 8 
lap: 6, tag: 100008020200000100000799, position: 7, time: 1004000, gap: 292000, best lap: 81000, start#: 5 
lap: 5, tag: 100008020200000100000790, position: 8, time: 900000, gap: 340000, best lap: 119000, start#: 10 
lap: 4, tag: 100008020200000100000796, position: 9, time: 986000, gap: 543000, best lap: 117000, start#: 9 
lap: 3, tag: 100008020200000100000792, position: 10, time: 464000, gap: 148000, best lap: 71000, start#: 4 
```

### chicha example json (/api/laps/results/byraceid/1):
```
[{"ID":69,"owner_id":0,"tag_id":"100008020200000100000792","discovery_unix_time":1632925602493,"discovery_time":"2021-09-29T17:26:42.493+03:00","antenna":2,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:26:43.860395+03:00","race_id":1,"current_race_postition":1,"time_behind_the_leader":0,"lap_number":8,"lap_time":105286,"lap_postition":1,"lap_is_current":1,"best_lap_time":84187,"best_lap_postition":5,"race_total_time":948606,"better_or_worse_lap_time":-21099},{"ID":61,"owner_id":0,"tag_id":"100008020200000100000799","discovery_unix_time":1632925492178,"discovery_time":"2021-09-29T17:24:52.178+03:00","antenna":2,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:24:53.201111+03:00","race_id":1,"current_race_postition":2,"time_behind_the_leader":110315,"lap_number":7,"lap_time":121365,"lap_postition":1,"lap_is_current":1,"best_lap_time":80163,"best_lap_postition":4,"race_total_time":838291,"better_or_worse_lap_time":-41202},{"ID":65,"owner_id":0,"tag_id":"100008020200000100000793","discovery_unix_time":1632925590346,"discovery_time":"2021-09-29T17:26:30.346+03:00","antenna":2,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:26:31.224135+03:00","race_id":1,"current_race_postition":3,"time_behind_the_leader":98168,"lap_number":7,"lap_time":66035,"lap_postition":3,"lap_is_current":1,"best_lap_time":34160,"best_lap_postition":1,"race_total_time":936459,"better_or_worse_lap_time":-31875},{"ID":66,"owner_id":0,"tag_id":"100008020200000100000794","discovery_unix_time":1632925592376,"discovery_time":"2021-09-29T17:26:32.376+03:00","antenna":3,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:26:33.220601+03:00","race_id":1,"current_race_postition":4,"time_behind_the_leader":100198,"lap_number":7,"lap_time":321863,"lap_postition":4,"lap_is_current":1,"best_lap_time":86166,"best_lap_postition":7,"race_total_time":938489,"better_or_worse_lap_time":-235697},{"ID":67,"owner_id":0,"tag_id":"100008020200000100000790","discovery_unix_time":1632925600436,"discovery_time":"2021-09-29T17:26:40.436+03:00","antenna":2,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:26:41.203934+03:00","race_id":1,"current_race_postition":5,"time_behind_the_leader":108258,"lap_number":7,"lap_time":124346,"lap_postition":5,"lap_is_current":1,"best_lap_time":86176,"best_lap_postition":8,"race_total_time":946549,"better_or_worse_lap_time":-38170},{"ID":60,"owner_id":0,"tag_id":"100008020200000100000796","discovery_unix_time":1632925485148,"discovery_time":"2021-09-29T17:24:45.148+03:00","antenna":1,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:24:46.204145+03:00","race_id":1,"current_race_postition":6,"time_behind_the_leader":214635,"lap_number":6,"lap_time":88227,"lap_postition":5,"lap_is_current":1,"best_lap_time":88227,"best_lap_postition":9,"race_total_time":831261,"better_or_worse_lap_time":0},{"ID":57,"owner_id":0,"tag_id":"100008020200000100000797","discovery_unix_time":1632925475059,"discovery_time":"2021-09-29T17:24:35.059+03:00","antenna":1,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:24:35.661412+03:00","race_id":1,"current_race_postition":7,"time_behind_the_leader":298796,"lap_number":5,"lap_time":73083,"lap_postition":7,"lap_is_current":1,"best_lap_time":68091,"best_lap_postition":2,"race_total_time":821172,"better_or_worse_lap_time":-4992},{"ID":54,"owner_id":0,"tag_id":"100008020200000100000791","discovery_unix_time":1632925399951,"discovery_time":"2021-09-29T17:23:19.951+03:00","antenna":0,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:23:20.202635+03:00","race_id":1,"current_race_postition":8,"time_behind_the_leader":342986,"lap_number":4,"lap_time":105221,"lap_postition":7,"lap_is_current":1,"best_lap_time":85176,"best_lap_postition":6,"race_total_time":746064,"better_or_worse_lap_time":-20045},{"ID":63,"owner_id":0,"tag_id":"100008020200000100000798","discovery_unix_time":1632925506228,"discovery_time":"2021-09-29T17:25:06.228+03:00","antenna":2,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:25:07.202575+03:00","race_id":1,"current_race_postition":9,"time_behind_the_leader":449263,"lap_number":4,"lap_time":100227,"lap_postition":9,"lap_is_current":1,"best_lap_time":75142,"best_lap_postition":3,"race_total_time":852341,"better_or_worse_lap_time":-25085},{"ID":68,"owner_id":0,"tag_id":"100008020200000100000795","discovery_unix_time":1632925602464,"discovery_time":"2021-09-29T17:26:42.464+03:00","antenna":2,"antenna_ip":"192.168.96.15","created_at":"2021-09-29T17:26:43.204096+03:00","race_id":1,"current_race_postition":10,"time_behind_the_leader":545499,"lap_number":4,"lap_time":124344,"lap_postition":10,"lap_is_current":1,"best_lap_time":97264,"best_lap_postition":10,"race_total_time":948577,"better_or_worse_lap_time":-27080}]
```


### Download chicha source code:
```
cd ~
git clone https://github.com/matveynator/chicha.git
export GOPATH=~/chicha/GOPATH
echo "export GOPATH=~/chicha/GOPATH" >> ~/.bash_profile
cd ~/chicha
```

### Run chicha (test):
``` 
cd ~/chicha
go run chicha.go
```

### Compile chicha:
```
cd ~/chicha
sh  Scripts/crosscompile.sh
ls ~/chicha/downloads
```

### Run race test (edit options inside: Scripts/RaceTest.sh)/
```
cd ~/chicha
sh Scripts/RaceTest.sh
```
 

