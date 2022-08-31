import { GraphEdge, GraphNode, SimpleGraphEdge } from "./graph";
import { Builder, IBuilder } from "builder-pattern";

const defaultEdge: GraphEdge = {
    deniedConnection: false,
    from: "A",
    fromDeployment: undefined,
    id: "1",
    label: "test",
    ports: [],
    to: "B",
    toDeployment: undefined,
    hidden: false,
}

const defaultSimpleEdge: SimpleGraphEdge = {
    deniedConnection: false,
    from: "A",
    fromDeployment: undefined,
    id: "1",
    label: "test",
    port: '80',
    to: "B",
    toDeployment: undefined,
    hidden: false,
    timestamp: 0,
}

const defaultNode: GraphNode = {
    type: "Pod",
    deployment: undefined,
    group: undefined,
    id: "1",
    label: "A",
    name: "A"

}

export function GraphEdgeBuilder(): IBuilder<GraphEdge> {
    return Builder<GraphEdge>(defaultEdge)
}

export function SimpleGraphEdgeBuilder(): IBuilder<SimpleGraphEdge> {
    return Builder<SimpleGraphEdge>(defaultSimpleEdge)
}

export function GraphNodeBuilder(): IBuilder<GraphNode> {
    return Builder<GraphNode>(defaultNode)
}
