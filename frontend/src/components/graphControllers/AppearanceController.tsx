import Switch from "react-switch";
import { AiOutlineDownload } from "react-icons/ai";
import React from "react";
import "./AppearanceController.scss";
import { Button } from "src/components/button/Button";
export interface AppearanceControllerProps {
  onChange: () => void;
  printMode: boolean;
  onJSONPrint: () => Promise<void>;
  onImagePrint: () => Promise<void>;
}

/**
 * Component for controlling Kubesonde graph appearance.
 *
 * @component
 * @example
 * const props: AppearanceControllerProps
 * return (
 *   <AppearanceController {{onChange, printMode, onJSONPrint, onImagePrint}: AppearanceControllerProps } />
 * )
 */
export const AppearanceController = ({
  onChange,
  printMode,
  onJSONPrint,
  onImagePrint,
}: AppearanceControllerProps): JSX.Element => {
  return (
    <>
      <div id={"AppearanceControllerContainer"}>
        <div id={"SwitchWrapper"}>
          <Switch
            role="switch"
            height={14}
            width={30}
            checkedIcon={false}
            uncheckedIcon={false}
            onColor="#219de9"
            offColor="#bbbbbb"
            checked={printMode}
            onChange={onChange}
          />
          <span> Print mode</span>
        </div>
        <span
          style={{
            padding: "4px",
          }}
        />
        <Button
          icon={<AiOutlineDownload size={18} />}
          title={"Download graph as JSON"}
          onClick={onJSONPrint}
        />
        <span
          style={{
            padding: "4px",
          }}
        />
        <Button
          icon={<AiOutlineDownload size={18} />}
          title={"Download graph as PNG"}
          onClick={onImagePrint}
        />
      </div>
    </>
  );
};
