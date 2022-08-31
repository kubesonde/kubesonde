import React from "react";
import { GraphTable, GraphTableCell } from "../table/graphTable/graphTable";
import "./graphController.css";
import Switch from "react-switch";
import { ColourOption, PortSelector } from "./PortSelector";
import { MultiValue } from "react-select";
import { Button } from "src/components/button/Button";
import { PodNetworkingInfoV2 } from "src/entities/probeOutput";

export interface GraphControllerProps {
  ports: {
    netstat: PodNetworkingInfoV2;
    declared: PodNetworkingInfoV2;
    probed: PodNetworkingInfoV2;
  };
  showDeniedConnections: boolean;
  showDeniedConnectionsHandler: () => void;
  tableData: GraphTableCell[];
  deploymentClickHandler: (key: string) => void;
  enableDeploymentHandler: (
    group: string,
    pods: string[],
    enable: boolean
  ) => void;
  handleReset: () => void;
  podClickHandler: (pod: string) => boolean;
  portsFiltered: string[];
  onPortClickHandler: (currentPorts: string[]) => void;
}

/**
 * Component for controlling Kubesonde graph.
 *
 * @component
 * @example
 * const props: GraphControllerProps
 * return (
 *   <GraphController {...{props}} />
 * )
 */
export const GraphController: React.FC<GraphControllerProps> = ({
  ports,
  showDeniedConnections,
  showDeniedConnectionsHandler,
  tableData,
  deploymentClickHandler,
  enableDeploymentHandler,
  handleReset,
  podClickHandler,
  portsFiltered,
  onPortClickHandler,
}: GraphControllerProps) => {
  const filterablePorts = Object.values(ports.probed)
    .flat()
    .filter((p) => p.protocol !== undefined);

  const data = filterablePorts.map((port) => ({
    value: port.port + "/" + port.protocol,
    label: port.port + "/" + port.protocol,
    color: "var(--secoColor)",
  }));
  const onChange = (values: MultiValue<ColourOption>) => {
    const ids = values.map((value) => value.value);
    return onPortClickHandler(ids);
  };
  const portsFilt = portsFiltered.map((port) => ({
    value: port,
    label: port,
    color: "var(--secoColor)",
  }));
  return (
    <>
      <div id={"ControlsWrapper"}>
        <div className="col">
          <PortSelector
            data={data}
            defaultValue={portsFilt}
            onChange={onChange}
          />
        </div>
        <div className="col">
          <div id="SwitchWrapper">
            <Switch
              role="switch"
              height={14}
              width={30}
              checkedIcon={false}
              uncheckedIcon={false}
              onColor="#219de9"
              offColor="#bbbbbb"
              checked={showDeniedConnections}
              onChange={showDeniedConnectionsHandler}
            />
            <span> Show denied connections</span>

            <span style={{ margin: "4px" }} />
            <Button
              title={"Reset graph"}
              onClick={handleReset}
              type={"alert"}
            />
          </div>
        </div>
      </div>
      <p></p>
      <div className="viewWrapper">
        <div id="tableWrapper" className="w100">
          <GraphTable
            ports={ports}
            data={tableData}
            onDeploymentClick={deploymentClickHandler}
            onEnabledClick={enableDeploymentHandler}
            onPodClick={podClickHandler}
          />
        </div>
      </div>
    </>
  );
};
