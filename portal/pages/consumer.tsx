import React, {useState, useEffect} from 'react';
import axios from 'axios';
import Head from "next/head";
import Button from "@/components/Button";
import {useRecoilState} from "recoil";
import {itemWithID} from "@/components/CommonProps";
import LogWindow from "@/components/LogWindow";
import MessageProducer from "@/components/MessageProducer";
import MessageConsumer from "@/components/MessageConsumer";

export default function Consumer() {

    // recoil items
    const consumerLogID = 'consumerLog'
    const statusLogID = 'statusText'

    const textAreaClassName = "p-2 overflow-y-auto w-full h-24 min-h-full"

    const [, setStatusText] = useRecoilState(itemWithID(statusLogID))
    const [, setConsumerText] = useRecoilState(itemWithID(consumerLogID))

    const [consumerUpdateCount, setConsumerUpdateCount] = useState(0);

    const updateStatus = (newLogEntry: string) => {
        setStatusText(currentLog => [...currentLog, newLogEntry]);
    }

    const updateConsumerLog = (newLogEntry: string) => {
        setConsumerText(currentLog => [...currentLog, newLogEntry]);
    }

    const refreshConsumer = () => {
        const url = `/api/status`;
        console.log(`get(${url})`)
        axios.get(url)
            .then((response) => {
                updateStatus(JSON.stringify(response))
                updateConsumerLog(JSON.stringify(response.data));
            })
            .catch((error) => {
                const axiosErrorStatus = `axiosErrorStatus(${error})`;
                console.error(axiosErrorStatus);
                updateStatus(axiosErrorStatus)
            });

    }

    useEffect(() => () => clearTimeout(setTimeout(refreshConsumer, 2000)), [consumerUpdateCount]);

    const postToProducer = (payload: object) => {

        const producerRequest = JSON.stringify(payload)

        updateStatus(`postToProducer(${producerRequest})`);

        const axiosHeaders = {
            headers: {
                'Content-Type': 'application/json'
            }
        }

        const url = `/api/outbox`;
        console.log(`post(${url})`)
        axios.post(url, producerRequest, axiosHeaders)
            .then((response) => {
                const producerRequestResponse = `producerRequestResponseStatusText(${response.statusText})`;
                updateStatus(producerRequestResponse);
                console.log(`triggering a consumer update`);
                setConsumerUpdateCount(consumerUpdateCount + 1);
            })
            .catch((error) => {
                const producerRequestError = `producerRequestError(${error})`;
                console.error(producerRequestError);
                updateStatus(producerRequestError);
            });

    }

    /**
        deliver() delivers the message to the producer
     */
    const deliver = (topic: string, message: string) => {
        updateStatus(`deliver to topic(${topic}): "${message}"`);
        postToProducer({
            topic: topic,
            key: new Date().toISOString(),
            message: message
        })
    }

    /**
        startConsumer() stops the current consumer process (if any) and starts
                        a new one, listening on the indicated topics
     */
    const startConsumer = (topics: string[]) => {
        // TODO do what the contract says above; must wait until we have an API or Kafka client (preferred / suggested)
        updateStatus(`startConsumer(${topics}); TODO implement`)
        refreshConsumer()
    }

    const title = (text: string) => {
        return (
            <h1 className="p-2 text-2xl decoration-2 font-bold">{text}</h1>
        );
    }

    return (
        <main className="flex items-center justify-around font-sans">
            <Head>
                <title>Consumer</title>
            </Head>

            <div className="p-2 m-5 w-full shadow-2xl drop-shadow-lg space-y-3">
                <div>
                    <Button name="test" onClick={() => {
                        setStatusText(oldArray => [...oldArray, 'testing setting the status text']);
                    }} />
                </div>
                <div className="border-2 border-black bg-blue-50">
                    {title("Status")}
                    <LogWindow
                        loggerId={statusLogID}
                        loggerDescription="Consumer status updates"
                        textAreaClassName={textAreaClassName}
                    />
                </div>
                <div className="border-2 border-black bg-blue-50">
                    {title("Messages")}
                    <div className="box-border border-dashed border-2">
                        <MessageProducer
                            knownTopics={['stream', 'seer-statement', 'seer-question']}
                            selectedTopic='stream'
                            textAreaClassName={textAreaClassName}
                            onDeliver={deliver}
                        />
                    </div>
                    <div className="box-border border-dashed border-2">
                        <MessageConsumer
                            knownTopics={['stream', 'seer-answer']}
                            selectedTopics={['stream', 'seer-answer']}
                            textAreaClassName={textAreaClassName}
                            outputItemID={consumerLogID}
                            onStartConsumer={startConsumer}
                        />
                    </div>
                </div>
            </div>

        </main>
    );
}

