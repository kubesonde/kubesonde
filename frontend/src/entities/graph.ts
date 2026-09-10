export interface GraphNode {
    id: string,
    type?: string, //"Pod" | "Service" | "Deployment"
    label: string,
    name: string,
    group?: string,
    deployment: string | undefined,
    title?: string
    shape?: string
    hidden?: boolean
}

export interface GraphEdge {
    from: string,
    to: string,
    id: string,
    toDeployment: string | undefined,
    fromDeployment: string | undefined,
    label: string,
    ports: string[],
    hidden?: boolean,
    deniedConnection: boolean,
    // Whether the connection was expected by the declarative config
    // (expectedAction === "Allow"). Drives the green/orange edge color.
    expected?: boolean

}

export interface SimpleGraphEdge {
    from: string,
    to: string,
    id: string,
    toDeployment: string | undefined,
    fromDeployment: string | undefined,
    label: string,
    port: string,
    hidden?: boolean,
    deniedConnection: boolean
    timestamp: number
    // Whether the connection was expected by the declarative config.
    expected?: boolean

}

export interface Graph {
    nodes: GraphNode[],
    edges: GraphEdge[]
}

export interface GenericGraph<T extends GraphEdge> extends Graph {
    nodes: GraphNode[],
    edges: T[]
}
