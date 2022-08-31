import './Button.css'
import React from "react";
import {IconType} from "react-icons";

export interface ButtonProps {
    title: string
    icon?: React.ReactElement<IconType>
    onClick?: () => void
    type?: "alert"
}

export const Button = ({title, icon,onClick,type}: ButtonProps) => {
    const Icon = icon
    return (
        <div className={`button ${type?? null}`} onClick={onClick}>
            {Icon ? Icon : null}
            <span
                style={{float: "right", paddingLeft: "2px", margin: "0px"}}
            >
                {title}
            </span>

        </div>
    )
}
