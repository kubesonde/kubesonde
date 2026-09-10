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
    // Portal the menu to <body> and lift it above the sticky table header.
    menuPortal: base => ({ ...base, zIndex: 500 }),
    control: (styles, { isFocused }) => ({
        ...styles,
        minHeight: 'var(--control-h)',
        backgroundColor: 'var(--panel)',
        borderColor: isFocused ? 'var(--brand-500)' : 'var(--border)',
        borderRadius: 'var(--radius-sm)',
        boxShadow: isFocused ? '0 0 0 1px var(--brand-500)' : 'var(--shadow-sm)',
        ':hover': { borderColor: 'var(--border-strong)' },
    }),
    placeholder: (styles) => ({ ...styles, color: 'var(--ink-faint)' }),
    input: (styles) => ({ ...styles, color: 'var(--ink)' }),
    menu: (styles) => ({
        ...styles,
        backgroundColor: 'var(--panel)',
        border: '1px solid var(--border)',
        borderRadius: 'var(--radius-md)',
        boxShadow: 'var(--shadow-lg)',
        overflow: 'hidden',
        zIndex: 500,
    }),
    option: (styles, { isDisabled, isFocused, isSelected }) => ({
        ...styles,
        backgroundColor: isSelected
            ? 'var(--brand-500)'
            : isFocused
                ? 'var(--brand-100)'
                : 'transparent',
        color: isSelected ? '#fff' : 'var(--ink)',
        cursor: isDisabled ? 'not-allowed' : 'pointer',
        ':active': {
            backgroundColor: isSelected ? 'var(--brand-600)' : 'var(--brand-100)',
        },
    }),
    multiValue: (styles, { data }) => ({
        ...styles,
        backgroundColor: data.color,
        borderRadius: 'var(--radius-sm)',
    }),
    multiValueLabel: (styles) => ({
        ...styles,
        color: '#fff',
        fontWeight: 600,
    }),
    multiValueRemove: (styles, { data }) => ({
        ...styles,
        color: '#fff',
        ':hover': {
            backgroundColor: data.color,
            color: '#fff',
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
    menuPortalTarget={typeof document !== 'undefined' ? document.body : undefined}
    isMulti
    options={data}
    styles={colourStyles}
    value={[...defaultValue]}
    onChange={onChange}
    placeholder={"Select port to filter out..."}
    />)
}
