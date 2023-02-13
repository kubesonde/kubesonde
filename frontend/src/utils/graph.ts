import { GenericGraph, Graph, GraphEdge, GraphNode, SimpleGraphEdge } from "../entities/graph";
// @ts-ignore
import stronglyConnectedComponents from './graph_manipulator'
import { GraphTableCell } from "../components/table/graphTable/graphTable";
import { BoolDict, Dict } from "../entities/types";
// @ts-ignore
import randomColor from "randomcolor"
import cytoscape, { ElementDefinition } from "cytoscape";

const INTERNET_COLOR = "#FFBF00";
const TEST_POD_COLOR = "red";
//const DEFAULT_COLOR = "#2B7CE9"

export const toClusteredNode = (node: GraphNode): GraphNode => ({
    ...node,
    id: node.deployment ?? node.id,
    label: node.deployment ?? node.label,
    shape: node.deployment ? "square" : "ellipse"
})

export const toClusteredEdge = (edge: GraphEdge): GraphEdge => ({
    ...edge,
    from: edge.fromDeployment ?? edge.from,
    to: edge.toDeployment ?? edge.to
})

export const uniqueNodesReducer = (acc: GraphNode[], curr: GraphNode): GraphNode[] => {
    if (acc.map((item) => item.id).includes(curr.id)) {

        return acc.map((item) => {
            if (item.id === curr.id) {
                let title
                if (item.title === undefined && curr.title !== undefined)
                    title = curr.title
                else if (item.title && !curr.title)
                    title = item.title
                else if (!item.title && !curr.title)
                    title = undefined
                else
                    title = item.title + "\n" + curr.title
                return { ...item, title }
            }
            return item
        })
    } else {
        return [...acc, curr]
    }
}

export const mergeEdges = (edges: GraphEdge[]) => {
    return edges.reduce(
        (acc, curr: GraphEdge) => {
            if (curr.to === curr.from) {
                return acc
            }
            const item = acc.find((item) => item.from === curr.from && item.to === curr.to)
            if (curr.hidden === true) {
                return [...acc, curr]
            }
            if (item) {
                const index = acc.indexOf(item)
                if (acc[index].label === curr.label) {
                    return acc
                }
                acc[index].ports = [...acc[index].ports, ...curr.ports]
                return acc
            }
            return [...acc, curr]
        }, [] as GraphEdge[])
        .map((edge: GraphEdge) => ({ ...edge, ports: Array.from(new Set(edge.ports)), label: Array.from(new Set(edge.ports)).join(", ") }))
}

export function getDeployments(nodes: GraphNode[]): string[] {
    // @ts-ignore
    return Array.from(new Set(nodes.filter((n) => n.deployment !== undefined).map((n) => n.deployment)))
}

export function getPortsFullGraph(edges: GraphEdge[]): string[] {
    return Array.from(new Set(edges.filter((edge) => !edge.deniedConnection).map((edge) => edge.ports).reduce((acc, curr) => ([...acc, ...curr]), [])
    ))
}

export function getPorts(edges: SimpleGraphEdge[]): string[] {
    return Array.from(new Set(edges.filter((edge) => !edge.deniedConnection).map((edge) => edge.port)))
}

export function toBoolDict(values: string[], defaultValue = false): BoolDict {
    const items: BoolDict[] = values.map((value) => ({ [value]: defaultValue }))
    return Object.assign({} as BoolDict, ...items)
}

export function clusterGraph(data: Graph): Graph {
    const nodes = data.nodes
        .map(toClusteredNode)
        .reduce(uniqueNodesReducer, [] as GraphNode[])

    const edges = data.edges
        .map(toClusteredEdge)
    return {
        edges: mergeEdges(edges),
        nodes: nodes as GraphNode[]
    }
}

/**
 * Summary. This function creates clusters only for the selected groups
 *
 **/
export function clusterGroups(data: Graph, groups: string[]): Graph {
    const nodesNotToCluster = data.nodes
        // @ts-ignore
        .filter((node: GraphNode) => !groups.includes(node.group))

    const nodesToCluster = data.nodes
        // @ts-ignore
        .filter((node: GraphNode) => groups.includes(node.group))
        .map((node: GraphNode) => ({
            ...node,
            id: node.deployment ?? node.id,
            label: node.deployment ?? node.label
        } as GraphNode))
        .reduce((acc, curr: GraphNode) => {
            if (acc.map((item) => item.id).includes(curr.id)) {
                return acc
            } else {
                return [...acc, { ...curr, shape: "square" }]
            }
        }, [] as GraphNode[])
    const notToClusterList = nodesNotToCluster.map((node) => node.id)
    const edges = data.edges
        .map((edge: GraphEdge) => ({
            ...edge,
            from: notToClusterList.includes(edge.from) ? edge.from : edge.fromDeployment ?? edge.from,
            to: notToClusterList.includes(edge.to) ? edge.to : edge.toDeployment ?? edge.to,
        })).filter((edge: GraphEdge) => edge.to !== edge.from)


    return { edges: mergeEdges(edges), nodes: [...nodesNotToCluster, ...nodesToCluster] }

}

