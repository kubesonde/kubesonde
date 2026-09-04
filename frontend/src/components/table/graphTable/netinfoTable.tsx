import {
  ColumnDef,
  useReactTable,
  getCoreRowModel,
  getExpandedRowModel,
  flexRender,
} from "@tanstack/react-table";
import React, { useMemo } from "react";

export interface NetInfoTableEntry {
  podName: string;
  port: string;
  protocol: string;
  ip: string;
}
export interface NetInfoTableProps {
  data: NetInfoTableEntry[];
}
export const NetInfoTable = ({ data }: NetInfoTableProps): React.JSX.Element => {
  const columns: ColumnDef<NetInfoTableEntry, any>[] = useMemo(
    () => [
      {
        id: "podName",
        header: "Name",
        accessorKey: "podName",
      },
      {
        id: "port",
        header: "Port",
        accessorKey: "port",
      },
      {
        id: "protocol",
        header: "Protocol",
        accessorKey: "protocol",
      },
      {
        id: "ip",
        header: "Interface",
        accessorKey: "ip",
      },
    ],
    []
  );

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
  });

  // Render Table UI
  return (
    <table>
      <thead>
        {table.getHeaderGroups().map((headerGroup) => (
          <tr key={headerGroup.id}>
            {headerGroup.headers.map((header) => (
              <th key={header.id}>
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
            <tr key={row.id}>
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
