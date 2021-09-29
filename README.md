### Chicha timekeeper. 

Free chronograf for motorcycle, cars, bycicle and other types of competitions. 
UHF-RFID compatible. 

<img src="https://raw.githubusercontent.com/matveynator/chicha/main/chicha.jpg" width="600">

Бесплатная программа хронометража для любого трека на базе UHF-RFID 905-912 MHZ. 

- ## [↓ Download latest version of CHICHA.](http://files.matveynator.ru/chicha/) 
- ## [↓ Скачать последнюю версию CHICHA.](http://files.matveynator.ru/chicha/)

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

### Chicha example log:
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
