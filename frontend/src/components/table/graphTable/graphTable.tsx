import { Column, useTable, useExpanded, CellProps } from "react-table";
import { useMemo } from "react";
import "./table.css";
import { PodCellRenderer } from "./renderers/PodCellRenderer";
import { PortCellWithToggleRenderer } from "./renderers/PortCellWithToggleRenderer";
import { DeploymentCellRenderer } from "./renderers/DeploymentCellRenderer";
import { BoolDict, Dict } from "../../../entities/types";
import { PodNetworkingInfoV2 } from "src/entities/probeOutput";

export interface GraphTableCell {
  deployment: string;
  isEnabled: boolean;
  isDeploymentExpanded: boolean;
  background: string;
  pods: string[];
  podsExpanded: BoolDict;
  podsNumber: string;
}

export interface GraphTableProps {
  data: GraphTableCell[];
  ports: {
    netstat: PodNetworkingInfoV2;
    declared: PodNetworkingInfoV2;
    probed: PodNetworkingInfoV2;
  };
  onDeploymentClick: (key: string) => void;
  onEnabledClick: (key: string, pods: string[], enable: boolean) => void;
  onPodClick: (pod: string) => boolean;
}

const tooltipsMessages: Dict = {
  isEnabled: "Hide/Show deployment from graph",
  deployment: "Deployment name. Click to ungroup and reveal its pods",
  podsNumber: "Number of pods in a deployment",
  pods: "List of pods in a deployment",
  ports:
    "List of ports open on a given pod.\n" +
    "A port is highlighted in green if appears in the declarative configuration and is open in the pod \n" +
    "A port is highlighted in red if does not appear in the declarative configuration but is open in the pod",
};

export const GraphTable = ({
  data,
  ports,
  onDeploymentClick,
  onEnabledClick,
  onPodClick,
}: GraphTableProps) => {
  const columns: Column<GraphTableCell>[] = useMemo(
    () => [
      {
        Header: "On",
        accessor: "isEnabled",
        Cell: (row) => {
          if (row.row.original.deployment.startsWith("none")) {
            return <div style={{ textAlign: "center" }}>-</div>;
          }

          return (
            <div style={{ textAlign: "center" }}>
              <input
                type="checkbox"
                defaultChecked={row.value}
                onChange={(e) => {
                  onEnabledClick(
                    row.row.original.deployment,
                    row.row.original.pods,
                    e.target.checked
                  );
                }}
              />
            </div>
          );
        },
      },
      {
        Header: "Deployment",
        accessor: "deployment",
        Cell: (row: CellProps<GraphTableCell>) =>
          DeploymentCellRenderer(row, onDeploymentClick),
      },
      {
        Header: "# Pods",
        accessor: "podsNumber",
        Cell: (row) => {
          return (
            <div style={{ textAlign: "center" }}>
              {row.row.original.podsNumber}
            </div>
          );
        },
      },
      {
        Header: "Pods",
        accessor: "pods",
        Cell: (row: CellProps<GraphTableCell>) =>
          PodCellRenderer(onPodClick, row),
      },
      {
        Header: "Ports exposed",
        Cell: (row: CellProps<GraphTableCell>) =>
          PortCellWithToggleRenderer(ports, row),
      },
    ],
    [onEnabledClick, onDeploymentClick, onPodClick, ports]
  );

  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } =
    useTable(
      {
        columns,
        data,
      },
      useExpanded
    );

  // Render Table UI
  return (
    <table {...{ ...getTableProps(), role: "graphTable" }}>
      <thead>
        {headerGroups.map((headerGroup) => (
          <tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column) => (
              <th
                {...{
                  ...column.getHeaderProps(),
                  "data-tooltip": tooltipsMessages[column.id],
                }}
              >
                {column.render("Header")}
              </th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody {...getTableBodyProps()}>
        {rows.map((row, i) => {
          prepareRow(row);
          const props = {
            ...row.getRowProps(),
            style: {
              ...row.getRowProps().style,
              backgroundColor: row.original.background,
            },
          };
          return (
            <tr {...props}>
              {row.cells.map((cell) => {
                return <td {...cell.getCellProps()}>{cell.render("Cell")}</td>;
              })}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};
