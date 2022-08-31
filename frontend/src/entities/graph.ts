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
    deniedConnection: boolean

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

}

export interface Graph {
    nodes: GraphNode[],
    edges: GraphEdge[]
}

export interface GenericGraph<T extends GraphEdge> extends Graph {
    nodes: GraphNode[],
    edges: T[]
}
