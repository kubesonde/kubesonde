import React, { useEffect, useMemo, useRef, useState } from "react";
import CytoscapeComponent from "react-cytoscapejs";
import {
  buildColorMap,
  buildInitialEnabledGroups,
  clusterGraph,
  clusterGroups,
  getColoredData,
  getDeployments,
  getGraphDataForTable,
  getPorts,
  mergeEdgesSimple,
  showAllGraph,
  toBoolDict,
  toCyEdge,
  toCyNode,
} from "src/utils/graph";
import {
  cytoscapeStylesheet,
  cytoscapeStylesheetPrintMode,
} from "./graphBaseOptions";
import { GraphNode, SimpleGraphEdge } from "src/entities/graph";
import { GraphTableCell } from "src/components/table/graphTable/graphTable";
import { BoolDict } from "src/entities/types";
import "./graph.css";
import {
  GraphController,
  GraphControllerProps,
} from "src/components/graphControllers/GraphController";
import cytoscape, { NodeSingular } from "cytoscape";
import popper from "cytoscape-popper";
import { BuildPopup } from "src/components/graphPopup/graphPopup";
import { downloadImage, downloadJSON } from "src/utils/downloads";
import { AppearanceController } from "src/components/graphControllers/AppearanceController";
import {
  onEnabledClick,
  onPodClick,
  onPortClick,
} from "src/components/graph/graphBase/GraphBaseAPIs";
import { PodNetworkingInfoV2 } from "src/entities/probeOutput";

cytoscape.use(popper);

export type GraphProps = BasicGraphProps & { title: string };
export interface BasicGraphProps {
  nodes: GraphNode[];
  edges: SimpleGraphEdge[];
  podNetworkInfo?: PodNetworkingInfoV2;
  declarativeConfiguration?: PodNetworkingInfoV2;
}

const headless = process.env.NODE_ENV === "test";

/**
 * Component for showing Kubesonde graph.
 * @type {BasicGraphProps}
 * @example
 * const nodes: GraphNode[] = []
 * const edges: SimpleGraphEdge[] = []
 * const netstat: NetstatInfo[] = []
 * const podNetworkInfo: PodNetworkingInfoV2 = {}
 * return (
 *   <GraphBase nodes={nodes} edges={edges} nestat={netstat} podNetworkInfo={podNetworkingInfo} />
 * )
 */
