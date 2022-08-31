import { PodNetworkingInfoV2, ProbeEndpointInfo, ProbeOutput, ProbeOutputError, ProbeOutputItem } from "../entities/probeOutput";
import { GraphNode, SimpleGraphEdge } from "../entities/graph";
import { NetstatInterface, parseNetstat } from "./netstat";
import { Dict } from "src/entities/types";

//const PUBLIC_DNS = "Public DNS"

export interface NetstatInfo {
    name: string,
    entries: NetstatInterface[]
}

export interface ProbeErrorInfo {
    podName: string,
    reason: string,
    timestamp: number
}

const buildGroupMap = (input_probes: ProbeOutput): Map<string, string> => {
    const groupMap = new Map<string, string>()

    input_probes.items.forEach((item) => {

        if (item.source.deploymentName) {
            groupMap.set(item.source.name, item.source.deploymentName)
        }
        if (item.destination.deploymentName) {
            groupMap.set(item.destination.name, item.destination.deploymentName)
        }
    })

    return groupMap
}

const toSimpleEdge = (groupMap: Map<string, string>) => (probe: ProbeOutputItem, index: number): SimpleGraphEdge => ({
    id: index.toString(),
    from: probe.source.name, //source.name.endsWith('DNS') ? PUBLIC_DNS : probe.source.name,
    to: probe.destination.name,//probe.destination.name.endsWith('DNS') ? PUBLIC_DNS : probe.destination.name,
    label: probe.port ? probe.port.toString() : "UNKNOWN",
    fromDeployment: groupMap.get(probe.source.name),
    toDeployment: groupMap.get(probe.destination.name),
    port: probe.forwardedPort ? `${probe.port}:${probe.forwardedPort}/${probe.protocol}` : `${probe.port}/${probe.protocol}`,
    timestamp: probe.timestamp,
    deniedConnection: probe.resultingAction === "Deny" ? true : false
})

const toErrorEdge = (groupMap: Map<string, string>) => (probe: ProbeOutputError, index: number): SimpleGraphEdge => ({
    id: index.toString() + "disallowed",
    from: probe.value.source.name, //source.name.endsWith('DNS') ? PUBLIC_DNS : probe.source.name,
    to: probe.value.destination.name,//probe.destination.name.endsWith('DNS') ? PUBLIC_DNS : probe.destination.name,
    label: probe.reason.toString(),
    fromDeployment: groupMap.get(probe.value.source.name),
    toDeployment: groupMap.get(probe.value.destination.name),
    port: probe.value.port ?? "None",
    hidden: true,
    deniedConnection: true,
    timestamp: probe.value.timestamp
})

export function getNetstatInfoFromProbes(input_probes: ProbeOutput): NetstatInfo[] | undefined {
    const netstatMap = input_probes.podNetworking?.map((item) =>
        ({ name: item.podName, entries: item.netstat.split('\n').slice(2) }))
    return netstatMap?.map((entry) => ({
        name: entry.name,
        entries: entry?.entries.map((e) => parseNetstat(e))
    }))
}

export function buildNodesFromProbes(input_probes: ProbeOutput): GraphNode[] {
    const regular_nodes: ProbeEndpointInfo[] = input_probes.items.map(probe => [probe.source, probe.destination]).flat()
    const error_nodes: ProbeEndpointInfo[] = input_probes.errors.map(probe => [probe.value.source, probe.value.destination]).flat()

    const individualNodes = Array.from(new Set(regular_nodes.concat(error_nodes)))
    const individualNodesNames = individualNodes.map(node => node.name)
    const buildTitle = (node: string): string | undefined => {
        const nodeInfo = (input_probes.items).concat(input_probes.errors.map(e => e.value)).find((item) => item.source.name === node)
        if (!nodeInfo) {
            return undefined
        }
        const namespace = nodeInfo.source.namespace
        const name = nodeInfo.source.name
        const hostNames = nodeInfo.destinationHostnames
        const deployment = nodeInfo.source.deploymentName
        const title = JSON.stringify({
            name,
            namespace,
            hostNames,
            deployment
        }, null, 2)
        return title
    }
    const groupMap = buildGroupMap(input_probes)
    const nodesFromAllowedProbes = individualNodes.map((item) => (
        {
            id: item.name,
            name: item.name,
            label: item.name,
            deployment: groupMap.get(item.name),
            group: groupMap.get(item.name),
            title: buildTitle(item.name),
            type: item.type as string
        }))

    const disallowedProbes = input_probes.items
        .filter((probe) => probe.resultingAction === "Deny")
        .map((probe) => ([probe.source, probe.destination])).flat()
    const nodesToBeAdded = Array.from(new Set(disallowedProbes.filter((node) => !individualNodesNames.includes(node.name))))
        .map((item) => (
            {
                id: item.name,
                name: item.name,
                label: item.name,
                hidden: true,
                deployment: groupMap.get(item.name),
                group: groupMap.get(item.name),
                title: buildTitle(item.name),
                type: item.type as string
            }))
    const retval = [...nodesFromAllowedProbes, ...nodesToBeAdded]
    const uniqueArray = retval.filter((value) => value.id
        !== undefined).filter((value, index) => {
            const _value = JSON.stringify(value);
            return index === retval.findIndex(obj => {
                return JSON.stringify(obj) === _value;
            });
        });
    return uniqueArray
}

