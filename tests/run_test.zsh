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
          "latency" : 100.0,
          "computationcapability" : 50.0,
          "queue" : 12.0
  }
}
EOF
}

if [[ ${3} = "--generate" ]]; then
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
    eval "${TERMINAL} \"${GOPATH}/bin/coordinator -c ${RUN_PATH}/config_coordinator${i}.json\"" &
  done
fi
eval "${TERMINAL} \"${GOPATH}/bin/coordinator -c ${RUN_PATH}/config_coordinator${NUM_COOR}.json\"" &
echo "[i] Coordinators running"

if [[ ${NUM_PEER} -gt 1 ]]; then
  for i in {1..$(( ${NUM_PEER}-1 ))}; do
    eval "${TERMINAL} \"${GOPATH}/bin/peer -c ${RUN_PATH}/config_peer${i}.json >& /dev/null\"" &
  done
fi
eval "${TERMINAL} \"${GOPATH}/bin/peer -c ${RUN_PATH}/config_peer${NUM_PEER}.json >& /dev/null\"" &
echo "[i] Peers running"
