import React, {useState, useEffect} from 'react';
import axios from 'axios';
import Head from "next/head";
import {toLocalLogTime} from "@/components/CommonProps";
import LogWindow from "@/components/LogWindow";
import MessageProducer from "@/components/MessageProducer";
import MessageConsumer from "@/components/MessageConsumer";
import {getServiceAddr, SERVICENAME_MONITOR} from "@/routes/info";

const Consumer = (): JSX.Element => {

    const initialConsumerText: string[] = []
    const initialStatusText: string[] = []

    const knownTopics = ['stream', 'seer-statement', 'seer-question', 'seer-answer']
    const initialProducerTopic = 'seer-question'
    const initialConsumerTopics = ['seer-answer']

    const [statusText, setStatusText] = useState(initialStatusText)
    const [consumerText, setConsumerText] = useState(initialConsumerText)
    const [consumerTopics, setConsumerTopics] = useState(initialConsumerTopics);

    const timestamp = (s: string) => {
        return `${toLocalLogTime(new Date())} ${s}`
    }

    const updateStatus = (newLogEntry: string) => {
        setStatusText(currentLog => [...currentLog, timestamp(newLogEntry)]);
    }

    const updateConsumerLog = (newLogEntry: string) => {
        setConsumerText(currentLog => [...currentLog, timestamp(newLogEntry)]);
    }

    const produce = (topic: string, message: string) => {

        const payload = {
            topic: topic,
            message: message,
            key: new Date().toISOString(),
        }

        const headers = {
            headers: {
                'Content-Type': 'application/json'
            }
        }

        const url = `/api/outbox`;
        console.log(`post(${url})`)
        const jsonPayload = JSON.stringify(payload);
        axios.post(url, jsonPayload, headers)
            .then((response) => {
                updateStatus(`delivered(${jsonPayload}) with response(${response.statusText})`)
            })
            .catch((error) => {
                const producerRequestError = `producerRequestError(${error})`;
                console.error(producerRequestError);
                updateStatus(producerRequestError);
            });

    }

    useEffect(() => {

        const serviceAddr = getServiceAddr(SERVICENAME_MONITOR);
        const url = `${serviceAddr}/consume?topics=${consumerTopics}`;

        const es = new EventSource(url);
        es.onmessage = console.log
        es.onerror = console.log
        es.onopen = console.log
        es.addEventListener('open', () => updateStatus(`opened ${url}`));
        es.addEventListener('streamer', (e) => updateConsumerLog(e.data));
        es.addEventListener('error', () => updateStatus(`error from ${url}`));
        updateStatus(`requesting ${url}`);
        return () => {
            updateStatus(`closing ${url}`);
            es.close();
        }

    }, [consumerTopics]);

    const title = (text: string) => {
        return (
            <h1 className="p-2 text-2xl decoration-2 font-bold">{text}</h1>
        );
    }

    const textAreaClassName = "p-2 overflow-y-auto w-full h-24 min-h-full"

    return (
        <main className="flex items-center justify-around font-sans">
            <Head>
                <title>Consumer</title>
            </Head>

            <div className="p-2 m-5 w-full shadow-2xl drop-shadow-lg space-y-3">
                <div className="border-2 border-black bg-blue-50">
                    {title("Status Updates")}
                    <LogWindow
                        loggerText={statusText}
                        loggerDescription="Consumer status updates"
                        textAreaClassName={textAreaClassName}
                    />
                </div>
                <div className="border-2 border-black bg-blue-50">
                    {title("Messages")}
                    <div className="box-border border-dashed border-2">
                        <MessageProducer
                            knownTopics={knownTopics}
                            initiallySelectedTopic={initialProducerTopic}
                            textAreaClassName={textAreaClassName}
                            onDeliver={produce}
                        />
                    </div>
                    <div className="box-border border-dashed border-2">
                        <MessageConsumer
                            knownTopics={knownTopics}
                            currentTopics={consumerTopics}
                            textAreaClassName={textAreaClassName}
                            loggerText={consumerText}
                            onTopicsChange={setConsumerTopics}
                            onClearConsumer={() => setConsumerText([])}
                        />
                    </div>
                </div>
            </div>

        </main>
    );
};

export default Consumer;