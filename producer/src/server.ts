import express, { Request, Response } from 'express';
import { Kafka, logLevel } from 'kafkajs';
import dotenv from 'dotenv';
import cors from 'cors';

// load the configuration
dotenv.config();

// set up the Kafka client
const KAFKA_BROKER_ADDRESS = process.env.KAFKA_BROKER!
const KAFKA_TRANSACTION_TIMEOUT = (process.env.KAFKA_TRANSACTION_TIMEOUT && parseInt(process.env.KAFKA_TRANSACTION_TIMEOUT))|| 30000
const SERVER_PORT = process.env.SERVER_PORT!

console.log(`using kafka broker(${KAFKA_BROKER_ADDRESS})`);
const kafka = new Kafka({brokers: [KAFKA_BROKER_ADDRESS], logLevel: logLevel.ERROR});
const producer = kafka.producer({
        allowAutoTopicCreation: true,
        transactionTimeout: KAFKA_TRANSACTION_TIMEOUT
    }
);

async function sendToKafka(topic, key, message) {
    await producer.connect();
    const promise = producer.send({
        topic: topic,
        messages: [{key: key, value: message}],
    });
    await promise.then(rm => {
        console.log(`RecordMetadata from Kafka is ${JSON.stringify(rm)}`);
    }).catch(r => {
        console.error(`Rejected Reason from Kafka is ${r}`);
    });
}

async function disconnectKafka() {
    await producer.disconnect();
}

// create & configure the server application
const app = express();
app.use(express.json({ strict: false }));
app.use(cors());

// set up the entry points

app.get('/_health', (req: Request, res: Response) => {
    res.send('Healthy!');
});

app.post('/produce/:topic/:key', cors(), (req: Request, res: Response) => {

    const topic = req.params["topic"];
    const key = req.params["key"];
    console.log(`delivering(${JSON.stringify(req.body)}) using topic(${topic}) and key(${key})`);

    let result = '';
    sendToKafka(topic, key, JSON.stringify(req.body))
    .then(() => {
        result = "delivered!";
    })
    .catch(t => {
        result = `not delivered: ${t}`;
    });

    res.send(result);
});

// start the server

console.log('starting service...');
exitOnSignal('SIGINT');
exitOnSignal('SIGTERM');
process.stdin.resume();

const server = app.listen(SERVER_PORT, () => {
    console.log(`The application is listening on port ${SERVER_PORT}!`);
});

function exitOnSignal(signal) {
    process.on(signal, reason => {
        console.log(`exiting due to(${reason})`);
        disconnectKafka()
        .then(() => { console.log(`kafka disconnected`); })
        .catch(e => { console.error(`error disconnecting from Kafka(${e})`); })
        .finally(() => {
            server.close(() => {
                console.log(`Server shutdown`);
                process.exit(1);
            });
        });
    });
}

