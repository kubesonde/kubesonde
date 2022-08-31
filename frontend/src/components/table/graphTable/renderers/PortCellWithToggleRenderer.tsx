import { CellProps } from "react-table";
import { GraphTableCell } from "../graphTable";
import { PodNetworkingInfoV2 } from "src/entities/probeOutput";

export const PortCellWithToggleRenderer = (
  ports: {
    netstat: PodNetworkingInfoV2;
    declared: PodNetworkingInfoV2;
    probed: PodNetworkingInfoV2;
  },
  row: CellProps<GraphTableCell>
) => {
  //const [toggle, setToggle] = useState<boolean>(true);
  const pods = row.row.original.pods;
  const portProtoNetstat = pods
    .filter((pod) => ports.netstat[pod])
    .map((pod) =>
      ports.netstat[pod].map((mapping) => mapping.port + "/" + mapping.protocol)
    )
    .flat();
  const portProtoDeclared = pods
    .filter((pod) => ports.declared[pod])
    .map((pod) =>
      ports.declared[pod].map(
        (mapping) => mapping.port + "/" + mapping.protocol
      )
    )
    .flat();
  const portProtoProbed = pods
    .filter((pod) => ports.probed[pod])
    .map((pod) =>
      ports.probed[pod].map((mapping) => mapping.port + "/" + mapping.protocol)
    )
    .flat();
  const okPorts = portProtoProbed.filter(
    (portProto) =>
      portProtoDeclared.includes(portProto) &&
      portProtoNetstat.includes(portProto)
  );
  const okPortsRendered = okPorts.map((item, index) => (
    <div key={"ns" + index} className="okPort">
      {item}
    </div>
  ));
  const errorPorts = portProtoProbed.filter(
    (portProto) =>
      !portProtoDeclared.includes(portProto) &&
      portProtoNetstat.includes(portProto)
  );
  const errorPortsRendered = errorPorts.map((item, index) => (
    <div key={"ns" + index + okPorts.length} className="errorPort">
      {item}
    </div>
  ));
  const warningPorts = portProtoProbed.filter(
    (portProto) =>
      portProtoDeclared.includes(portProto) &&
      !portProtoNetstat.includes(portProto)
  );
  const warningPortsRendered = warningPorts.map((item, index) => (
    <div
      key={"ns" + index + okPorts.length + errorPorts.length}
      className="warningPort"
    >
      {item}
    </div>
  ));
  return okPortsRendered
    .concat(errorPortsRendered)
    .concat(warningPortsRendered);
  /*const netstatPorts = row.row.original.netstatPorts;
  const itemsInNetstatAndNotInProbe = netstatPorts
    ?.filter((port) => !row.value.includes(port))
    .map((item, index) => (
      <div key={"ns" + index} className="portItemNetstatOnly">
        {item}
      </div>
    ));
  let rendered;
  if (isOpen(toggle)) {
    if (moreThanOnePort(row.value)) {
      rendered = [
        <div key={-1} onClick={() => setToggle(!toggle)} className={"arrow"}>
          ⬇
        </div>,
      ]
        .concat(
          (row.value as string[]).map((value, index) => {
            return (
              <div
                key={index}
                className={classNameForPort(netstatPorts, value)}
              >
                {value}
              </div>
            );
          })
        )
        .concat(itemsInNetstatAndNotInProbe ?? []);
    } else {
      rendered = (
        <div key={0} className={classNameForPort(netstatPorts, row.value[0])}>
          {row.value[0]}
        </div>
      );
    }
  } else {
    if (
      moreThanOnePort(row.value) ||
      onePortAndMultipleNetstat(
        row.value,
        itemsInNetstatAndNotInProbe as unknown as string[] | undefined
      )
    ) {
      rendered = (
        <div key={0} onClick={() => setToggle(!toggle)} className={"arrow"}>
          ➡
        </div>
      );
    } else {
      rendered = (
        <div key={0} className="portItem">
          {row.value[0]}
        </div>
      );
    }
  }
  return <div>{rendered}</div>;*/
};