export function computeMetrics(graph: Graph): { ssc: string[], avgOutDegree: number } {
    const nodesToIds = graph.nodes.filter((node) => !node.hidden).reduce((acc: { [key: string]: number }, curr: GraphNode, index) => {
        acc[curr.id] = index
        return acc
    }, {} as { [key: string]: number })
    const idsToNodes = graph.nodes.filter((node) => !node.hidden).reduce((acc, curr: GraphNode, index) => {
        acc[index] = curr.id
        return acc
    }, {} as string[])
    const adjRaw = graph.edges.reduce((acc: number[][], curr) => {
        if (nodesToIds[curr.from] === undefined || nodesToIds[curr.to] === undefined) {
            return acc
        }
        // @ts-ignore
        acc[nodesToIds[curr.from]] = acc[nodesToIds[curr.from]] ?
            [...acc[nodesToIds[curr.from]], nodesToIds[curr.to]]
            : [nodesToIds[curr.to]]
        return acc
    }, [])
    const adj = Object.values(nodesToIds).map((index) => {
        if (!adjRaw[index]) {
            return []
        }
        return adjRaw[index]
    })
    const connComponents = (stronglyConnectedComponents(adj).components as number[][])
        .map((component) => component.map((item) => idsToNodes[item]).join(", "))

    const outDegrees = adj.map((entry) => entry.length)
    const avgOutDegree = outDegrees.reduce((a, b) => a + b, 0) / outDegrees.length

    return {
        ssc: connComponents,
        avgOutDegree
    }
}

export const getGraphDataForTable =
    (data: Graph,
        colorMap: Dict,
        enabledGroups: BoolDict,
        clusteredGroups: string[],
        expandedPods: BoolDict,
    ): GraphTableCell[] => {
        const nodesThatAreNotServices = data.nodes.filter(node => node.type !== "Service")
        const graphTableData = getDeployments(nodesThatAreNotServices).map((deployment) => {
            const pods = data.nodes.filter((node) => node.deployment === deployment).map((d) => d.name)
            return {
                deployment: deployment,
                background: colorMap[deployment],
                isEnabled: enabledGroups[deployment],
                isDeploymentExpanded: !clusteredGroups.includes(deployment),
                pods,
                podsExpanded: Object.assign({}, ...pods.map((p) => ({ [p]: expandedPods[p] }))),
                podsNumber: pods.length.toString(),
            } as GraphTableCell
        })
        const nodesWithNoDeployment = nodesThatAreNotServices.filter((node) => node.deployment === undefined).map((d) => d.name)
        const graphDataWithNoneDeployments = [...graphTableData,
        ...nodesWithNoDeployment.map((node, idx) => {

            return ({
                deployment: 'none' + idx,
                isEnabled: enabledGroups['none' + idx],
                isDeploymentExpanded: false,
                background: colorMap[node],
                pods: [node],
                podsExpanded: { [node]: expandedPods[node] },
                podsNumber: '1',
            }) as GraphTableCell
        })]
        return graphDataWithNoneDeployments
    }
export function buildColorMap(nodes: GraphNode[]) {
    const deployments = getDeployments(nodes)
    const nodesWithNoDeployment = nodes.filter((node) => node.deployment === undefined).map((d) => d.name)
    const nodeColors = randomColor({
        seed: 2,
        luminosity: 'light',
        count: deployments.length + nodesWithNoDeployment.length
    });
    const deps = [...deployments, ...nodesWithNoDeployment].reduce((prev, curr: string, index) => {
        const newP = {
            ...prev,
        }
        // @ts-ignore
        newP[curr] = nodeColors[index]
        return newP
    }, {} as Dict)
    return {
        ...deps,
        'test-pod': TEST_POD_COLOR,
        'Internet': INTERNET_COLOR
    }
}

export function getColoredData(initialData: GenericGraph<any>, colorMap: { [key: string]: string }): GenericGraph<any> {
    const deploymentColored = {
        ...initialData,
        nodes: initialData.nodes?.map((nod) => (
            {
                ...nod,
                // @ts-ignore
                color: colorMap[nod.deployment] ?? colorMap[nod.name]
            }))
    }

    return {
        ...deploymentColored,
        // @ts-ignore
        nodes: deploymentColored.nodes.map((node) => {
            if (node.id === "test-pod") {
                return { ...node, color: TEST_POD_COLOR }
                // @ts-ignore
            } else if (node.id === 'Internet') {
                return { ...node, color: INTERNET_COLOR }
            }

            return node
        })
    }

}

