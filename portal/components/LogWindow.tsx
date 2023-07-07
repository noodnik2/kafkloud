import {itemWithID, TextAreaProps, toLocalLogTime} from "@/components/CommonProps";
import React from "react";
import {useRecoilState} from "recoil";

interface LogWindowProps extends TextAreaProps {
    loggerDescription?: string
    loggerId: string
}

export default ({textAreaClassName = "", loggerDescription = '', loggerId}: LogWindowProps): JSX.Element => {

    const [logText] = useRecoilState(itemWithID(loggerId))

    return (
        <div>
            <textarea
                className={textAreaClassName}
                readOnly
                placeholder={loggerDescription}
                value={logText.map(s => `${toLocalLogTime(new Date())} ${s}`).join('\n')}
            />
        </div>
    );
}