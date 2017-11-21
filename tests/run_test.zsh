#!/usr/bin/env zsh
#
#

local RUN_PATH=$(dirname $0:A)
local GOMATH="${GOPATH}/src/github.com/mlisa/gomath"
local NUM_COOR=${1}
local NUM_PEER=${2}

if [[ $(command -v terminator) ]]; then
  local TERMINAL="terminator -e"
else
  local TERMINAL="${GOMATH}/tests/term.scpt"
fi

function generateCoordinatorConfig {
  local port=$(( 8000 + ${2} ))
cat << EOF > config_coordinator${1}.json
{
  "myself" : {
          "id" : "coordinator${1}",
          "address" : "127.0.0.1:${port}"
  }
}
EOF
}

function generatePeerConfig {
  local port=$(( 8100 + ${2} ))
cat << EOF > config_peer${1}.json
{
  "myself" : {
          "id" : "peer${1}",
          "address" : "127.0.0.1:${port}",
          "computecapability" : 50.0
  }
}
EOF
}

if [[ ${3} = "--generate" ]]; then
  rm -f ${RUN_PATH}/*.json
  rm -f /tmp/coordinator*.log /tmp/peer*.log 2> /dev/null
  curl -S -s -A "Go-http-client/1.1" http://gomath.duckdns.org:8080/generate.php\?c\=${NUM_COOR} > /dev/null
  echo "[i] Generating new coordinator configs"
  for i in {1..${NUM_COOR}}; do
    generateCoordinatorConfig ${i} ${i}
  done
  
  echo "[i] Generating new peer configs"
  for i in {1..${NUM_PEER}}; do
    generatePeerConfig ${i} ${i}
  done
fi

echo "[i] Compiling..."
cd ${GOMATH}/peer/ && \
  go build -i -o peer *.go && \
  mv ${GOMATH}/peer/peer ${GOPATH}/bin/peer

cd ${GOMATH}/coordinator/ && \
  go build -i -o coordinator *.go && \
  mv ${GOMATH}/coordinator/coordinator ${GOPATH}/bin/coordinator

if [[ ${NUM_COOR} -gt 1 ]]; then
  for i in {1..$(( ${NUM_COOR}-1 ))}; do
    eval "${TERMINAL} \"${GOPATH}/bin/coordinator -c ${RUN_PATH}/config_coordinator${i}.json >& /tmp/coordinator${i}.log\"" &
    sleep 0.2
  done
fi
eval "${TERMINAL} \"${GOPATH}/bin/coordinator -c ${RUN_PATH}/config_coordinator${NUM_COOR}.json >& /tmp/coordinator${NUM_COOR}.log\"" &
sleep 0.2
echo "[i] Coordinators running"

if [[ ${NUM_PEER} -gt 1 ]]; then
  for i in {1..$(( ${NUM_PEER}-1 ))}; do
    eval "${TERMINAL} \"${GOPATH}/bin/peer -c ${RUN_PATH}/config_peer${i}.json >& /tmp/peer${i}.log\"" &
    sleep 0.2
  done
fi
eval "${TERMINAL} \"${GOPATH}/bin/peer -c ${RUN_PATH}/config_peer${NUM_PEER}.json >& /tmp/peer${NUM_PEER}.log\"" &
echo "[i] Peers running"

sleep 0.5
\tail -f /tmp/coordinator*.log /tmp/peer*.log


#######################################################
# Generate file for mirror using bash
##!/usr/bin/env bash
#out="["
#for i in $(seq 1 ${1}); do
#  out="$out{
#  	\"id\": \"coordinator${i}\",
#  	\"address\": \"127.0.0.1:800${i}\"
#	},"
#done
#out="${out%?}]"
#echo $out
#######################################################
