import { CellContext } from "@tanstack/react-table";
//import { useState } from "react";
import { GraphTableCell } from "../graphTable";

export const PodCellRenderer = (
  onPodClick: (key: string) => boolean,
  row: CellContext<GraphTableCell, string[]>
) => {
  //const [toggle, setToggle] = useState<boolean>(true);
  const toggle = true;
  const value = row.getValue();
  let rendered;
  if (toggle) {
    rendered = value.map((value, index) => (
      <div key={"podItem" + index} className="podItem">
        <input
          key={"podItem_checkbox" + index}
          type="checkbox"
          disabled={
            row.row.original.deployment.startsWith("none")
              ? false
              : !row.row.original.isDeploymentExpanded
          }
          defaultChecked={!row.row.original.podsExpanded[value]}
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
    ));
  } else {
    if (value.length > 1) {
      rendered = (
        <div
          key={0}
          // onClick={() => setToggle(!toggle)}
          className={"arrow"}
        >
          ➡
        </div>
      );
    } else {
      rendered = (
        <div key={0} className="podItem">
          <input type="checkbox" defaultChecked={true} />
          {value[0]}
        </div>
      );
    }
  }
  return <div className="podGroup">{rendered}</div>;
};
