import Select, {MultiValue, StylesConfig} from 'react-select';
import React from "react";

export interface ColourOption {
    readonly value: string;
    readonly label: string;
    readonly color: string;
    readonly isFixed?: boolean;
    readonly isDisabled?: boolean;
}
const colourStyles: StylesConfig<ColourOption, true> = {
    menuPortal: base => ({ ...base, zIndex: 3 }),
    control: (styles) => ({ ...styles, backgroundColor: 'white'}),
    option: (styles, { data, isDisabled, isFocused, isSelected }) => {
        return {
            ...styles,
            backgroundColor: isDisabled
                ? undefined
                : isSelected
                    ? "var(--secoColo)"
                    : isFocused
                        ? 'var(--primColor)'
                        : "white",
            color: isSelected
                    ? "white"
                    : 'var(--secoColor)',
            cursor: isDisabled ? 'not-allowed' : 'default',

            ':active': {
                ...styles[':active'],
                backgroundColor: !isDisabled
                    ? isSelected
                        ? "red"
                        : "var(--secoColor)"
                    : undefined,
                color: !isDisabled
                    ? isSelected
                        ? "var(--primColor)"
                        : "var(--primColor)"
                    : undefined,
            },
        };
    },
    multiValue: (styles, { data }) => {
        return {
            ...styles,
            backgroundColor: data.color
        };
    },
    multiValueLabel: (styles, { data }) => ({
        ...styles,
        color: "white",
    }),
    multiValueRemove: (styles, { data }) => ({
        ...styles,
        color: "white",
        ':hover': {
            backgroundColor: data.color,
            color: 'white',
        },
    }),
};

interface PortSelectorProps {
    data: ColourOption[]
    defaultValue: ColourOption[]
    onChange: (value: MultiValue<ColourOption>) => void
}
export const PortSelector: React.FC<PortSelectorProps>  = ({data,defaultValue,onChange}) => {

    return (<Select
    closeMenuOnSelect={false}
    menuPosition={'fixed'}
    isMulti
    options={data}
    styles={colourStyles}
    value={[...defaultValue]}
    onChange={onChange}
    placeholder={"Select port to filter out..."}
    />)
}
