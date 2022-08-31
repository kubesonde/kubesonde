export interface ProbeOutput {
    start: string,
    end: string,
    items: ProbeOutputItem[],
    errors: ProbeOutputError[],
    podNetworking?: PodNetworkingInfo[]
    podNetworkingv2: PodNetworkingInfoV2
    podConfigurationNetworking: PodNetworkingInfoV2
}

export interface PodNetworkingInfoV2 { [name: string]: PodNetwotkingItem[] }
export interface PodNetwotkingItem {
    ip: string,
    port: string,
    protocol: string
}
export interface PodNetworkingInfo {
    podName: string,
    netstat: string
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
}
