#!/usr/bin/env zsh
#
#

local RUN_PATH=$(dirname $0:A)
local GOMATH="${GOPATH}/src/github.com/mlisa/gomath"
local NUM_COOR=${1}
local NUM_PEER=${2}

if [[ $(command -v terminator) ]]; then
  local TERMINAL="terminator"
else
  echo "No terminal found"
  exit -1
fi

rm ${RUN_PATH}/config_*.json

function generateCoordinatorConfig {
  local port=$(( 8000 + ${2} ))
cat << EOF > config_coordinator${1}.json
{
  "myself" : {
          "id" : "${1}",
          "address" : "127.0.0.1:${port}"
  }
}
EOF
}

for i in {1..${NUM_COOR}}; do
  generateCoordinatorConfig ${i} ${i}
done

function generatePeerConfig {
  local port=$(( 8100 + ${2} ))
cat << EOF > config_peer${1}.json
{
  "myself" : {
          "id" : "${1}",
          "address" : "127.0.0.1:${port}",
          "latency" : 100.0,
          "computationcapability" : 50.0,
          "queue" : 12.0
  }
}
EOF
}

for i in {1..${NUM_PEER}}; do
  generatePeerConfig ${i} ${i}
done

cd ${GOMATH}/peer/ && \
  go build -o peer *.go && \
  mv ${GOMATH}/peer/peer ${GOPATH}/bin/peer

cd ${GOMATH}/coordinator/ && \
  go build -o coordinator *.go && \
  mv ${GOMATH}/coordinator/coordinator ${GOPATH}/bin/coordinator

for i in {1..$(( ${NUM_COOR}-1 ))}; do
  ${TERMINAL} -e "coordinator -c ${RUN_PATH}/config_coordinator1.json" &
done
  ${TERMINAL} -e "coordinator -c ${RUN_PATH}/config_coordinator${NUM_COOR}.json &> /dev/null" &

for i in {1..$(( ${NUM_PEER}-1 ))}; do
  ${TERMINAL} -e "peer -c ${RUN_PATH}/config_peer${i}.json &> /dev/null" &
done
${TERMINAL} -e "peer -c ${RUN_PATH}/config_peer${NUM_PEER}.json &> /dev/null" &
