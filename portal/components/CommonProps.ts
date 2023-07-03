import {atom} from "recoil";
import memoize from "@emotion/memoize";

export interface TextAreaProps {
    textAreaClassName?: string
}

export const itemWithID = memoize((id: string) => {
    const defaultValue: string[] = []
    return atom({
            key: `item${id}`,
            default: defaultValue
        }
    )
});

export const toLocalLogTime = (d: Date) => {
    const tzoffset = d.getTimezoneOffset() * 60000; //offset in milliseconds
    return (new Date(d.getTime() - tzoffset)).toISOString().slice(0, -1)
}

export interface Option {
    readonly label: string;
    readonly value: string;
}

export const createOption = (label: string) => ({
    label,
    value: label.toLowerCase().replace(/\W/g, ''),
});
