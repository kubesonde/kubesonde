import {
  ColumnDef,
  useReactTable,
  getCoreRowModel,
  getExpandedRowModel,
  flexRender,
  CellContext,
} from "@tanstack/react-table";
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
    "A port is highlighted in red if does not appear in the declarative configuration but is open in the pod \n" +
    "A port is highlighted in orange if appears in the declarative configuration but is not open in the pod",
};

export const GraphTable = ({
  data,
  ports,
  onDeploymentClick,
  onEnabledClick,
  onPodClick,
}: GraphTableProps) => {
  const columns: ColumnDef<GraphTableCell, any>[] = useMemo(
    () => [
      {
        id: "isEnabled",
        header: "On",
        accessorKey: "isEnabled",
        cell: (info: CellContext<GraphTableCell, boolean>) => {
          if (info.row.original.deployment.startsWith("none")) {
            return <div style={{ textAlign: "center" }}>-</div>;
          }

          return (
            <div style={{ textAlign: "center" }}>
              <input
                type="checkbox"
                defaultChecked={info.getValue()}
                onChange={(e) => {
                  onEnabledClick(
                    info.row.original.deployment,
                    info.row.original.pods,
                    e.target.checked
                  );
                }}
              />
            </div>
          );
        },
      },
      {
        id: "deployment",
        header: "Deployment",
        accessorKey: "deployment",
        cell: (info: CellContext<GraphTableCell, string>) =>
          DeploymentCellRenderer(info, onDeploymentClick),
      },
      {
        id: "podsNumber",
        header: "# Pods",
        accessorKey: "podsNumber",
        cell: (info: CellContext<GraphTableCell, string>) => {
          return (
            <div style={{ textAlign: "center" }}>
              {info.row.original.podsNumber}
            </div>
          );
        },
      },
      {
        id: "pods",
        header: "Pods",
        accessorKey: "pods",
        cell: (info: CellContext<GraphTableCell, string[]>) =>
          PodCellRenderer(onPodClick, info),
      },
      {
        id: "ports",
        header: "Ports exposed",
        cell: (info: CellContext<GraphTableCell, unknown>) =>
          PortCellWithToggleRenderer(ports, info),
      },
    ],
    [onEnabledClick, onDeploymentClick, onPodClick, ports]
  );

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
  });

  // Render Table UI
  return (
    <table role="graphTable">
      <thead>
        {table.getHeaderGroups().map((headerGroup) => (
          <tr key={headerGroup.id}>
            {headerGroup.headers.map((header) => (
              <th
                key={header.id}
                data-tooltip={tooltipsMessages[header.column.id]}
              >
                {header.isPlaceholder
                  ? null
                  : flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
              </th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody>
        {table.getRowModel().rows.map((row) => {
          return (
            <tr
              key={row.id}
              style={{ backgroundColor: row.original.background }}
            >
              {row.getVisibleCells().map((cell) => {
                return (
                  <td key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                );
              })}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};
