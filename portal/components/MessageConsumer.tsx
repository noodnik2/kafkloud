import React, {useState} from 'react';
import Button from "@/components/Button";
import MultiSelector from "@/components/MultiSelector";
import {TextAreaProps} from "@/components/CommonProps";
import {TbArrowBigRightLines} from "react-icons/tb";
import LogWindow from "@/components/LogWindow";

interface MessageConsumerProps extends TextAreaProps {
    knownTopics?: string[]
    currentTopics?: string[]
    loggerText: string[]
    onTopicsChange: (topics: string[]) => void
    onClearConsumer: () => void
}

const MessageConsumer = ({textAreaClassName = "", knownTopics = [], currentTopics = [], loggerText, onTopicsChange, onClearConsumer}: MessageConsumerProps): JSX.Element => {

    return (
        <div>
            <span>
                <Button name="Clear" onClick={onClearConsumer} />
                <span className="h-1 p-1 m-1">Consume From Topics <TbArrowBigRightLines className="inline"/></span>
                <span className="w-96 w-fit float-right">
                    <MultiSelector
                        labels={knownTopics}
                        selectedLabels={currentTopics}
                        onChange={
                            topics => {
                                if (topics) {
                                    onTopicsChange(topics)
                                }
                            }
                        }
                    />
                </span>
            </span>
            <div>
                <LogWindow
                    loggerText={loggerText}
                    textAreaClassName={textAreaClassName}
                    loggerDescription="Text received from the selected topic(s)"
                />
            </div>
        </div>
    );
};

export default MessageConsumer;