interface EdgeDict { [key: string]: SimpleGraphEdge }
function removeDuplicates(allEdges: SimpleGraphEdge[]): SimpleGraphEdge[] {

    const maps = allEdges.reduce((acc, curr) => {
        const id: string = `${curr.from}-${curr.to}-${curr.port}-${curr.label}`
        if (id in acc) {
            const prev = acc[id]
            if (prev.timestamp > curr.timestamp) {
                return acc
            } else {
                acc[id] = curr
                return acc
            }
        } else {
            acc[id] = curr
            return acc
        }

    }, {} as EdgeDict)

    return Object.values(maps).flat()
}

export function buildEdgesFromProbes(probes: ProbeOutput): SimpleGraphEdge[] {
    const groupMap = buildGroupMap(probes)

    const allowedEdges: SimpleGraphEdge[] = probes.items
        .filter((probe) => probe.resultingAction !== "Deny")
        .map(toSimpleEdge(groupMap))
    const disallowedEdges: SimpleGraphEdge[] = probes.items
        .filter((probe) => probe.resultingAction === "Deny")
        .map(toSimpleEdge(groupMap))
        .map((edge) => ({ ...edge, id: edge.id + "disallowed", hidden: true }))
    const errorEdges = probes.errors.map(toErrorEdge(groupMap))
    const retval = removeDuplicates([...allowedEdges, ...disallowedEdges, ...errorEdges])
    /*const retvalSameId = retval.map(item => ({ ...item, id: '1' }))
    const uniqueArray = retvalSameId.filter((value, index) => {
        const _value = JSON.stringify(value);
        return index === retvalSameId.findIndex(obj => {
            return JSON.stringify(obj) === _value;
        });
    });*/
    const uniqueId = retval.map((item, index) => ({ ...item, id: index.toString() }))
    return uniqueId
}

export function getErrorLogs(input_probes: ProbeOutput): ProbeErrorInfo[] | undefined {
    return Array.from(new Set(input_probes.errors?.filter((error) => error.reason !== undefined)
        .map((probe) => ({
            reason: probe.reason,
            podName: probe.value.source.name,
            timestamp: probe.value.timestamp
        }))))
}

export function getPodIPMappingFromProbes(input_probes: ProbeOutput): Dict | undefined {
    return input_probes.items.map((item) => [item.source, item.destination])
        .flat(1)
        .filter((item) => item.IPAddress !== undefined)
        .map((entry) => ({
            name: entry.name,
            ip: entry.IPAddress
        })).reduce((acc, curr) => {
            acc[curr.ip as string] = curr.name
            return acc
        }, {} as Dict)


}

export function cleanupNetInfo(data: PodNetworkingInfoV2): PodNetworkingInfoV2 {
    return Object.entries(data).map(([key, value]) => {
        return { key, value: value.filter((item) => item.ip !== "127.0.0.1") }
    }).reduce((acc, curr) => {
        acc[curr.key] = curr.value
        return acc
    }, {} as PodNetworkingInfoV2)

}

interface ProbeDict { [key: string]: ProbeOutputItem }
export function cleanupProbes(items: ProbeOutputItem[]): ProbeOutputItem[] {
    const probesDict = items.reduce((acc: ProbeDict, curr: ProbeOutputItem) => {
        const key = `${curr.destination.name}-${curr.source.name}-${curr.port}-${curr.protocol}`
        if (key in acc) {
            const previousItem = acc[key]
            if (curr.timestamp > previousItem.timestamp) {
                acc[key] = curr
            }
            return acc
        }
        acc[key] = curr
        return acc
    }, {} as ProbeDict)

    return Object.values(probesDict)
}


/**This function removes the noise created by the backend when running tests.
 * It: 
 * - Removes the ports listening to localhost 
 * - Removes conflicting entries in the probes by selecting always the newest probes 
 */
export function cleanupProbeOutput(input_probes: ProbeOutput): ProbeOutput {
    return {
        ...input_probes,
        items: cleanupProbes(input_probes.items),
        podNetworkingv2: cleanupNetInfo(input_probes.podNetworkingv2)
    }
}