import { ProbeErrorInfo } from "../../../utils/probes";
import {
  ColumnDef,
  useReactTable,
  getCoreRowModel,
  getExpandedRowModel,
  flexRender,
} from "@tanstack/react-table";
import React, { useMemo } from "react";

export interface ErrorLogTableProps {
  errorLog: ProbeErrorInfo[];
}
export const ErrorLogTable = ({
  errorLog,
}: ErrorLogTableProps): React.JSX.Element => {
  const data = errorLog;
  const columns: ColumnDef<ProbeErrorInfo, any>[] = useMemo(
    () => [
      {
        id: "podName",
        header: "Name",
        accessorKey: "podName",
      },
      {
        id: "reason",
        header: "Reason",
        accessorKey: "reason",
      },
      {
        id: "timestamp",
        header: "Timestamp",
        accessorKey: "timestamp",
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
