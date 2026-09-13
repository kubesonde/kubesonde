export interface ProbeOutput {
    start: string,
    end: string,
    items: ProbeOutputItem[],
    errors: ProbeOutputError[],
    podNetworking: PodNetworkingInfo
    podConfigurationNetworking: PodNetworkingInfo
}

export interface PodNetworkingInfo { [name: string]: PodNetwotkingItem[] }
export interface PodNetwotkingItem {
    ip: string,
    port: string,
    protocol: string
}

export enum ProbeOutputType {
    PROBE = "Probe",
    INFORMATION = "Information"
}
export interface ProbeOutputItem {
    type: ProbeOutputType,
    expectedAction: string,
    resultingAction: string,
    source: ProbeEndpointInfo,
    destination: ProbeEndpointInfo,
    destinationHostnames: string[]
    port: string
    forwardedPort?: string
    protocol: string,
    timestamp: number

}

export interface ProbeOutputError {
    value: ProbeOutputItem,
    reason: string

}

export enum ProbeEndpointType {
    POD = "Pod",
    SERVICE = "Service",
    INTERNET = "Internet"

}
export interface ProbeEndpointInfo {
    type: ProbeEndpointType
    name: string,
    namespace: string,
    IPAddress?: string,
    deploymentName?: string,
    replicaSetName?: string,
    selector?: string,
}
