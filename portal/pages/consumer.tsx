import React, {useState, useEffect, useCallback} from 'react';
import axios from 'axios';
import Head from "next/head";
import Button from "@/components/Button";
import {useRecoilState} from "recoil";
import {itemWithID, toLocalLogTime} from "@/components/CommonProps";
import LogWindow from "@/components/LogWindow";
import MessageProducer from "@/components/MessageProducer";
import MessageConsumer from "@/components/MessageConsumer";

const Consumer = (): JSX.Element => {

    // recoil items
    const consumerLogID = 'consumerLog'
    const statusLogID = 'statusText'

    const textAreaClassName = "p-2 overflow-y-auto w-full h-24 min-h-full"

    const [, setStatusText] = useRecoilState(itemWithID(statusLogID))
    const [, setConsumerText] = useRecoilState(itemWithID(consumerLogID))

    const [consumerUpdateCount, setConsumerUpdateCount] = useState(0);

    const timestamp = (s: string) => {
        return `${toLocalLogTime(new Date())} ${s}`
    }

    const updateStatus = (newLogEntry: string) => {
        setStatusText(currentLog => [...currentLog, timestamp(newLogEntry)]);
    }

    const updateConsumerLog = (newLogEntry: string) => {
        setConsumerText(currentLog => [...currentLog, timestamp(newLogEntry)]);
    }

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

    const runConsumer = async (topics: string[]) => {

        // Relevant:
        // - For XMLHttpRequest2:
        //   - https://streams.spec.whatwg.org/#other-specs-rs-reading
        //   - https://streams.spec.whatwg.org/#read-loop
        const url = `/api/consume?topics=${topics}`;

        // see:
        // - https://developer.chrome.com/articles/fetch-streaming-requests/

        const backendResponse = await fetch(url);
        if (!backendResponse.ok) {
            updateStatus(`fetch error: ${backendResponse.status} ${backendResponse.statusText}`)
            return
        }
        if (!backendResponse.body) {
            updateStatus(`error: no body from fetch to backend`)
            return
        }
        const reader = backendResponse.body
            .pipeThrough(new TextDecoderStream('utf-8'))
            .getReader()
        while(true) {
            const {value, done} = await reader.read()
            if (done) {
                break
            }
            updateConsumerLog(`${value}`)
        }

    }

    /**
        startConsumer() stops the current consumer process (if any) and starts
                        a new one, listening on the indicated topics
     */
    const reStartConsumer = async (topics: string[]) => {
        updateStatus(`(re-)start consumer started`)
        // TODO do what the contract says above; should cancel any currently running consumer first
        await runConsumer(topics)
        updateStatus(`(re-)start consumer finished`)
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
                            onStartConsumer={reStartConsumer}
                        />
                    </div>
                </div>
            </div>

        </main>
    );
};

export default Consumer;