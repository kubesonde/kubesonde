import { IBuilder, Builder } from "builder-pattern";
import { ProbeOutputItem, ProbeOutputType, ProbeEndpointType } from "src/entities/probeOutput";

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


export function ProbeOutputItemBuilder(): IBuilder<ProbeOutputItem> {
    return Builder<ProbeOutputItem>(allow8080)
}