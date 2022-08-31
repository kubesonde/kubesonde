import { PodNetworkingInfoV2, ProbeEndpointType, ProbeOutput, ProbeOutputItem, ProbeOutputType } from "../entities/probeOutput";
import { buildEdgesFromProbes, cleanupNetInfo, cleanupProbes } from "./probes";
import { ProbeOutputItemBuilder } from "./probes.builder";

describe('Edge creation', function () {
    it('Groups edges with same source an destination', () => {
        const output: ProbeOutput = {
            podConfigurationNetworking: {},
            podNetworkingv2: {},
            start: "start",
            end: "end",
            errors: [],
            items: [
                {
                    type: ProbeOutputType.PROBE,
                    expectedAction: "Allow",
                    resultingAction: "Allow",
                    destinationHostnames: [],
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
                    port: "8080",
                    protocol: "TCP",
                    timestamp: 1234
                }, {
                    type: ProbeOutputType.PROBE,
                    expectedAction: "Allow",
                    resultingAction: "Allow",
                    destinationHostnames: [],
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
                    port: "80",
                    protocol: "TCP",
                    timestamp: 1234
                }]
        }

        const result = buildEdgesFromProbes(output)
        expect(result).toEqual([{
            from: "pod1",
            to: "pod2",
            id: "0",
            fromDeployment: undefined,
            toDeployment: undefined,
            label: "8080",
            port: "8080/TCP",
            timestamp: 1234,
            deniedConnection: false
        },
        {
            from: "pod1",
            to: "pod2",
            fromDeployment: undefined,
            toDeployment: undefined,
            id: "1",
            label: "80",
            port: "80/TCP",
            timestamp: 1234,
            deniedConnection: false
        }
        ])
    })
    it('Does not group edges with different source an destination', () => {
        const output: ProbeOutput = {
            podConfigurationNetworking: {},
            podNetworkingv2: {},
            start: "start",
            end: "end",
            errors: [],
            items: [
                {
                    type: ProbeOutputType.PROBE,
                    expectedAction: "Allow",
                    resultingAction: "Allow",
                    destinationHostnames: [],
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
                    port: "8080",
                    protocol: "TCP",
                    timestamp: 1234
                }, {
                    type: ProbeOutputType.PROBE,
                    expectedAction: "Allow",
                    resultingAction: "Allow",
                    destinationHostnames: [],
                    source: {
                        type: ProbeEndpointType.POD,
                        name: "pod2",
                        namespace: "default"
                    },
                    destination: {
                        type: ProbeEndpointType.POD,
                        name: "pod1",
                        namespace: "default"
                    },
                    port: "80",
                    protocol: "TCP",
                    timestamp: 1234
                }]
        }

        const result = buildEdgesFromProbes(output)
        expect(result).toEqual([{
            from: "pod1",
            to: "pod2",
            id: "0",
            label: "8080",
            port: "8080/TCP",
            timestamp: 1234,
            deniedConnection: false
        },
        {
            from: "pod2",
            to: "pod1",
            id: "1",
            label: "80",
            port: "80/TCP",
            timestamp: 1234,
            deniedConnection: false
        }
        ])
    })
});

describe('Cleanup functions', () => {
    it('removes localhost interfaces from netinfo', () => {
        const netinfo: PodNetworkingInfoV2 = {
            'pod1': [{
                ip: "127.0.0.1",
                port: "80",
                protocol: "TCP"
            }, {
                ip: "0.0.0.0",
                port: "8080",
                protocol: "TCP"
            }]
        }
        const expectedNetinfo: PodNetworkingInfoV2 = {
            'pod1': [{
                ip: "0.0.0.0",
                port: "8080",
                protocol: "TCP"
            }]
        }
        expect(cleanupNetInfo(netinfo)).toEqual(expectedNetinfo)

    })
    it('cleanup probeItems', () => {
        const items: ProbeOutputItem[] = [
            ProbeOutputItemBuilder().port("80").timestamp(1).resultingAction("Deny").build(),
            ProbeOutputItemBuilder().port("80").timestamp(2).resultingAction("Allow").build()
        ]
        expect(cleanupProbes(items)).toEqual([ProbeOutputItemBuilder().port("80").timestamp(2).resultingAction("Allow").build()])
    })
})