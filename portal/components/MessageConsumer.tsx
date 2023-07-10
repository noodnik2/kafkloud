import React, {useState} from 'react';
import Button from "@/components/Button";
import MultiSelector from "@/components/MultiSelector";
import {TextAreaProps} from "@/components/CommonProps";
import {TbArrowBigRightLines} from "react-icons/tb";
import LogWindow from "@/components/LogWindow";

interface MessageConsumerProps extends TextAreaProps {
    knownTopics?: string[]
    selectedTopics?: string[]
    outputItemID: string
    onStartConsumer: (topics: string[]) => void
}

const MessageConsumer = ({textAreaClassName = "", knownTopics = [], selectedTopics = [], outputItemID, onStartConsumer}: MessageConsumerProps): JSX.Element => {

    const [currentlySelectedTopics, setCurrentlySelectedTopics] = useState(selectedTopics);

    return (
        <div>
            <span>
                <Button
                    name="Consume"
                    onClick={
                        () => {
                            onStartConsumer(currentlySelectedTopics)
                        }
                    }
                />
                <span className="h-1 p-1 m-1">From Topics <TbArrowBigRightLines className="inline"/></span>
                <span className="w-96 w-fit float-right">
                    <MultiSelector
                        labels={knownTopics}
                        selectedLabels={currentlySelectedTopics}
                        onChange={
                            topics => {
                                if (topics) {
                                    setCurrentlySelectedTopics(topics)
                                }
                            }
                        }
                    />
                </span>
            </span>
            <div>
                <LogWindow
                    loggerId={outputItemID}
                    textAreaClassName={textAreaClassName}
                    loggerDescription="Text received from the selected topic(s)"
                />
            </div>
        </div>
    );
};

export default MessageConsumer;