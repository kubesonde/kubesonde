import { useLocation } from "react-router-dom";
import { GraphBase } from "src/components/graph/graphBase/GraphBase";
import {
  buildEdgesFromProbes,
  buildNodesFromProbes,
  getErrorLogs,
  cleanupProbeOutput,
  ProbeErrorInfo,
} from "src/utils/probes";
import { PodNetworkingInfoV2, ProbeOutput } from "src/entities/probeOutput";
import { ErrorLogTable } from "src/components/table/graphTable/errorLogTable";
import React from "react";
import { WithSwitch } from "src/hoc/WithSwitch";
import { StatsTable } from "src/components/table/statsTable/statsTable";
import { mergeEdgesSimple } from "src/utils/graph";
import { NetInfoTable } from "../table/graphTable/netinfoTable";

const isPresent = (errorLog: ProbeErrorInfo[] | undefined): boolean => {
  return errorLog !== null && errorLog !== undefined && errorLog.length > 0;
};

const netinfo2Table = (netInfo: PodNetworkingInfoV2) => {
  const unsortedEntries = Object.entries(netInfo)
    .map(([key, value]) => {
      if (value === undefined || value === null) {
        return [];
      }
      const retval = value.map((entry) => ({ ...entry, podName: key }));
      return retval;
    })
    .flat();

  const sortedEntries = unsortedEntries.sort((a, b) => {
    if (a === undefined) {
      return 1;
    }
    if (b === undefined) {
      return -1;
    }
    const cmp1 = a.podName.localeCompare(b.podName);
    if (cmp1 === 0) {
      return parseInt(a.port) - parseInt(b.port);
    }
    return 0;
  });
  return sortedEntries;
};

export const GraphFromLocation = () => {
  const { state } = useLocation();
  // @ts-ignore
  const title = state.title;
  // @ts-ignore
  const data = cleanupProbeOutput(state.data as ProbeOutput);
  const nodes = buildNodesFromProbes(data);
  const edges = buildEdgesFromProbes(data);
  const errorLog = getErrorLogs(data);
  const netInfoContainers = data.podNetworkingv2;
  const netInfoContainersData =
    netInfoContainers && netinfo2Table(netInfoContainers);

  const netInfoDescriptionBase =
    data.podConfigurationNetworking &&
    netinfo2Table(data.podConfigurationNetworking);
  const netInfoDescriptionWithDuplicates = Array.from(
    new Set(
      netInfoDescriptionBase?.map((item) => ({
        port: item.port,
        protocol: item.protocol,
        podName: item.podName,
      }))
    )
  ).map((entry) => ({ ...entry, ip: "0.0.0.0" }));
  const netInfoDescription = netInfoDescriptionWithDuplicates.filter(
    (value, index) => {
      const _value = JSON.stringify(value);
      return (
        index ===
        netInfoDescriptionWithDuplicates.findIndex((obj) => {
          return JSON.stringify(obj) === _value;
        })
      );
    }
  );
  const ErrorLogTableWithSwitch = WithSwitch(ErrorLogTable, "error log");
  const StatsTableWithSwitch = WithSwitch(StatsTable, "graph stats");
  const NetInfoTableWithSwitch = WithSwitch(NetInfoTable, "open ports in pods");
  const NetInfoKubesondeTableWithSwitch = WithSwitch(
    NetInfoTable,
    "declarative network configuration"
  );
  return (
    <>
      <GraphBase
        title={title}
        nodes={nodes}
        edges={edges}
        podNetworkInfo={netInfoContainers}
        declarativeConfiguration={data.podConfigurationNetworking}
      />
      <StatsTableWithSwitch edges={mergeEdgesSimple(edges)} nodes={nodes} />
      {isPresent(errorLog) && (
        <ErrorLogTableWithSwitch errorLog={errorLog ?? []} />
      )}
      {netInfoContainersData && (
        <NetInfoTableWithSwitch {...{ data: netInfoContainersData }} />
      )}
      {netInfoDescription && (
        <NetInfoKubesondeTableWithSwitch {...{ data: netInfoDescription }} />
      )}
    </>
  );
};
