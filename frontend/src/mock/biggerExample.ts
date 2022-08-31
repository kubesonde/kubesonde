import { ProbeEndpointType, ProbeOutput, ProbeOutputItem, ProbeOutputType } from "../entities/probeOutput";


const allowHTTP = {
    type: ProbeOutputType.PROBE,
    port: "80",
    protocol: "TCP",
    timestamp: 1234,
    expectedAction: "Allow",
    resultingAction: "Allow",
    source: {
        type: ProbeEndpointType.POD,
        name: "pod1",
        namespace: "default"
    },
    destination: {
        type: ProbeEndpointType.POD,
        name: "pod2",
        namespace: "default"
    },
    destinationHostnames: []
}

const allowHTTPS = {
    type: ProbeOutputType.PROBE,
    port: "443",
    protocol: "TCP",
    timestamp: 1234,
    expectedAction: "Allow",
    resultingAction: "Allow",
    source: {
        type: ProbeEndpointType.POD,
        name: "pod1",
        namespace: "default",
        deploymentName: "Deployment-2"
    },
    destination: {
        type: ProbeEndpointType.POD,
        name: "pod2",
        namespace: "default"
    },
    destinationHostnames: []
}

const allow8080: ProbeOutputItem = {
    type: ProbeOutputType.PROBE,
    port: "8080",
    protocol: "TCP",
    timestamp: 1234,
    expectedAction: "Allow",
    resultingAction: "Allow",
    source: {
        type: ProbeEndpointType.POD,
        name: "pod1",
        deploymentName: "Deployment-2",
        namespace: "default"
    },
    destination: {
        type: ProbeEndpointType.POD,
        name: "pod2",
        deploymentName: "Deployment-1",
        namespace: "default"
    },
    destinationHostnames: []
}

const deny8888: ProbeOutputItem = {
    type: ProbeOutputType.PROBE,
    port: "8888",
    protocol: "TCP",
    timestamp: 1234,
    expectedAction: "Deny",
    resultingAction: "Deny",
    source: {
        type: ProbeEndpointType.POD,
        name: "pod1",
        namespace: "default",
        deploymentName: "Deployment-2"
    },
    destination: {
        type: ProbeEndpointType.POD,
        name: "pod2",
        namespace: "default",
        deploymentName: "Deployment-1"
    },
    destinationHostnames: []
}

export const CompleteExample: ProbeOutput = {
    podConfigurationNetworking: {},
    podNetworkingv2: {},
    start: "now",
    end: "then",
    errors: [],
    items: [
        deny8888,
        { ...allowHTTPS, source: { type: ProbeEndpointType.POD, name: "pod2", deploymentName: "Deployment-1", namespace: "default" }, destination: { type: ProbeEndpointType.POD, name: "pod4", deploymentName: "Deployment-2", namespace: "default" } },
        { ...allow8080, source: { type: ProbeEndpointType.POD, name: "pod2", deploymentName: "Deployment-1", namespace: "default" }, destination: { type: ProbeEndpointType.POD, name: "pod4", deploymentName: "Deployment-2", namespace: "default" } },
        { ...allow8080, source: { type: ProbeEndpointType.POD, name: "pod4", deploymentName: "Deployment-2", namespace: "default" }, destination: { type: ProbeEndpointType.POD, name: "pod2", deploymentName: "Deployment-1", namespace: "default" } },
        { ...allowHTTP, source: { type: ProbeEndpointType.POD, name: "pod3", deploymentName: "Deployment-1", namespace: "default" }, destination: { type: ProbeEndpointType.POD, name: "pod4", deploymentName: "Deployment-2", namespace: "default" } },
        { ...allowHTTP, source: { type: ProbeEndpointType.POD, name: "pod2", deploymentName: "Deployment-1", namespace: "default" }, destination: { type: ProbeEndpointType.POD, name: "pod3", deploymentName: "Deployment-1", namespace: "default" } },
        { ...allow8080, destination: { type: ProbeEndpointType.POD, name: "pod4", deploymentName: "Deployment-2", namespace: "default" } },
    ]
}
