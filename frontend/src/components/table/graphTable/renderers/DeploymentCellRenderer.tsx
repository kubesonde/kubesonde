import { CellContext } from "@tanstack/react-table";
import { GraphTableCell } from "../graphTable";

export const DeploymentCellRenderer = (
  row: CellContext<GraphTableCell, string>,
  onDeploymentClick: (key: string) => void
) => {
  if (row.row.original.deployment.startsWith("none")) {
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
        onDeploymentClick(row.row.original.deployment);
      }}
    >
      <input
        type="checkbox"
        disabled={row.row.original.isEnabled ? false : true}
        defaultChecked={!row.row.original.isDeploymentExpanded}
      />
      {row.getValue()}
    </div>
  );
};
