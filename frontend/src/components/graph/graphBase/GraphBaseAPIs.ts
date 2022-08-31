import {BoolDict} from "src/entities/types";
import {BasicGraphProps} from "src/components/graph/graphBase/GraphBase";
import {
    clusterGraph,
    getPodIdFromGraphData,
    hideEdgeWithinPod,
    hidePod,
    hidePorts, mergeEdgesSimple,
    showEdgeWithinPod,
    showPod, toCyEdge, toCyNode
} from "src/utils/graph";
import cytoscape from "cytoscape";
import {GraphEdge, GraphNode} from "src/entities/graph";

export const onPodClick = (
    graphData: cytoscape.ElementDefinition[],
    expandedPods: BoolDict,
    setExpandedPods: (props: BoolDict) => void,
    rawData: BasicGraphProps,
    setRawData: (props: BasicGraphProps) => void,
    cy: cytoscape.Core | undefined
) => (pod: string) => {
    const isClustered = !graphData.map(getPodIdFromGraphData).includes(pod)
    if (isClustered) {
        return false
    }
    const currExpandedPods = expandedPods
    const nodes = rawData.nodes.map(currExpandedPods[pod] ? showPod(pod) : hidePod(pod))
    const edges = rawData.edges.map(currExpandedPods[pod] ? showEdgeWithinPod(pod) : hideEdgeWithinPod(pod))
    // @ts-ignore
    setRawData({nodes, edges})
    cy?.nodes(`node[id= "${pod}"]`).style("visibility", currExpandedPods[pod] ? "visible" : "hidden")
    cy?.nodes(`node[id= "${pod}"]`).connectedEdges().style("visibility", currExpandedPods[pod] ? "visible" : "hidden")


    currExpandedPods[pod] = !currExpandedPods[pod]
    setExpandedPods(currExpandedPods)
    return true
}

export const onPortClick = (
    filteredPorts: BoolDict,
    setFilteredPorts: (props: BoolDict) => void,
    rawData: BasicGraphProps,
    setRawData: (props: BasicGraphProps) => void,
    setGraphData: (props: cytoscape.ElementDefinition[]) => void,
    data:  {nodes: GraphNode[], edges: GraphEdge[]}
) => (clickedPorts: string[]) => {
    const oldPorts = filteredPorts
    Object.entries(filteredPorts)
        .forEach(((value) => {
            if (clickedPorts.includes(value[0])) {
                oldPorts[value[0]] = true
            } else {
                oldPorts[value[0]] = false
            }
        }))
    setFilteredPorts(oldPorts)
    const newEdges = hidePorts(rawData.edges, clickedPorts)
    setRawData({
        nodes: rawData.nodes,
        edges: newEdges
    })
    const currData = {
        nodes: data.nodes,
        edges: mergeEdgesSimple(newEdges)
    }
    const cGraph = clusterGraph(currData)
    setGraphData([...cGraph.nodes.map(toCyNode), ...cGraph.edges.map(toCyEdge)])
}

export const onEnabledClick = (
    rawData: BasicGraphProps,
    setRawData: (props: BasicGraphProps) => void,
    data:  {nodes: GraphNode[], edges: GraphEdge[]},
    activeDeployments: BoolDict,
    setActiveDeployments: (props: BoolDict) => void,
    cy: cytoscape.Core | undefined,
) => (group: string, pods: string[], enable: boolean) => {
    let newNodes;
    if (group.startsWith('none')) {
        newNodes = handleEnabledWhenNoDeployment(data,pods, enable)
    } else {
        newNodes = handleEnabledWhenDeployment(data,group, enable)
    }
    const prevEnabled = activeDeployments
    prevEnabled[group] = !prevEnabled[group]
    setActiveDeployments(prevEnabled)
    const newData = {...rawData, nodes: newNodes}
    setRawData(newData)
    cy?.nodes(`node[id= "${group}"]`).style("visibility", !prevEnabled[group] ? "hidden" : "visible")
    cy?.nodes(`node[id= "${group}"]`).connectedEdges().style("visibility", !prevEnabled[group] ? "hidden" : "visible")
    return
}

const handleEnabledWhenNoDeployment = (data:  {nodes: GraphNode[], edges: GraphEdge[]},pods: string[], enable: boolean) => {
    return data.nodes.map((node) => {
        if (pods.includes(node.id)) {
            return {...node, hidden: !enable}
        }
        return node
    })
}
const handleEnabledWhenDeployment = (data:  {nodes: GraphNode[], edges: GraphEdge[]},group: string, enable: boolean) => {
    return data.nodes.map((node) => {
        if (node.group === group) {
            return {...node, hidden: !enable}
        }
        return node
    })
}
