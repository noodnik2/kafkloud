# Kafka Producer

## Conda Environment

This incarnation of a producer is written in Python, and uses `conda`
to manage its dependencies.  The `conda` environment was created with:

```shell
$ conda create -c conda-forge -n streamer-kafka-producer python-confluent-kafka
$ conda activate streamer-kafka-producer
```

## Python

The [kafka-producer.py](./kafka-producer.py) and [kafka-consumer.py](./kafka-consumer.py)
examples here were gleaned from [this](https://developer.confluent.io/get-started/python/)
tutorial.

After successfully starting the Kafka broker (using `docker-compose`) today,
I was able to produce & consume events:

In one shell:
```shell
$ conda activate streamer-kafka-producer
$ python kafka-consumer.py getting_started.ini
Waiting...
```

In another shell:
```shell
$ conda activate streamer-kafka-producer
$ python kafka-producer.py getting_started.ini 
Produced event to topic stream: key = jsmith       value = book        
Produced event to topic stream: key = sgarcia      value = t-shirts    
...
```

Quickly, the first shell noted receipt of the messages:
```shell
Consumed event from topic stream: key = sgarcia      value = book        
Consumed event from topic stream: key = eabara       value = gift card   
...
```


### See Also
- [Intro to Confluent Kafka Python Producer](https://www.geeksforgeeks.org/introduction-to-confluent-kafka-python-producer/)
- [confluent-kafka-python](https://github.com/confluentinc/confluent-kafka-python)
- [confluent-kafka](https://pypi.org/project/confluent-kafka/)
