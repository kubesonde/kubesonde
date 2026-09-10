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
        // Highlight the focused (keyboard/hover) option with a solid brand fill
        // so it clearly stands out; selected sits one shade darker.
        backgroundColor: isSelected
            ? 'var(--brand-600)'
            : isFocused
                ? 'var(--brand-500)'
                : 'transparent',
        color: isSelected || isFocused ? '#fff' : 'var(--ink)',
        cursor: isDisabled ? 'not-allowed' : 'pointer',
        ':active': {
            backgroundColor: 'var(--brand-700)',
            color: '#fff',
        },
    }),
    // Selected port chips. Port colors are light pastels (they mirror the graph
    // palette), so a consistent brand-tinted chip stays readable rather than
    // white-on-pastel.
    multiValue: (styles) => ({
        ...styles,
        backgroundColor: 'var(--brand-100)',
        borderRadius: 'var(--radius-sm)',
    }),
    multiValueLabel: (styles) => ({
        ...styles,
        color: 'var(--brand-700)',
        fontWeight: 600,
    }),
    multiValueRemove: (styles) => ({
        ...styles,
        color: 'var(--brand-600)',
        borderRadius: '0 var(--radius-sm) var(--radius-sm) 0',
        ':hover': {
            backgroundColor: 'var(--brand-200)',
            color: 'var(--brand-700)',
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
