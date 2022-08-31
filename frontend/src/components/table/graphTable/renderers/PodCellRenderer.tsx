import { CellProps } from "react-table";
import { useState } from "react";
import { GraphTableCell } from "../graphTable";

export const PodCellRenderer = (
  onPodClick: (key: string) => boolean,
  row: CellProps<GraphTableCell>
) => {
  const [toggle, setToggle] = useState<boolean>(true);
  let rendered;
  if (toggle) {
    rendered = (
      row.value.length > 1
        ? [
            <div
              key={-1}
              onClick={() => setToggle(!toggle)}
              className={"arrow"}
            >
              ⬇
            </div>,
          ]
        : []
    ).concat(
      (row.value as string[]).map((value, index) => (
        <div key={index} className="podItem">
          <input
            type="checkbox"
            disabled={
              row.row.original.deployment.startsWith("none")
                ? false
                : !row.row.original.isDeploymentExpanded
            }
            defaultChecked={!row.cell.row.original.podsExpanded[value]}
            onClick={(e) => {
              e.preventDefault();
              const success = onPodClick(value);
              if (success === false) {
                e.stopPropagation();
              }
            }}
          />
          {value}
        </div>
      ))
    );
  } else {
    if (row.value.length > 1) {
      rendered = (
        <div key={0} onClick={() => setToggle(!toggle)} className={"arrow"}>
          ➡
        </div>
      );
    } else {
      rendered = (
        <div key={0} className="podItem">
          <input type="checkbox" defaultChecked={true} />
          {row.value[0]}
        </div>
      );
    }
  }
  return <div>{rendered}</div>;
};
