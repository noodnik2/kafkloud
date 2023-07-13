import React, {useEffect, useState} from 'react';
import Button from "@/components/Button";
import MultiSelector from "@/components/MultiSelector";
import {TextAreaProps, logEntry} from "@/components/CommonProps";
import {TbArrowBigRightLines} from "react-icons/tb";
import LogWindow from "@/components/LogWindow";
import {getServiceAddr, SERVICENAME_MONITOR} from "@/routes/info";

interface MessageConsumerProps extends TextAreaProps {
    knownTopics?: string[]
    initialConsumerTopics?: string[]
    onStatusUpdate: (statusText: string) => void
}

const MessageConsumer = ({textAreaClassName = "", knownTopics = [], initialConsumerTopics = [], onStatusUpdate}: MessageConsumerProps): JSX.Element => {

    const [consumerText, setConsumerText] = useState([] as string[])
    const [consumerTopics, setConsumerTopics] = useState(initialConsumerTopics);

    const updateConsumerLog = (newLogEntry: string) => {
        setConsumerText(currentLog => [...currentLog, logEntry(newLogEntry)]);
    }

    useEffect(() => {

        const serviceAddr = getServiceAddr(SERVICENAME_MONITOR);
        const url = `${serviceAddr}/consume?topics=${consumerTopics}`;
        onStatusUpdate(`consumer URL: ${url}`);

        const es = new EventSource(url);
        es.onmessage = console.log
        es.onerror = console.log
        es.onopen = console.log
        es.addEventListener('open', () => onStatusUpdate(`consumer opened`));
        es.addEventListener('error', () => onStatusUpdate(`consumer error`));
        es.addEventListener('streamer', (e) => updateConsumerLog(e.data));
        return () => {
            onStatusUpdate(`closing consumer`);
            es.close();
        }

    }, [consumerTopics]);

    return (
        <div>
            <span>
                <Button name="Clear" onClick={() => setConsumerText([])} />
                <span className="h-1 p-1 m-1">Consume From Topics <TbArrowBigRightLines className="inline"/></span>
                <span className="w-96 w-fit float-right">
                    <MultiSelector
                        labels={knownTopics}
                        selectedLabels={consumerTopics}
                        onChange={
                            topics => {
                                if (topics) {
                                    setConsumerTopics(topics)
                                }
                            }
                        }
                    />
                </span>
            </span>
            <div>
                <LogWindow
                    loggerText={consumerText}
                    textAreaClassName={textAreaClassName}
                    loggerDescription="Text received from the selected topic(s)"
                />
            </div>
        </div>
    );
};

export default MessageConsumer;