export function showAllGraph(graph: Graph): Graph {
    const allnodesshown =
        graph.nodes.map((node) => {
            if (node.hidden === true) {
                return { ...node, hidden: false }
            }
            return node
        })
    const alledgesshown =
        graph.edges.map((edge) => {
            if (edge.hidden === true) {
                return {
                    ...edge, hidden: false,
                }
            }
            return edge
        })

    return {
        nodes: allnodesshown,
        edges: alledgesshown
    }
}

export function buildInitialEnabledGroups(nodes: GraphNode[]) {
    const initialEnabledGroups = getDeployments(nodes)
        .reduce((acc: BoolDict, curr: string) => {
            acc[curr] = true
            return acc
        }, {} as BoolDict)
    const initialEnabledGroupsWithNoDeployments = initialEnabledGroups
    const nodesWithNoDeployment = nodes.filter((node) => node.deployment === undefined).map((d) => d.name)
    nodesWithNoDeployment.forEach((_, idx) => initialEnabledGroupsWithNoDeployments['none' + idx] = true)
    return initialEnabledGroups
}

/**
 * @param {SimpleGraphEdge[]} edges   Edges of an unmerged graph
 * @param {string[]}    portsToHide Ports that need to be hidden
 *
 * @return SimpleGraphEdge[] Edges with filtered ports
 * Summary. This function hides the specified ports from the list of edges and
 * shows all the others.
 * */
export function hidePorts(edges: SimpleGraphEdge[], portsToHide: string[]) {

    const newEdges = edges.map((edge) => {
        if (portsToHide.includes(edge.port)) {
            return { ...edge, hidden: true }
        } else if (!edge.deniedConnection) {
            return { ...edge, hidden: false }
        }
        return edge
    })

    return newEdges




}

export function mergeEdgesSimple(simpleEdges: SimpleGraphEdge[]): GraphEdge[] {
    const initialEdges: GraphEdge[] = simpleEdges.reduce(
        (acc, curr: SimpleGraphEdge) => {
            if (curr.hidden === true && !curr.deniedConnection) {
                return acc
            }
            if (curr.to === curr.from) {
                return acc
            }
            const item = acc.find((item) => item.from === curr.from && item.to === curr.to)
            if (curr.hidden === true && curr.deniedConnection) {
                const newItem = { ...curr, ports: [curr.port] }
                // @ts-ignore
                delete newItem['port']
                // @ts-ignore
                delete newItem['timestamp']
                return [...acc, newItem]
            }
            if (item) {
                const index = acc.indexOf(item)
                if (acc[index].ports.includes(curr.port)) {
                    return acc
                }
                acc[index].ports = [...acc[index].ports, curr.port]
                return acc
            }
            const newItem = { ...curr, ports: [curr.port] }
            // @ts-ignore
            delete newItem['port']
            // @ts-ignore
            delete newItem['timestamp']
            return [...acc, newItem]
        }, [] as GraphEdge[])

    return initialEdges.map((edge: GraphEdge) => ({ ...edge, ports: Array.from(new Set(edge.ports)), label: Array.from(new Set(edge.ports)).join(", ") }))
}

export const hidePod = (pod: string) => (node: GraphNode) => {
    if (node.name === pod) {
        return { ...node, hidden: true }
    }
    return node
}

export const showEdgeWithinPod = (pod: string) => (edge: GraphEdge | SimpleGraphEdge) => {
    if ((edge.from === pod || edge.to === pod) && edge.deniedConnection === false) {
        return { ...edge, hidden: false }
    }
    return edge
}

export const showPod = (pod: string) => (node: GraphNode) => {
    if (node.name === pod) {
        return { ...node, hidden: false }
    }
    return node
}

export const getPodIdFromGraphData = (node: ElementDefinition) => node.data.id

/**
 * Internal functions to give pseudo-random colors to the graph
 * */
export function toCyNode(node: GraphNode): cytoscape.ElementDefinition {
    const get_type = (node: GraphNode) => {
        if (node.id === node.deployment) {
            return "deployment"
        } else if (node.type === "Internet") {
            return "internet"
        } else if (node.type === "Service") {
            return "service"
        } else {
            return "pod"
        }
    }
    return Object.assign({}, {
        data: {
            id: node.id,
            label: node.label,
            // @ts-ignore
            bg: node.color,
            type: get_type(node),
            hidden: node.hidden ? "true" : "false"
        }
    })


}

export function toCyEdge(edge: GraphEdge): cytoscape.ElementDefinition {
    return Object.assign({}, {
        data:
        {
            id: Math.random().toString(),
            source: edge.from,
            target: edge.to, label: edge.label,
            hidden: edge.hidden ? "true" : "false",
            denied: edge.deniedConnection ? "true" : "false"
        },
    })
}

export const hideEdgeWithinPod = (pod: string) => (edge: GraphEdge | SimpleGraphEdge) => {
    if (edge.from === pod || edge.to === pod) {
        return { ...edge, hidden: true }
    }
    return edge
}

