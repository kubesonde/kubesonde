import { PodNetworkingInfoV2, PodNetwotkingItem } from "src/entities/probeOutput";
import { SimpleGraphEdge } from "../entities/graph";
import { NetstatInfo } from "./probes";

export interface PortMapping {
    podName: string
    ports: string[]
}

export const getMappingFromEdges = (edges: SimpleGraphEdge[]): PortMapping[] => {
    const pods = Array.from(new Set(edges.map((e) => e.to)))
    const mapping = pods.map((podName) => ({
        podName,
        ports: edges.filter((e) => e.to === podName).map((e) => e.port)
    }))
    return mapping
}

export const getMappingFromNetstat = (netstat: NetstatInfo[]): PortMapping[] => {
    const mapping = netstat.map((netstatLine) => ({
        podName: netstatLine.name,
        // FIXME: this is a hack because I am not sure of how to parse the output
        ports: Array.from(new Set(netstatLine.entries.filter((e) => e !== undefined && e.local !== undefined && e.state === "LISTEN").map((e) => `${e.local.port}/${e.protocol.toUpperCase()}`)))
    }))
    return mapping
}

export const getMappingFromNetInfo = (netstat: PodNetworkingInfoV2): PortMapping[] => {
    const mapping: PortMapping[] = Object.entries<PodNetwotkingItem[]>(netstat).map(value => {
        return {
            podName: value[0],
            ports: value[1].map(entry => entry.port + "/" + entry.protocol) // TODO: FIX this 
        }
    })

    return mapping
}