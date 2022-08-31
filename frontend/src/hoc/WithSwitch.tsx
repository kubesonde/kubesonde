import Switch from "react-switch";
import React, {ComponentProps, ComponentType, useState} from "react";
import './WithSwitch.css'
export interface WithSwitchProps {
    title: string
}

export function WithSwitch<T> (Component: ComponentType<T>,title: string) {
    const [checked, setChecked] = useState<boolean>(false)
    return (hocProps: ComponentProps<typeof Component>) => {
        return (
<>
            <div className="withSwitch">
                <Switch
                    role="switch"
                    height={14}
                    width={30}
                    checkedIcon={false}
                    uncheckedIcon={false}
                    onColor="#219de9"
                    offColor="#bbbbbb"
                    checked={checked}
                    onChange={() => setChecked(!checked)}/>
                <span className="switchText">{ (checked ? "Hide " : "Show ") + title}</span>
            </div>
                <div></div>
                {checked && <Component {...(hocProps as T)}/>}
</>
        )

    }

}