export const GraphBase: React.FC<GraphProps> = (props: GraphProps) => {
  const [rawData, setRawData] = useState<BasicGraphProps>(props);
  const ports = useMemo(() => getPorts(props.edges), [props]);
  const data = useMemo(
    () => ({
      nodes: rawData.nodes,
      edges: mergeEdgesSimple(rawData.edges),
    }),
    [rawData]
  );

  const probedNetinfo = props.edges.reduce((acc, curr) => {
    const portProto = curr.port.split("/");
    const entry = {
      ip: "none",
      port: portProto[0],
      protocol: portProto[1],
    };
    if (acc[curr.to]) {
      // Remove duplicates
      if (
        acc[curr.to].find((value) => {
          return (
            value.ip === entry.ip &&
            value.port === entry.port &&
            value.protocol === entry.protocol
          );
        })
      ) {
        return acc;
      }

      acc[curr.to] = [...acc[curr.to], entry];
      return acc;
    }
    acc[curr.to] = [entry];
    return acc;
  }, {} as PodNetworkingInfoV2);
  const [filteredPorts, setFilteredPorts] = useState<BoolDict>(
    toBoolDict(ports)
  );
  /* clustered groups keeps track of the expanded deployments*/
  const [expandedDeployments, setExpandedDeployments] = useState<string[]>(
    getDeployments(props.nodes)
  );
  const [expandedPods, setExpandedPods] = useState<BoolDict>(
    toBoolDict(props.nodes.map((node) => node.id))
  );
  /* enabled groups keeps track of the visible deployments*/
  const [activeDeployments, setActiveDeployments] = useState<BoolDict>(
    buildInitialEnabledGroups(props.nodes)
  );
  const [showDenied, setShowDenied] = useState<boolean>(false);
  const [graphData, setGraphData] = useState<cytoscape.ElementDefinition[]>([]);
  const [colorMap, setColorMap] = useState<{ [key: string]: string }>(
    buildColorMap(props.nodes)
  );
  const [printMode, setPrintMode] = useState<boolean>(false);
  const cyRef = useRef<cytoscape.Core>();
  const graphTableData: GraphTableCell[] = useMemo(
    () =>
      getGraphDataForTable(
        data,
        colorMap,
        activeDeployments,
        expandedDeployments,
        expandedPods
      ),
    [data, colorMap, activeDeployments, expandedDeployments, expandedPods]
  );

  const handlePortClick = onPortClick(
    filteredPorts,
    setFilteredPorts,
    rawData,
    setRawData,
    setGraphData,
    data
  );

  const handleEnabledClicked = onEnabledClick(
    rawData,
    setRawData,
    data,
    activeDeployments,
    setActiveDeployments,
    cyRef.current
  );

  /**
   * Show / hide pod from graph. This should work only if the Deployment is ungrouped
   * */
  const handlePodClicked = onPodClick(
    graphData,
    expandedPods,
    setExpandedPods,
    rawData,
    setRawData,
    cyRef.current
  );

  function showDeniedRules() {
    setShowDenied(!showDenied);
    const newData = clusterGroups(data, expandedDeployments);
    if (showDenied) {
      setGraphData([
        ...newData.nodes.map(toCyNode),
        ...newData.edges.map(toCyEdge),
      ]);
    } else {
      const { nodes, edges } = showAllGraph(newData);
      setGraphData([...nodes.map(toCyNode), ...edges.map(toCyEdge)]);
    }
  }

  function handleGraphReset() {
    setExpandedPods(toBoolDict(props.nodes.map((node) => node.id)));
    setFilteredPorts(toBoolDict(getPorts(props.edges)));
    setExpandedDeployments(getDeployments(props.nodes));
    setShowDenied(false);
    setActiveDeployments(buildInitialEnabledGroups(props.nodes));
    const coloredData = getColoredData(props, colorMap);
    setRawData(coloredData);
    const cGraph = clusterGraph({
      ...coloredData,
      edges: mergeEdgesSimple(coloredData.edges),
    });
    setGraphData([
      ...cGraph.nodes.map(toCyNode),
      ...cGraph.edges.map(toCyEdge),
    ]);
  }

  function handleDeploymentClick(group: string) {
    // We cannot process deployments without expandedGroups
    if (group === "none") {
      window.alert("You cannot ungroup pods without deployments");
      return;
    }
    if (activeDeployments[group] === false) {
      return;
    }
    const newGroups = expandedDeployments.includes(group)
      ? expandedDeployments.filter((item) => item !== group)
      : [...expandedDeployments, group];
    setExpandedDeployments(newGroups);
    const cGraph = clusterGroups(data, newGroups);
    setGraphData([
      ...cGraph.nodes.map(toCyNode),
      ...cGraph.edges.map(toCyEdge),
    ]);
  }

  /* This runs only during the first render  */
  useEffect(() => {
    const newColors = buildColorMap(rawData.nodes);
    const coloredData = getColoredData(rawData, newColors);
    setColorMap(newColors);
    setRawData(coloredData);

    const cGraph = clusterGraph({
      ...coloredData,
      edges: mergeEdgesSimple(coloredData.edges),
    });
    setGraphData([
      ...cGraph.nodes.map(toCyNode),
      ...cGraph.edges.map(toCyEdge),
    ]);

    cyRef.current?.on("click", "node", (event) => {
      const node = event.target as NodeSingular;
      const textToDisplay = data.nodes.find(
        (nod) => nod.id === node.id()
      )?.title;
      const reference: DOMRect = node.popperRef()
        .getBoundingClientRect as unknown as DOMRect;
      BuildPopup(reference, textToDisplay || "").show();
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    const maxW = Math.max(...data.nodes.map((node) => node.id.length));
    cyRef.current?.nodes().style("width", () => {
      return maxW * 8 + "px;"; // return x._private.data.id.length * 6 + "px;"; //  x.data('name').length + 'px;'
    });
  });

  const layout = {
    name: "concentric",
    boxSelectionEnabled: false,
    autounselectify: true,
    position() {
      return null;
    },
  };

  const graphControllerProps: GraphControllerProps = {
    ports: {
      netstat: props.podNetworkInfo ?? {},
      declared: props.declarativeConfiguration ?? {},
      probed: probedNetinfo,
    },
    showDeniedConnections: showDenied,
    showDeniedConnectionsHandler: showDeniedRules,
    tableData: graphTableData,
    deploymentClickHandler: handleDeploymentClick,
    enableDeploymentHandler: handleEnabledClicked,
    handleReset: handleGraphReset,
    podClickHandler: handlePodClicked,
    portsFiltered: Object.entries(filteredPorts)
      .filter((value) => value[1] === true)
      .map((value) => value[0]),
    onPortClickHandler: handlePortClick,
  };

  const cytoscapeComponentProps = {
    id: printMode ? "graphIdPrintMode" : "graphId",
    cy: (cp: cytoscape.Core) => (cyRef.current = cp),
    layout: layout,
    elements: [...graphData],
    stylesheet: printMode ? cytoscapeStylesheetPrintMode : cytoscapeStylesheet,
  };

  return (
    <>
      <div className="title"> {props.title}</div>
      <div id="graphContainer">
        {/* @ts-ignore */}
        {headless ? (
          <div />
        ) : (
          <CytoscapeComponent {...cytoscapeComponentProps} />
        )}
        <p />
        <AppearanceController
          onChange={() => setPrintMode(!printMode)}
          printMode={printMode}
          onJSONPrint={async () => {
            await downloadJSON(
              "kubesondeGraph",
              cyRef.current?.json() as unknown as string
            );
          }}
          onImagePrint={async () => {
            await downloadImage(
              "kubesondeGraph",
              cyRef.current?.png({ output: "blob" }) as unknown as Blob
            );
          }}
        />
        <p />
        <GraphController {...graphControllerProps} />
        <p />
      </div>
    </>
  );
};
