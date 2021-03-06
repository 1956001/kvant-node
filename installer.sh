#!/bin/sh

mkdir -p /kvant
cd /kvant
ver=`wget -O - https://raw.githubusercontent.com/kvant-development/kvant-node/master/bin_name 2>/dev/null`
echo $ver
wget https://github.com/kvant-development/kvant-node/raw/master/release/$ver
chmod +x ./$ver
./$ver show_node_id
wget -O genesis.json https://raw.githubusercontent.com/kvant-development/kvant-node/master/genesis/mainnet/genesis.json
cp genesis.json ~/.kvant/config/genesis.json

echo "screen -S node -d -m ./node_start.sh" > daemon.sh
echo "screen -r -d node" >> daemon.sh
chmod +x daemon.sh


mkdir -p /kvant/logs

echo > node_start.sh
echo "log=\"/kvant/logs/daemon.`hostname`.[$ver].`date '+%Y-%m-%d_%H_%M_%S'`.txt\"" >>  node_start.sh
echo "./$ver version" >> node_start.sh
echo "./$ver show_node_id" >> node_start.sh
echo "./$ver show_validator" >> node_start.sh
echo 'echo $log'  >> node_start.sh
echo "./$ver node >> \$log 2>&1 " >> node_start.sh
chmod +x node_start.sh
./node_start.sh

