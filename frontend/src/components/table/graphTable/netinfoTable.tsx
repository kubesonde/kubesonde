import { Column, useExpanded, useTable } from "react-table";
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
export const NetInfoTable = ({ data }: NetInfoTableProps): JSX.Element => {
  const columns: Column<typeof data[0]>[] = useMemo(
    () => [
      {
        Header: "Name",
        accessor: "podName",
      },
      {
        Header: "Port",
        accessor: "port",
      },
      {
        Header: "Protocol",
        accessor: "protocol",
      },
      {
        Header: "Interface",
        accessor: "ip",
      },
    ],
    []
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
    <table {...getTableProps()}>
      <thead>
        {headerGroups.map((headerGroup) => (
          <tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column) => (
              <th {...column.getHeaderProps()}>{column.render("Header")}</th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody {...getTableBodyProps()}>
        {rows.map((row, i) => {
          prepareRow(row);
          return (
            <tr {...row.getRowProps()}>
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
