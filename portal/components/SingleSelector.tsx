import React, {Dispatch, useState} from 'react';

import CreatableSelect from 'react-select/creatable';
import {Option, createOption} from "@/components/CommonProps";

interface SingleSelectorProps {
    labels: string[]
    selectedLabel?: string
    onChange: Dispatch<string | undefined | null>
}

const SingleSelector = ({labels, onChange, selectedLabel}: SingleSelectorProps): JSX.Element => {

    const [isLoading, setIsLoading] = useState(false);
    const [options, setOptions] = useState(labels.map(t => createOption(t)));
    const [selectedOption, setSelectedOption] = useState<Option | undefined | null>(selectedLabel? createOption(selectedLabel): null);

    const setAndDispatchValue = (option: Option | null | undefined) => {
        setSelectedOption(option)
        onChange(option?.label)
    }

    const handleCreate = (inputValue: string) => {
        setIsLoading(true);
        const newOption = createOption(inputValue);
        setIsLoading(false);
        setOptions((prev) => [...prev, newOption]);
        setSelectedOption(newOption)
        setAndDispatchValue(newOption)
    };

    return (
        <CreatableSelect
            isClearable
            isDisabled={isLoading}
            isLoading={isLoading}
            onChange={(newValue) => setAndDispatchValue(newValue)}
            onCreateOption={handleCreate}
            options={options}
            value={selectedOption}
        />
    );
};

export default SingleSelector;