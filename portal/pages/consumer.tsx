import React, {useState} from 'react';
import Head from "next/head";
import {logEntry} from "@/components/CommonProps";
import LogWindow from "@/components/LogWindow";
import MessageProducer from "@/components/MessageProducer";
import MessageConsumer from "@/components/MessageConsumer";

const Consumer = (): JSX.Element => {

    const knownTopics = ['stream', 'seer-statement', 'seer-question', 'seer-answer']
    const initialProducerTopic = 'seer-question'
    const initialConsumerTopics = ['seer-answer']

    const [statusText, setStatusText] = useState([] as string[])

    const updateStatus = (statusUpdate: string) => {
        setStatusText(currentStatus => [...currentStatus, logEntry(statusUpdate)]);
    }

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
                            onStatusUpdate={(statusText) => updateStatus(statusText)}
                        />
                    </div>
                    <div className="box-border border-dashed border-2">
                        <MessageConsumer
                            knownTopics={knownTopics}
                            initialConsumerTopics={initialConsumerTopics}
                            textAreaClassName={textAreaClassName}
                            onStatusUpdate={(statusText) => updateStatus(statusText)}
                        />
                    </div>
                </div>
            </div>

        </main>
    );
};

export default Consumer;