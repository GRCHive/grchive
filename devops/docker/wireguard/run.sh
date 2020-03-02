#/bin/bash
wg-quick up wg0
sleep infinity &
wait $!
