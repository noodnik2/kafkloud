#! /bin/bash

#
# E.g., run the tests using localhost using port forwarding; e.g.:
#  kubectl port-forward -n kafkloud pod/consumer-dpl-7857df4f9-bkh4l 30072:8072 &
#  kubectl port-forward -n kafkloud pod/producer-dpl-dc9fdd4cc-5jhlh 30000:8000 &
#

PRODUCER_HOST=localhost
PRODUCER_URL=${PRODUCER_HOST}:30000/produce/stream/key

RECEIVER_HOST=localhost
RECEIVER_URL=${RECEIVER_HOST}:30072/status

sendMessage() {
  if curl -d '{"text":"w1 w2"}' -H 'Content-Type: "application/json"' $PRODUCER_URL; then
#    echo "message sent"
    return 0
  else
    rc=$?
    echo "message not sent: $rc"
    return $rc
  fi
}

error() {
  echo $* 1>&2
}

fatal() {
  error $*
  exit 1
}

receiveMessage() {
  message=$(curl -s -H 'Content-Type: "application/json"' $RECEIVER_URL)
  echo $message
}


# send initial request
sendMessage || fatal "can't send initial message"
m1=$(receiveMessage)

# send subsequent request
sendMessage || fatal "can't send next message"
m2=$(receiveMessage)

# responses should differ
[ "$m1" == "$m2" ] && fatal "m1($m1) same as m2($m2)"

echo "test passes: m1($m1), m2($m2)"
