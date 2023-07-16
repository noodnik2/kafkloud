import React, {useState} from 'react';
import Button from "@/components/Button";
import SingleSelector from "@/components/SingleSelector";
import {TextAreaProps} from "@/components/CommonProps";
import { TbArrowBigRightLines } from 'react-icons/tb';
import axios from 'axios';

interface MessageProducerProps extends TextAreaProps {
    knownTopics?: string[]
    initiallySelectedTopic: string
    onStatusUpdate: (statusText: string) => void
}

const MessageProducer = ({textAreaClassName = "", knownTopics = [], initiallySelectedTopic, onStatusUpdate}: MessageProducerProps): JSX.Element => {

    const [currentlySelectedTopic, setCurrentlySelectedTopic] = useState(initiallySelectedTopic)
    const [message, setMessage] = useState('')

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
            .then((response) => onStatusUpdate(`delivered(${jsonPayload}) with response(${response.statusText})`))
            .catch((error) => onStatusUpdate(`delivery failure: ${error}`));

    }

    return (
        <div>
            <span>
                <Button name="Clear" onClick={() => setMessage('')} />
                <Button name="Produce" onClick={() => produce(currentlySelectedTopic, message)} />
                <span className="h-1 p-1 m-1">To Topic <TbArrowBigRightLines className="inline"/></span>
                <span className="w-96 w-fit float-right">
                    <SingleSelector
                        labels={knownTopics}
                        selectedOption={currentlySelectedTopic}
                        onChange={
                            topic => {
                                if (topic) {
                                    setCurrentlySelectedTopic(topic)
                                }
                            }
                        }
                    />
                </span>
            </span>
            <div>
                <textarea
                    className={textAreaClassName}
                    placeholder="Enter the value you wish to deliver to the selected topic here"
                    value={message}
                    onChange={e => setMessage(e.target.value)}
                />
            </div>
        </div>
    );
};

export default MessageProducer;