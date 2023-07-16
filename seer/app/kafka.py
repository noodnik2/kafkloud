import logging
from confluent_kafka import Consumer, Producer, KafkaException


KAFKA_TOPIC_STATEMENT = 'seer-statement'
KAFKA_TOPIC_QUESTION = 'seer-question'
KAFKA_TOPIC_ANSWER = 'seer-answer'


class KafkaSeer:

    def __init__(self, seer, broker):

        consumer_conf = {
            'group.id': 'seer',
            'bootstrap.servers': broker,
            'session.timeout.ms': 6000,
            'auto.offset.reset': 'earliest',
            'enable.auto.offset.store': False,
        }

        producer_conf = {
            'bootstrap.servers': broker,
        }

        self.logger = logging.getLogger(__name__)
        self.seer = seer
        self.cancel_flag = False

        self.consumer = Consumer(consumer_conf, logger=self.logger)
        self.producer = Producer(**producer_conf, logger=self.logger)

    def cancel(self):
        self.cancel_flag = True

    def run(self):

        try:
            consumer_topics = [KAFKA_TOPIC_STATEMENT, KAFKA_TOPIC_QUESTION]
            self.logger.debug(f"subscribing to topics({consumer_topics})")
            self.consumer.subscribe(consumer_topics, on_assign=self._notice_assignment)
            self.logger.debug("starting kafka consumer loop")
            while not self.cancel_flag:
                self._poll()
            self.logger.debug("ending kafka consumer loop")

        except KeyboardInterrupt:
            self.logger.debug('kafka consumer interrupted')

        finally:
            self.logger.debug(f"unsubscribing from consumer topics")
            self.consumer.unsubscribe()
            self.logger.debug(f"waiting for {len(self.producer)} outstanding delivery/ies")
            self.producer.flush()
            # commit final offsets
            self.consumer.close()
            self.logger.debug('kafka consumer stopped')

    def _poll(self):
        msg = self.consumer.poll(timeout=1.0)
        if msg is None:
            return
        if msg.error():
            raise KafkaException(msg.error())
        topic = msg.topic()
        data = msg.value()
        key = msg.key()
        partition = msg.partition()
        self.logger.debug(f"incoming message on topic({topic}) via partition({partition})"
                          f" at offset({msg.offset()}) with key({str(key)}), data({data})")
        if topic == "seer-statement":
            self.seer.accept([str(data)])
            self.logger.debug(f"seer accepted({data})")
        elif topic == "seer-question":
            seer_answers = self.seer.ask([str(data)])
            self.logger.debug(f"seer was asked({data}) and responded({seer_answers})")
            for seer_answer in seer_answers:
                self._deliver_answers(seer_answer)
        self.consumer.store_offsets(msg)

    def _deliver_answers(self, answer):
        self.logger.debug(f"delivering answer({answer}) to topic({KAFKA_TOPIC_ANSWER})")
        try:
            self.producer.produce(KAFKA_TOPIC_ANSWER, key="seer-answer", value=str(answer), callback=self._delivery_callback)

        except BufferError:
            self.logger.warning(f"producer queue is full ({len(self.producer)} message(s) awaiting delivery)")

        self.producer.poll(0)

    def _notice_assignment(self, consumer, partitions):
        self.logger.debug(f"partition assignment({partitions})")

    def _delivery_callback(self, err, msg):
        if err:
            self.logger.debug(f"message failed delivery: {err}")
        else:
            self.logger.debug(f"message delivered to {msg.topic()} [{msg.partition()}] @ {msg.offset()}")
