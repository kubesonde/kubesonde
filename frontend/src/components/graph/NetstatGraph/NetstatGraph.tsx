import { NetstatInfo } from "src/utils/probes";
import cytoscape from "cytoscape";
import { Dict } from "src/entities/types";
import { cytoscapeStylesheet } from "src/components/graph/graphBase/graphBaseOptions";
import CytoscapeComponent from "react-cytoscapejs";

import React, { useEffect, useRef } from "react";
import { SimpleGraphEdge } from "src/entities/graph";
import { NetstatInterface } from "src/utils/netstat";

export interface NetstatGraphProps {
  netstatMappings: NetstatInfo[];
  podIPMappings: Dict;
}

function toCyNode(id: string, label: string): cytoscape.ElementDefinition {
  return Object.assign(
    {},
    {
      data: {
        id: id,
        label: label,
        bg: "#84C3BE",
      },
    }
  );
}

function toCyEdge(
  from: string,
  to: string,
  label: string
): cytoscape.ElementDefinition {
  return Object.assign(
    {},
    {
      data: {
        id: Math.random().toString(),
        source: from,
        target: to,
        label: label,
      },
    }
  );
}

const defineSourceAndDst = (
  podIPMappings: Dict,
  name: string,
  entry: NetstatInterface
) => {
  if (entry.state === "LISTEN") {
    const from = name;
    const to =
      entry.local.address !== undefined
        ? podIPMappings[entry.remote.address] ?? from
        : from;
    const label = (entry.local.address ?? "0.0.0.0") + ":" + entry.local.port;
    const port = entry.local.port.toString();

    return {
      from,
      to,
      label,
      port,
    };
  } else if (entry.state === "SYN_SENT") {
    const from = name;
    const to = podIPMappings[entry.remote.address] ?? entry.remote.address;
    const label = (entry.remote.address ?? "0.0.0.0") + ":" + entry.remote.port;
    const port = entry.remote.port.toString();

    return {
      from,
      to,
      label,
      port,
    };
  } else {
    const remoteAddress = entry.remote.address;
    const from =
      podIPMappings[entry.remote.address] ?? remoteAddress === "127.0.0.1"
        ? name
        : remoteAddress;
    const to = podIPMappings[entry.local.address] ?? name;
    const port = entry.remote.port.toString();
    const label = (entry.remote.address ?? "0.0.0.0") + ":" + entry.remote.port;
    return {
      from,
      to,
      label,
      port,
    };
  }
};

export const NetstatGraph = ({
  netstatMappings,
  podIPMappings,
}: NetstatGraphProps): JSX.Element => {
  const cyRef = useRef<cytoscape.Core>();
  const edges = netstatMappings
    .map((mapping) => {
      const entries = mapping.entries
        .filter((e) => e !== undefined)
        .map((entry, index) => {
          const { from, to, label, port } = defineSourceAndDst(
            podIPMappings,
            mapping.name,
            entry
          );
          return {
            id: "id" + index,
            from,
            to,
            label,
            port,
            deniedConnection: false,
            fromDeployment: undefined,
            toDeployment: undefined,
            hidden: false,
          } as SimpleGraphEdge;
        });
      return entries;
    })
    .flat(1)
    .reduce((acc, curr) => {
      const idx = acc.findIndex(
        (item) => item.to === curr.to && item.from === curr.from
      );
      if (idx === -1) {
        return [...acc, curr];
      } else if (
        acc[idx].label === curr.label ||
        acc[idx].label.includes(curr.label)
      ) {
        return acc;
      } else {
        acc[idx].label = acc[idx].label + ", " + curr.label;
        return acc;
      }
    }, [] as SimpleGraphEdge[])
    .map((e) => toCyEdge(e.from, e.to, e.label));
  const nodes_with_duplicates = edges
    .map((entry) => [entry.data.source as string, entry.data.target as string])
    .flat(1);
  const nodes = Array.from(new Set(nodes_with_duplicates)).map((e) =>
    toCyNode(e, e)
  );
  const layout = {
    name: "grid",
    position() {
      return null;
    },
  };

  useEffect(() => {
    // @ts-ignore
    cyRef.current?.nodes().style("width", (x: any) => {
      return x._private.data.id.length * 6 + "px;";
    });
  });
  const cytoscapeComponentProps = {
    id: "graphId2",
    cy: (cp: cytoscape.Core) => (cyRef.current = cp),
    layout: layout,
    elements: [...nodes, ...edges],
    stylesheet: cytoscapeStylesheet,
  };

  return (
    <div id="graphContainer">
      {/* @ts-ignore */}
      <CytoscapeComponent {...cytoscapeComponentProps} />
    </div>
  );
};
