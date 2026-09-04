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

// Probes exercising the "unauthorized port" case: the port responds to a
// probe and is open according to netstat, but is absent from the pod's
// declarative (manifest) configuration.
const allowUnauthorizedOnPod2: ProbeOutputItem = {
    ...allow8080,
    port: "22",
    source: { type: ProbeEndpointType.POD, name: "pod1", namespace: "default" },
    destination: { type: ProbeEndpointType.POD, name: "pod2", deploymentName: "Deployment-1", namespace: "default" },
}

const allowUnauthorizedOnPod4: ProbeOutputItem = {
    ...allow8080,
    port: "9999",
    source: { type: ProbeEndpointType.POD, name: "pod2", deploymentName: "Deployment-1", namespace: "default" },
    destination: { type: ProbeEndpointType.POD, name: "pod4", deploymentName: "Deployment-2", namespace: "default" },
}

// Probes exercising the "declared but closed" case: the manifest declares
// the port, but netstat shows nothing actually listening on it.
const allowDeclaredButClosedOnPod2: ProbeOutputItem = {
    ...allow8080,
    port: "3000",
    source: { type: ProbeEndpointType.POD, name: "pod1", namespace: "default" },
    destination: { type: ProbeEndpointType.POD, name: "pod2", deploymentName: "Deployment-1", namespace: "default" },
}

const allowDeclaredButClosedOnPod4: ProbeOutputItem = {
    ...allow8080,
    port: "6000",
    source: { type: ProbeEndpointType.POD, name: "pod1", deploymentName: "Deployment-2", namespace: "default" },
    destination: { type: ProbeEndpointType.POD, name: "pod4", deploymentName: "Deployment-2", namespace: "default" },
}

export const CompleteExample: ProbeOutput = {
    // Declarative configuration (e.g. containerPorts from the Kubernetes manifests).
    podConfigurationNetworking: {
        pod2: [
            { ip: "10.0.0.2", port: "80", protocol: "TCP" },
            { ip: "10.0.0.2", port: "3000", protocol: "TCP" },
        ],
        pod4: [
            { ip: "10.0.0.4", port: "8080", protocol: "TCP" },
            { ip: "10.0.0.4", port: "6000", protocol: "TCP" },
        ],
    },
    // Ports actually observed listening on the pods (e.g. via netstat).
    podNetworkingv2: {
        pod2: [
            { ip: "10.0.0.2", port: "80", protocol: "TCP" },
            { ip: "10.0.0.2", port: "22", protocol: "TCP" },
        ],
        pod4: [
            { ip: "10.0.0.4", port: "8080", protocol: "TCP" },
            { ip: "10.0.0.4", port: "9999", protocol: "TCP" },
        ],
    },
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
        allowUnauthorizedOnPod2,
        allowUnauthorizedOnPod4,
        allowDeclaredButClosedOnPod2,
        allowDeclaredButClosedOnPod4,
    ]
}
