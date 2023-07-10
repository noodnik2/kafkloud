import {itemWithID, TextAreaProps, toLocalLogTime} from "@/components/CommonProps";
import React from "react";
import {useRecoilState} from "recoil";

interface LogWindowProps extends TextAreaProps {
    loggerDescription?: string
    loggerId: string
}

const LogWindow = ({textAreaClassName = "", loggerDescription = '', loggerId}: LogWindowProps): JSX.Element => {

    const [logText] = useRecoilState(itemWithID(loggerId))

    return (
        <div>
            <textarea
                className={textAreaClassName}
                readOnly
                placeholder={loggerDescription}
                value={logText.join('\n')}
            />
        </div>
    );
};

export default LogWindow;