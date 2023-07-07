import React, {useState} from 'react';
import Button from "@/components/Button";
import SingleSelector from "@/components/SingleSelector";
import {TextAreaProps} from "@/components/CommonProps";
import { TbArrowBigRightLines } from 'react-icons/tb';

interface MessageProducerProps extends TextAreaProps {
    knownTopics?: string[]
    selectedTopic: string
    onDeliver: (topic: string, message: string) => void
}

export default ({textAreaClassName = "", knownTopics = [], selectedTopic, onDeliver}: MessageProducerProps): JSX.Element => {
    const [currentlySelectedTopic, setCurrentlySelectedTopic] = useState(selectedTopic)
    const [message, setMessage] = useState('')
    return (
        <div>
            <span>
                <Button
                    name="Produce"
                    onClick={
                        () => {
                            onDeliver(currentlySelectedTopic, message)
                        }
                    }
                />
                <span className="h-1 p-1 m-1">To Topic <TbArrowBigRightLines className="inline"/></span>
                <span className="w-96 w-fit float-right">
                    <SingleSelector
                        labels={knownTopics}
                        selectedLabel='stream'
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
                    onChange={e => setMessage(e.target.value)}
                />
            </div>
        </div>
    );
}
