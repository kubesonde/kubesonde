import { CellProps } from "react-table";
import { GraphTableCell } from "../graphTable";

export const DeploymentCellRenderer = (
  row: CellProps<GraphTableCell>,
  onDeploymentClick: (key: string) => void
) => {
  if (row.cell.row.original.deployment.startsWith("none")) {
    return (
      <div style={{ textAlign: "center" }} className={"deployment"}>
        -
      </div>
    );
  }

  return (
    <div
      style={{ textAlign: "left" }}
      className={"deployment"}
      onClick={() => {
        onDeploymentClick(row.cell.row.original.deployment);
      }}
    >
      <input
        type="checkbox"
        disabled={row.row.original.isEnabled ? false : true}
        defaultChecked={!row.cell.row.original.isDeploymentExpanded}
      />
      {row.value}
    </div>
  );
};
