import { ProbeErrorInfo } from "../../../utils/probes";
import { Column, useExpanded, useTable } from "react-table";
import React, { useMemo } from "react";

export interface ErrorLogTableProps {
  errorLog: ProbeErrorInfo[];
}
export const ErrorLogTable = ({
  errorLog,
}: ErrorLogTableProps): JSX.Element => {
  const data = errorLog;
  const columns: Column<ProbeErrorInfo>[] = useMemo(
    () => [
      {
        Header: "Name",
        accessor: "podName",
      },
      {
        Header: "Reason",
        accessor: "reason",
      },
      {
        Header: "Timestamp",
        accessor: "timestamp",
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
